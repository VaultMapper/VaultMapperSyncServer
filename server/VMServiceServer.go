package server

import (
	"fmt"
	pb "github.com/NodiumHosting/VaultMapperSyncServer/proto"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"log"
	"regexp"

	"net/http"
)

// HUB stores all Vault structures
var HUB = Hub{}

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
	vaultID := r.URL.Query().Get("vaultID") // if checks pass, upgrade
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

	ok := HUB.AddConnectionToVault(vaultID, uuid, conn)
	if !ok { // if not ok -> connection exists -> return/close connection
		_ = conn.WriteMessage(websocket.CloseMessage, nil)
		err := conn.Close()
		if err != nil {
			return
		}
		return
	}

	defer onClose(uuid, vaultID)

	// this should basically be the onMessage thingy
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read: " + err.Error())
			return
		}
		var msg pb.Message
		err2 := proto.Unmarshal(data, &msg)
		if err2 != nil {
			log.Println("Marshal problem")
			return
		}
		onMessage(vaultID, uuid, &msg)
	}
}

func onMessage(vaultID string, uuid string, msg *pb.Message) {
	log.Printf("\nOn message from %s\ntype: %v\ndata: %v\n", uuid, msg.GetType(), msg.GetContent())
	msgType := msg.GetType()
	switch msgType {
	case pb.MessageType_VAULT_PLAYER:
		// this case handles accepted player packet
		handlePlayerMovement(vaultID, uuid, msg)
		break
	case pb.MessageType_VAULT_CELL:
		// This case should handle accepted VaultCell
		handleVaultCell(vaultID, uuid, msg)
		break
	case pb.MessageType_PLAYER_DISCONNECT:
		// this shouldn't happen as PlayerDisconnect is only S2C
		log.Println(uuid + " tried to send PlayerDisconnect which shouldn't happen")
		break
	case pb.MessageType_VAULT:
		// this shouldn't happen as the Vault is only S2C
		log.Println(uuid + " tried to send Vault which shouldn't happen")
		break
	default:
		log.Println(uuid + " sent unknown packet")
		break
	}
}

func onClose(uuid string, vaultID string) { // need to send down PlayerDisconnect to others in vault here
	log.Println(uuid + " closed connection to vault: " + vaultID)
	msg := pb.Message{
		Type: pb.MessageType_PLAYER_DISCONNECT,
	}
	msg.Content = &pb.PlayerDisconnect{Uuid: uuid}
	broadcastMessage(vaultID, uuid, &msg)
	HUB.RemoveConnectionFromVault(vaultID, uuid)
}

// handlePlayerMovement handles incoming PlayerMovement packets from clients and broadcasts them to the other players
func handlePlayerMovement(vaultID string, uuid string, msg *pb.Message) {
	log.Println("Handling PlayerMovement")
	broadcastMessage(vaultID, uuid, msg)
}

// handleVaultCell handles incoming VaultCell packets from clients, broadcasts them to the other players and adds them to internal structures
func handleVaultCell(vaultID string, uuid string, msg *pb.Message) {
	log.Println("Handling VaultCell")
	broadcastMessage(vaultID, uuid, msg)
}

// broadcastMessage is used to broadcast Message to a vault, with excludeUUID being excluded
func broadcastMessage(vaultID string, excludeUUID string, msg *pb.Message) {
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
}

func Run(ip string, port int) {
	fmt.Println("HELLO FROM SERVER")

	http.HandleFunc("/", handshakeHandler)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), nil); err != nil {
		log.Fatal(err)
	}
}
