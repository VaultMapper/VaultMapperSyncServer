package server

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

// Vault is a helper struct that is used to store connections to vault and other runtime data
//
// If other data is added, there should be a Mutex created with it to make sure writing is thread-safe
type Vault struct {
	UUID        string
	Connections sync.Map // stores a map of current connections inside the vault, key is uuid, value is Connection
}

// AddConnection adds the connection to Vault structure and starts up the WritePump
//
// ok is false if the connection already exists, else false
func (v *Vault) AddConnection(playerUUID string, conn *websocket.Conn) bool {
	_, ok := v.Connections.Load(playerUUID)
	if ok {
		log.Println("Tried to add connection but it already exists")
		return false // connection already exists
	}
	c := &Connection{ // create connection
		uuid: playerUUID,
		conn: conn,
		Send: make(chan []byte, 256), // buffered channel of 256 bytes
	}

	v.Connections.Store(playerUUID, c) // store the connection inside vault
	go c.WritePump()                   // Start the write pump
	return true
}

// RemoveConnection removes the connection from Vault structure and closes the Send channel
//
// return true if the Vault is empty after connection removal, false otherwise
//
// Send channel needs to be closed for WritePump to exit properly!
func (v *Vault) RemoveConnection(playerUUID string) bool {
	value, ok := v.Connections.Load(playerUUID)
	if !ok {
		return false
	}

	c := value.(*Connection)
	close(c.Send)         // close send channel
	err := c.conn.Close() // close connection
	if err != nil {
		return false
	}
	v.Connections.Delete(playerUUID) // remove connection

	// check if the vault is empty now
	isEmpty := true
	v.Connections.Range(func(k, v interface{}) bool {
		isEmpty = false
		return false
	})
	if isEmpty {
		return true
	}

	return false
}
