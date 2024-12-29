package server

import (
	"github.com/gorilla/websocket"
	"log"
)

// Connection is a helper struct which wraps the websocket connection and adds a buffered channel to it for outgoing messages
//
// to write to this connection, send data to the Send channel
type Connection struct {
	uuid string
	conn *websocket.Conn
	Send chan []byte // buffered channel for outgoing messages
}

// WritePump manages asynchronous writing of messages to the client connection using buffered channel Send
//
// source https://brojonat.com/posts/websockets/, use if ping needed
func (c *Connection) WritePump() {
	log.Println("starting write pump on " + c.uuid)
	for {
		msgBytes, ok := <-c.Send
		// ok will be false in case the Send channel is closed
		if !ok {
			// channel is closed, send close message and return
			_ = c.conn.WriteMessage(websocket.CloseMessage, nil)
			return
		}
		// write a message to the connection
		log.Println("writing message to connection")
		if err := c.conn.WriteMessage(websocket.BinaryMessage, msgBytes); err != nil {
			// close in case of errors
			log.Println("error, closing pump")
			return
		}

	}
}
