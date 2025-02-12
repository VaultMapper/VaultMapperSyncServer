package server

import (
	"log"
	"net/http"
	"regexp"
	"time"

	pb "github.com/NodiumHosting/VaultMapperSyncServer/proto"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var uuidRegex, _ = regexp.Compile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$")
var vaultIDRegex, _ = regexp.Compile("^vault_[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$")

func handshakeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("upgrade") != "websocket" { // guard against non websocket requests
		return
	}

	uuid := r.URL.Query().Get("uuid")
	vaultID := r.URL.Query().Get("vaultID")      // if checks pass, upgrade
	sendViewIDToast := r.URL.Query().Get("view") // if any value present(key is found), sends a toast to the player when joining vault
	log.Printf(vaultID + ": " + uuid)

	if !uuidRegex.MatchString(uuid) || !vaultIDRegex.MatchString(vaultID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		log.Println("Check not passed: " + uuid)
		return // close the ws basically..
	}
	log.Printf("Connection successful: " + uuid)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	// after this point is the websocket loop
	// conn.ReadMessage() reads the message, works like onMessage
	// use onClose to do stuff after closing socket

	SendVault(vaultID, conn) // send vault to client

	c, vault := HUB.AddConnectionToVault(vaultID, uuid, conn)
	if c == nil { // close connection/return if
		_ = conn.WriteMessage(websocket.CloseMessage, nil)
		err := conn.Close()
		if err != nil {
			return
		}
		return
	}

	defer onClose(uuid, vaultID)

	_ = AddPlayerToVault(uuid, vaultID) // add player to vault db

	// send viewer code
	msg := &pb.Message{
		Type:       pb.MessageType_VIEWER_CODE,
		ViewerCode: &pb.ViewerCode{Code: GetVaultViewCode(vaultID)},
	}
	msgBuffer, err := proto.Marshal(msg)
	c.Send <- msgBuffer

	if sendViewIDToast != "" { // send Toast with ViewerID if found
		c.SendToast("Viewer Code: " + vault.ViewerCode)
	}

	isDead := 0

	// this should basically be the onMessage thingy
	for {
		err := conn.SetReadDeadline(time.Now().Add(15 * time.Second))
		if err != nil {
			return
		}
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read: " + err.Error())
			return
		}

		if string(data) == "keep_me_alive" {
			//log.Println("Keep alive received")
			isDead++

			if isDead > 30 { // if connection keeps sending keepalives, but no other data arives, the player is afk or the client bugged out - disconnect them
				log.Println("CONNECTION IS DEAD: " + uuid)
				conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			continue
		}
		isDead = 0

		var msg pb.Message
		err2 := proto.Unmarshal(data, &msg)
		if err2 != nil {
			log.Println("Marshal problem")
			return
		}

		onMessage(vaultID, uuid, &msg, c)
		//log.Println("onmessage adding in")
		//inPacketCounterChan <- 1
		//log.Println("onmessage adding in done")

	}
}

func onMessage(vaultID string, uuid string, msg *pb.Message, conn *Connection) {
	msgType := msg.GetType()
	switch msgType {
	case pb.MessageType_VAULT_PLAYER:
		// this case handles accepted player packet
		HandlePlayerMovement(vaultID, uuid, msg)
		break
	case pb.MessageType_VAULT_CELL:
		// This case should handle accepted VaultCell
		HandleVaultCell(vaultID, uuid, msg)
		break
	case pb.MessageType_PLAYER_DISCONNECT:
		// this shouldn't happen as PlayerDisconnect is only S2C
		log.Println(uuid + " tried to send PlayerDisconnect which shouldn't happen")
		break
	case pb.MessageType_VAULT:
		// this shouldn't happen as the Vault is only S2C
		log.Println(uuid + " tried to send Vault which shouldn't happen")
		break
	case pb.MessageType_TOAST:
		// this shouldn't happen as the Toast is only S2C
		log.Println(uuid + " tried to send Toast which shouldn't happen")
		break
	case pb.MessageType_VIEWER_CODE_REQUEST:
		HandleViewerCodeRequest(vaultID, conn)
		break
	default:
		log.Println(uuid + " sent unknown packet")
		break
	}
}

func onClose(uuid string, vaultID string) { // need to send down PlayerDisconnect to others in vault here
	log.Println(uuid + " closed connection to vault: " + vaultID)
	msg := pb.Message{
		Type:             pb.MessageType_PLAYER_DISCONNECT,
		PlayerDisconnect: &pb.PlayerDisconnect{Uuid: uuid},
	}
	BroadcastMessage(vaultID, uuid, &msg)
	HUB.RemoveConnectionFromVault(vaultID, uuid)
}

func HandleViewerCodeRequest(vaultID string, conn *Connection) {
	code := GetVaultViewCode(vaultID)
	if code != "" {
		log.Println("Sending viewer code: " + code)
		msg := pb.Message{
			Type:       pb.MessageType_VIEWER_CODE,
			ViewerCode: &pb.ViewerCode{Code: code},
		}
		messageBuffer, err := proto.Marshal(&msg)
		if err != nil {
			return
		}
		conn.Send <- messageBuffer
	}
}

// HandlePlayerMovement handles incoming PlayerMovement packets from clients and broadcasts them to the other players
func HandlePlayerMovement(vaultID string, uuid string, msg *pb.Message) {
	//log.Println("Handling PlayerMovement")
	BroadcastMessage(vaultID, uuid, msg)
}

// HandleVaultCell handles incoming VaultCell packets from clients, broadcasts them to the other players and adds them to internal structures
func HandleVaultCell(vaultID string, uuid string, msg *pb.Message) {
	log.Println("Handling VaultCell")
	cell := msg.GetVaultCell()

	BroadcastMessage(vaultID, uuid, msg)
	vault := HUB.GetVault(vaultID)
	if vault == nil {
		return
	}

	vault.AddOrReplaceCell(cell)
}

// BroadcastMessage is used to broadcast Message to a vault, with excludeUUID being excluded
func BroadcastMessage(vaultID string, excludeUUID string, msg *pb.Message) {
	vault := HUB.GetVault(vaultID) // get vault
	if vault == nil {
		return
	}
	messageBuffer, err := proto.Marshal(msg) // serialize message into buffer
	if err != nil {
		return
	}

	vault.Connections.Range(func(key, val interface{}) bool { // go through connections and add to their Send channels
		if key != excludeUUID {
			conn := val.(*Connection)
			conn.Send <- messageBuffer
		}
		return true
	})
	vault.Viewers.Range(func(key, val interface{}) bool { // go through viewers and add to their Send channels
		conn := val.(*Connection)
		conn.Send <- messageBuffer
		//log.Println("sending to viewer")

		return true
	})
}

// SendVault sends all the Vault.Cells using the Vault message type initially to sync vault to client if joined after start
func SendVault(vaultID string, conn *websocket.Conn) {
	vault := HUB.GetVault(vaultID)
	if vault == nil {
		log.Println("Tried to send vault that doesn't exist, it will be created now")
		return // if vault doesn't exist, do nothing - this can happen when this is the first player joining a fresh vault
	}
	log.Println("Sending vault to client")
	var cells []*pb.VaultCell
	vault.Cells.Range(func(key, val interface{}) bool {
		cells = append(cells, val.(*pb.VaultCell))
		//log.Println("appended cell")
		return true
	})

	msg := pb.Message{
		Type:  pb.MessageType_VAULT,
		Vault: &pb.Vault{Cells: cells},
	}

	messageBuffer, err := proto.Marshal(&msg)
	if err != nil {
		return
	}

	errr := conn.WriteMessage(websocket.BinaryMessage, messageBuffer)
	if errr != nil {
		return
	}
}
