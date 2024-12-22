package server

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

// Hub is a helper struct that keeps all currently running vaults inside
//
// should be thread-safe thanks to sync.Map use
type Hub struct {
	Vaults sync.Map
}

// GetOrCreateVault is a helper method that gets and optionally creates Vault inside Hub
func (h *Hub) GetOrCreateVault(vaultID string) *Vault {
	vault, _ := h.Vaults.LoadOrStore(vaultID, &Vault{
		UUID: vaultID,
	})
	v := vault.(*Vault) // assert type of vault
	return v
}

// GetVault is used to get Vault by UUID, returns nil if not found
func (h *Hub) GetVault(vaultID string) *Vault {
	vault, ok := h.Vaults.Load(vaultID)
	if !ok {
		log.Println("Tried to access vault that doesn't exist")
		return nil
	}
	v := vault.(*Vault)
	return v
}

// RemoveVault is a helper method that removes Vault from Hub
//
// # Does not do any checks for if vault is empty
//
// Only call this if the Vault is empty, otherwise will leave dangling connections and send channels
func (h *Hub) RemoveVault(vaultID string) {
	h.Vaults.Delete(vaultID)
}

func (h *Hub) AddConnectionToVault(vaultID string, playerUUID string, conn *websocket.Conn) bool {
	vault := h.GetOrCreateVault(vaultID)
	ok := vault.AddConnection(playerUUID, conn)
	if !ok {
		return false
	}
	return true
}

func (h *Hub) RemoveConnectionFromVault(vaultID string, playerUUID string) {
	vault := h.GetVault(vaultID)
	if vault == nil {
		log.Println("Tried to remove connection from vault that doesn't exist")
		return
	}
	empty := vault.RemoveConnection(playerUUID)
	if empty {
		h.RemoveVault(vaultID)
	}
}
