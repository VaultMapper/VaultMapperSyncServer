package server

import (
	"errors"
	"github.com/NodiumHosting/VaultMapperSyncServer/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("file:vaults.db?_journal_mode=WAL"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	err1 := DB.AutoMigrate(&models.Vault{}, &models.VaultCell{}, &models.PlayerVault{}, &models.Player{})
	if err1 != nil {
		return
	}
}

// SaveVault updates the vault record in the database, creates a new vault record if needed
//
// This function is incredibly expensive for larger vaults and should not be used unless absolutely necessary
func SaveVault(vault *models.Vault) error {
	var existingVault models.Vault
	result := DB.First(&existingVault, "vault_id = ?", vault.VaultID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		if err := DB.Create(vault).Error; err != nil {
			return err
		}
	} else {
		if err := DB.Where("vault_id = ?", vault.VaultID).Delete(&models.VaultCell{}).Error; err != nil {
			return err
		}
		existingVault.Cells = vault.Cells
		if err := DB.Save(&existingVault).Error; err != nil {
			return err
		}
	}

	return nil
}

func AddPlayer(uuid string) error {
	player := models.Player{UUID: uuid}
	return DB.Create(&player).Error
}

// AddPlayerToVault adds a player to a vault inside the db for stats keeping purposes, if needed creates the player record
func AddPlayerToVault(playerUUID string, vaultID string) error {
	res := DB.Find(&models.Player{}, "uuid = ?", playerUUID)
	if res.RowsAffected == 0 {
		err := AddPlayer(playerUUID)
		if err != nil {
			return err
		}
	}
	var playerVault models.PlayerVault
	result := DB.First(&playerVault, "player_uuid = ? AND vault_id = ?", playerUUID, vaultID)
	if result.RowsAffected == 0 {
		playerVault = models.PlayerVault{PlayerUUID: playerUUID, VaultID: vaultID}
		return DB.Create(&playerVault).Error
	}

	return nil
}
