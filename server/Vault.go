package server

import (
	"errors"
	"fmt"
	"github.com/NodiumHosting/VaultMapperSyncServer/models"
	pb "github.com/NodiumHosting/VaultMapperSyncServer/proto"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"log"
	"sync"
)

// Vault is a helper struct that is used to store connections to vault and other runtime data
//
// If other data is added, there should be a Mutex created with it to make sure writing is thread-safe
type Vault struct {
	UUID        string
	Connections sync.Map // stores a map of current connections inside the vault, key is uuid, value is Connection
	Viewers     sync.Map // stores a map of viewers, same as Connections but separated
	ViewerCode  string   // code used to identify Vault for viewing purposes. Mutex shouldn't be needed for ViewerCode as the value is only read after it's creation which cannot result in race conditions
	Cells       sync.Map // stores a map of cells inside the vault, key is x,z, value is pb.VaultCell
}

// AddViewer adds the connection to Vault structure and starts up the WritePump
//
// ok is false if the connection already exists, else false
func (v *Vault) AddViewer(UUID string, conn *websocket.Conn) bool {
	_, ok := v.Viewers.Load(UUID)
	if ok {
		log.Println("Tried to add connection but it already exists")
		return false // connection already exists
	}
	c := &Connection{ // create connection
		uuid: UUID,
		conn: conn,
		Send: make(chan []byte, 256), // buffered channel of 256 bytes
	}

	v.Viewers.Store(UUID, c) // store the connection inside vault
	go c.WritePump()         // Start the write pump
	return true
}

// RemoveViewer removes the connection from Vault structure and closes the Send channel
//
// return true if the Vault is empty after connection removal, false otherwise
//
// Send channel needs to be closed for WritePump to exit properly!
func (v *Vault) RemoveViewer(viewerUUID string) bool {
	value, ok := v.Viewers.Load(viewerUUID)
	if !ok {
		return false
	}

	c := value.(*Connection)

	close(c.Send)         // close send channel
	err := c.conn.Close() // close connection
	if err != nil {
		log.Println("Error closing connection: ", err)
		v.Viewers.Delete(viewerUUID) // still try to remove the connection even after error
		return false
	}
	v.Viewers.Delete(viewerUUID) // remove connection

	// check if the vault is empty now
	isEmpty := true
	v.Viewers.Range(func(k, v interface{}) bool {
		isEmpty = false
		return false
	})
	if isEmpty {
		return true
	}

	return false
}

// ViewerCount allows easy access to the number of viewers in a vault
func (v *Vault) ViewerCount() int {
	cnt := 0
	v.Viewers.Range(func(k, v interface{}) bool {
		cnt += 1
		return true
	})
	return cnt
}

// ClearViewers safely clears all viewers including closing write pumps
func (v *Vault) ClearViewers() {
	log.Println("clearing viewers")
	v.Viewers.Range(func(k, val interface{}) bool {
		c := val.(*Connection)
		v.RemoveViewer(c.uuid)
		return true
	})
	log.Println("DONE")
}

// AddConnection adds the connection to Vault structure and starts up the WritePump
//
// Returns *Connection if operation was successful, else nil
func (v *Vault) AddConnection(playerUUID string, conn *websocket.Conn) *Connection {
	_, ok := v.Connections.Load(playerUUID)
	if ok {
		log.Println("Tried to add connection but it already exists")
		return nil // connection already exists
	}
	c := &Connection{ // create connection
		uuid: playerUUID,
		conn: conn,
		Send: make(chan []byte, 256), // buffered channel of 256 bytes
	}

	v.Connections.Store(playerUUID, c) // store the connection inside vault
	go c.WritePump()                   // Start the write pump
	return c
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
		log.Println("Error closing connection: ", err)
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

// AddOrReplaceCell adds a VaultCell to the Cells map
func (v *Vault) AddOrReplaceCell(cell *pb.VaultCell) {
	key := fmt.Sprintf("%d,%d", cell.X, cell.Z)
	v.Cells.Store(key, cell)

	// Add/Update in db
	// Convert the proto cell to a database model cell
	dbCell := &models.VaultCell{
		VaultID:   v.UUID,
		X:         cell.X,
		Z:         cell.Z,
		CellType:  int32(cell.CellType),
		RoomType:  int32(cell.RoomType),
		RoomName:  cell.RoomName,
		Explored:  cell.Explored,
		Inscribed: cell.Inscribed,
		Marked:    cell.Marked,
	}

	// Check if the cell already exists in the database
	var existingCell models.VaultCell
	result := DB.First(&existingCell, "vault_id = ? AND x = ? AND z = ?", v.UUID, cell.X, cell.Z)
	//log.Println("HELLO")
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("fuck off")
		log.Println(result.Error)
		log.Println("fuck off")
		return
	}

	if result.RowsAffected == 0 {
		// Cell does not exist, create a new one
		log.Println("Creating cell")
		if err := DB.Create(dbCell).Error; err != nil {
			log.Println(err)
		}
	} else {
		// Cell exists, update it
		log.Println("found, updating")
		existingCell.CellType = dbCell.CellType
		existingCell.RoomType = dbCell.RoomType
		existingCell.RoomName = dbCell.RoomName
		existingCell.Explored = dbCell.Explored
		existingCell.Inscribed = dbCell.Inscribed
		existingCell.Marked = dbCell.Marked
		if err := DB.Save(&existingCell).Error; err != nil {
			log.Println(err)
		}
	}
}

// RemoveCell removes a VaultCell from the Cells map
func (v *Vault) RemoveCell(x, z int) {
	key := fmt.Sprintf("%d,%d", x, z)
	v.Cells.Delete(key)
}

// GetCell retrieves a VaultCell from the Cells map
func (v *Vault) GetCell(x, z int) (*pb.VaultCell, bool) {
	key := fmt.Sprintf("%d,%d", x, z)
	value, ok := v.Cells.Load(key)
	if !ok {
		return nil, false
	}
	return value.(*pb.VaultCell), true
}

// IterateCells runs a provided function on each cell in the vault
func (v *Vault) IterateCells(f func(key string, cell *pb.VaultCell)) { // iterate over cells
	v.Cells.Range(func(k, v interface{}) bool {
		key := k.(string)
		cell := v.(*pb.VaultCell)
		f(key, cell)
		return true
	})
}

func (v *Vault) BroadcastToast(text string) {
	message := pb.Message{Type: pb.MessageType_TOAST, Toast: &pb.Toast{Message: text}}
	messageBuffer, err := proto.Marshal(&message)
	if err != nil {
		return
	}
	v.Connections.Range(func(k, v interface{}) bool {
		conn := v.(*Connection)
		conn.Send <- messageBuffer
		return true
	})
}
