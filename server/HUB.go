package server

import "sync"

// Hub is a helper struct that keeps all currently running vaults inside
//
// should be thread-safe thanks to sync.Map use
type Hub struct {
	Vaults sync.Map
}
