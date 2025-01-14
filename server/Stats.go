package server

import (
	"github.com/NodiumHosting/VaultMapperSyncServer/models"
	pb "github.com/NodiumHosting/VaultMapperSyncServer/proto"
	"github.com/disgoorg/json"
	"net/http"
	"os"
	"sync"
	"time"
)

func getActiveVaults() int {
	var i int
	HUB.Vaults.Range(func(k, v interface{}) bool {
		i++
		return true
	})
	return i
}

func getActiveConnections() int {
	var i int
	HUB.Vaults.Range(func(k, v interface{}) bool {
		vault := v.(*Vault)
		vault.Connections.Range(func(k, v interface{}) bool {
			i++
			return true
		})
		return true
	})
	return i
}

func getActiveCells() int {
	var i int
	HUB.Vaults.Range(func(k, v interface{}) bool {
		vault := v.(*Vault)
		vault.Cells.Range(func(k, v interface{}) bool {
			i++
			return true
		})
		return true
	})
	return i
}

func getActiveRooms() int {
	var i int
	HUB.Vaults.Range(func(k, v interface{}) bool {
		vault := v.(*Vault)
		vault.Cells.Range(func(k, v interface{}) bool {
			cell := v.(*pb.VaultCell)
			if cell.CellType == pb.CellType_CELLTYPE_ROOM {
				i++
			}
			return true
		})
		return true
	})
	return i
}

func GetTotalPlayerCount() (int64, error) {
	var count int64
	err := DB.Model(&models.Player{}).Count(&count).Error
	return count, err
}

func GetPlayerCountInVault(vaultID string) (int64, error) {
	var count int64
	err := DB.Model(&models.PlayerVault{}).Where("vault_id = ?", vaultID).Count(&count).Error
	return count, err
}

func GetTotalVaults() (int64, error) {
	var count int64
	err := DB.Model(&models.Vault{}).Count(&count).Error
	return count, err
}

func GetTotalRooms() (int64, error) {
	var count int64
	err := DB.Model(&models.VaultCell{}).Where("cell_type = ?", pb.CellType_CELLTYPE_ROOM).Count(&count).Error
	return count, err
}

// GetLargestVault returns the vault with the most cells in it
/*func GetLargestVault() (int64, error) {
	var count int64
	err := DB.Model(&models.VaultCell{}).Select("vault_id").Group("vault_id").Count(&count).Error
	return count, err
}*/

func GetLargestVault() (int64, error) {
	var result struct {
		VaultID string
		Count   int64
	}
	err := DB.Model(&models.VaultCell{}).
		Select("vault_id, COUNT(*) as count").
		Group("vault_id").
		Order("count DESC").
		Limit(1).
		Scan(&result).Error
	if err != nil {
		return 0, err
	}
	return result.Count, nil
}

/*func GetBiggestParty() (int64, error) {
	var count int64
	err := DB.Model(&models.PlayerVault{}).Select("player_uuid").Group("vault_id").Count(&count).Error
	return count, err
}*/

func GetBiggestParty() (int64, error) {
	var result struct {
		VaultID string
		Count   int64
	}
	err := DB.Model(&models.PlayerVault{}).
		Select("vault_id, COUNT(player_uuid) as count").
		Group("vault_id").
		Order("count DESC").
		Limit(1).
		Scan(&result).Error
	if err != nil {
		return 0, err
	}
	return result.Count, nil
}

// GetTotalRoomsBasic returns the total number of basic rooms in db
func GetTotalRoomsBasic() (int64, error) {
	var count int64
	err := DB.Model(&models.VaultCell{}).Where("room_type = ? & cell_type = ?", pb.RoomType_ROOMTYPE_BASIC, pb.CellType_CELLTYPE_ROOM).Count(&count).Error
	return count, err
}

func GetTotalRoomsOre() (int64, error) {
	var count int64
	err := DB.Model(&models.VaultCell{}).Where("room_type = ?", pb.RoomType_ROOMTYPE_ORE).Count(&count).Error
	return count, err
}

func GetTotalRoomsChallenge() (int64, error) {
	var count int64
	err := DB.Model(&models.VaultCell{}).Where("room_type = ?", pb.RoomType_ROOMTYPE_CHALLENGE).Count(&count).Error
	return count, err
}

func GetTotalRoomsOmega() (int64, error) {
	var count int64
	err := DB.Model(&models.VaultCell{}).Where("room_type = ?", pb.RoomType_ROOMTYPE_OMEGA).Count(&count).Error
	return count, err
}

// GetActivity returns a list of vaults with the players inside them
func GetActivity() map[string][]string {
	vaults := make(map[string][]string)
	HUB.Vaults.Range(func(k, v interface{}) bool {
		vaultID := k.(string)
		vault := v.(*Vault)
		var players []string
		vault.Connections.Range(func(k, v interface{}) bool {
			players = append(players, k.(string))
			return true
		})
		vaults[vaultID+" : "+vault.ViewerCode] = players
		return true
	})

	return vaults
}

var (
	stats           = make(map[string]interface{})
	cacheExpiration time.Time
	cacheMutex      sync.Mutex
	mainUpdateIndex int                // some stats are quite heavy to update so they'll be updated slower
	cacheDuration   = 10 * time.Second // Cache duration
)

func updateStatsCache() {
	//log.Println("Updating stats cache")
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	//stats := make(map[string]interface{})

	uniquePlayerCount, err := GetTotalPlayerCount()
	if err == nil {
		stats["unique_player_count"] = uniquePlayerCount
	}

	activeVaults := getActiveVaults()
	stats["active_vaults"] = activeVaults

	activeConnections := getActiveConnections()
	stats["active_connections"] = activeConnections

	activeCells := getActiveCells()
	stats["active_cells"] = activeCells

	activeRooms := getActiveRooms()
	stats["active_rooms"] = activeRooms

	// heavy updates from here
	switch mainUpdateIndex {
	case 0:
		biggestParty, err := GetBiggestParty()
		if err == nil {
			stats["biggest_party"] = biggestParty
		}

	case 1:
		totalVaults, err := GetTotalVaults()
		if err == nil {
			stats["total_vaults"] = totalVaults
		}

	case 2:
		totalRooms, err := GetTotalRooms()
		if err == nil {
			stats["total_rooms"] = totalRooms
		}

	case 3:
		totalRoomsBasic, err := GetTotalRoomsBasic()
		if err == nil {
			stats["total_rooms_basic"] = totalRoomsBasic
		}

	case 4:
		totalRoomsOre, err := GetTotalRoomsOre()
		if err == nil {
			stats["total_rooms_ore"] = totalRoomsOre
		}

	case 5:
		totalRoomsChallenge, err := GetTotalRoomsChallenge()
		if err == nil {
			stats["total_rooms_challenge"] = totalRoomsChallenge
		}

	case 6:
		totalRoomsOmega, err := GetTotalRoomsOmega()
		if err == nil {
			stats["total_rooms_omega"] = totalRoomsOmega
		}

	case 7:
		largestVaultCount, err := GetLargestVault()
		if err == nil {
			stats["largest_vault"] = largestVaultCount
		}
	}
	mainUpdateIndex++
	mainUpdateIndex = mainUpdateIndex % 8

	stats["activity"] = GetActivity()

	//statsCache = stats
	cacheExpiration = time.Now().Add(cacheDuration)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token != os.Getenv("TOKEN") || token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cacheMutex.Lock()
	if time.Now().After(cacheExpiration) {
		cacheMutex.Unlock()
		updateStatsCache()
		cacheMutex.Lock()
	}
	st := stats
	cacheMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(st); err != nil {
		http.Error(w, "Failed to encode stats", http.StatusInternalServerError)
	}
}
