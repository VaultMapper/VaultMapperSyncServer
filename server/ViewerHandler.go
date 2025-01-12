package server

import (
	"github.com/NodiumHosting/VaultMapperSyncServer/util"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var viewerUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var VIEWER_LIMIT = 10

func viewerHandshakeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("upgrade") != "websocket" { // guard against non websocket requests
		return
	}

	code := r.URL.Query().Get("code")
	log.Println("e: ", code)

	// get vault from code here
	// each vault will have it's own viewer code
	vault := HUB.GetVaultByCode(code)
	if vault == nil {
		http.Error(w, "Unknown vault", http.StatusForbidden)
		log.Println("tried to connect to nonexistent vault")
		return
	}
	if vault.ViewerCount() > VIEWER_LIMIT {
		http.Error(w, "Viewer limit reached", http.StatusForbidden)
		log.Println("too many viewers")
		return
	}
	vaultUUID := vault.UUID

	log.Printf("Connection successful to vault code: " + code)

	conn, err := viewerUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	// after this point is the websocket loop
	// conn.ReadMessage() reads the message, works like onMessage
	// use onClose to do stuff after closing socket

	SendVault(vaultUUID, conn) // send vault to client

	connUUID := util.RandSeq(32)
	ok := vault.AddViewer(connUUID, conn)

	if !ok { // if not ok -> connection exists -> return/close connection
		_ = conn.WriteMessage(websocket.CloseMessage, nil)
		err := conn.Close()
		if err != nil {
			return
		}
		return
	}

	defer onViewerClose(connUUID, vault)

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

		if string(data) == "keep_me_alive" { // the only incoming packet we expect here is a keepalive
			//log.Println("Keep alive received")
			continue
		} else {
			return // closes handshake handler which triggers onClose
		}
	}
}

func onViewerClose(uuid string, vault *Vault) { // need to send down PlayerDisconnect to others in vault here
	log.Println(uuid + " closed connection to vault: " + vault.UUID)

	vault.RemoveViewer(uuid)
}
