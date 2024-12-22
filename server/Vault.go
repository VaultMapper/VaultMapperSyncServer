package server

import "sync"

// Vault is a helper struct that is used to store connections to vault and other runtime data
//
// If other data is added, there should be a Mutex created with it to make sure writing is thread-safe
type Vault struct {
	UUID        string
	Connections sync.Map // stores a map of current connections inside the vault, key is uuid, value is Connection
}
