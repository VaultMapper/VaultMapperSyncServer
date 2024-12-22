package server

import (
	"sync"
)

// Hub is a helper struct that keeps all currently running vaults inside
//
// should be thread-safe thanks to sync.Map use
type Hub struct {
	Vaults sync.Map
}

// CreateVault is a helper method that creates Vault inside Hub
func (h *Hub) CreateVault(vaultID string) {
	h.Vaults.LoadOrStore(vaultID, &Vault{
		UUID: vaultID,
	})
}

// RemoveVault is a helper method that removes Vault from Hub
//
// # Does not do any checks for if vault is empty
//
// Only call this if the Vault is empty, otherwise will leave dangling connections and send channels
func (h *Hub) RemoveVault(vaultID string) {
	h.Vaults.Delete(vaultID)
}
