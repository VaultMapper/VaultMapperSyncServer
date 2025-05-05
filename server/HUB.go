package server

import (
	"errors"
	"fmt"
	"github.com/NodiumHosting/VaultMapperSyncServer/dswh"
	"github.com/NodiumHosting/VaultMapperSyncServer/models"
	pb "github.com/NodiumHosting/VaultMapperSyncServer/proto"
	"github.com/NodiumHosting/VaultMapperSyncServer/util"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"log"
	"sync"
)

// Hub is a helper struct that keeps all currently running vaults inside
//
// should be thread-safe thanks to sync.Map use
type Hub struct {
	Vaults sync.Map
}

func GetVaultViewCode(vaultID string) string {
	vault := HUB.GetVault(vaultID)
	if vault == nil {
		return ""
	}
	return vault.ViewerCode
}

// GetOrCreateVault is a helper method that gets and optionally creates Vault inside Hub
func (h *Hub) GetOrCreateVault(vaultID string) *Vault {
	vault, loaded := h.Vaults.LoadOrStore(vaultID, &Vault{
		UUID: vaultID,
	})
	v := vault.(*Vault) // assert type of vault

	if !loaded { // if vault was created, check if there's stuff to load
		var dbVault models.Vault
		result := DB.Preload("Cells").First(&dbVault, "vault_id = ?", vaultID)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				log.Printf("Vault %s not found in database, creating new\n", vaultID)
				dbVault = models.Vault{
					VaultID: vaultID,
					Cells:   []models.VaultCell{},
				}
				if err := DB.Create(&dbVault).Error; err != nil {
					log.Println("Error creating new vault in DB: ", err)
					return v
				}
				v.ViewerCode = util.RandSeq(6)
				log.Println("New vault with viewer code: ", v.ViewerCode)
				return v
			}
			log.Println("Error querying DB: ", result.Error)
			return v
		}
		var cells []models.VaultCell
		DB.Find(&cells, "vault_id = ?", vaultID)
		// if vault was found, load cells
		for _, cell := range cells {
			pbCell := &pb.VaultCell{
				X:         cell.X,
				Z:         cell.Z,
				CellType:  pb.CellType(cell.CellType),
				RoomType:  pb.RoomType(cell.RoomType),
				RoomName:  cell.RoomName,
				Explored:  cell.Explored,
				Inscribed: cell.Inscribed,
				Marked:    cell.Marked,
			}
			key := fmt.Sprintf("%d,%d", pbCell.X, pbCell.Z)
			v.Cells.Store(key, pbCell)
			//log.Printf("Loaded cell from DB: %s\n", key)
		}
		log.Printf("Loaded vault from DB: %s\n", vaultID)
		v.ViewerCode = util.RandSeq(6)
		log.Println("New vault with viewer code: ", v.ViewerCode)
	}
	return v
}

func (h *Hub) GetVaultByCode(viewerCode string) *Vault {
	var vlt *Vault = nil
	h.Vaults.Range(func(key, v interface{}) bool {
		vault := v.(*Vault)
		if vault.ViewerCode == viewerCode {
			vlt = vault
			return false
		}

		return true
	})
	return vlt
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

	// Sending the vault image to discord
	vault := HUB.GetVault(vaultID)
	if vault == nil {
		log.Println("Tried to upload vault that doesn't exist")
		return // if vault doesn't exist, do nothing - this can happen when this is the first player joining a fresh vault
	}
	log.Println("Sending vault to discord")
	var cells []*pb.VaultCell
	vault.Cells.Range(func(key, val interface{}) bool {
		cells = append(cells, val.(*pb.VaultCell))
		//log.Println("appended cell")
		return true
	})
	go dswh.SendMap(cells, vaultID, DB)

	// make sure to clear and close all viewers
	vault.ClearViewers()

	h.Vaults.Delete(vaultID)
}

// AddConnectionToVault is a helper method that adds vault connection including vault creation if needed
//
// Returns *Connection or null if operation wasn't successful as first, *Vault that was accessed as second
func (h *Hub) AddConnectionToVault(vaultID string, playerUUID string, conn *websocket.Conn) (*Connection, *Vault) {
	vault := h.GetOrCreateVault(vaultID)
	c := vault.AddConnection(playerUUID, conn)

	return c, vault
}

// RemoveConnectionFromVault is a helper method that removes connection from vault including checks for empty vault
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

func (h *Hub) BroadcastToast(text string) {
	message := pb.Message{Type: pb.MessageType_TOAST, Toast: &pb.Toast{Message: text}}
	messageBuffer, err := proto.Marshal(&message)
	if err != nil {
		return
	}
	h.Vaults.Range(func(k, v interface{}) bool {
		vault := v.(*Vault)
		vault.Connections.Range(func(k, v interface{}) bool {
			conn := v.(*Connection)
			conn.Send <- messageBuffer
			return true
		})
		return true
	})
}

func (h *Hub) BroadcastToastInVault(vaultID string, text string) {
	message := pb.Message{Type: pb.MessageType_TOAST, Toast: &pb.Toast{Message: text}}
	messageBuffer, err := proto.Marshal(&message)
	if err != nil {
		return
	}
	vault := h.GetVault(vaultID)
	if vault == nil {
		log.Println("Tried to broadcast toast in vault that doesn't exist")
		return
	}
	vault.Connections.Range(func(k, v interface{}) bool {
		conn := v.(*Connection)
		conn.Send <- messageBuffer
		return true
	})
}
