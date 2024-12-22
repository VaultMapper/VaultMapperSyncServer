package server

import "github.com/gorilla/websocket"

// Connection is a helper struct which wraps the websocket connection and adds a buffered channel to it for outgoing messages
type Connection struct {
	conn *websocket.Conn
	Send chan []byte // buffered channel for outgoing messages
}
