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
	DB, err := gorm.Open(sqlite.Open("file:vaults.db?_journal_mode=WAL"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	err1 := DB.AutoMigrate(&models.Vault{}, &models.VaultCell{})
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
