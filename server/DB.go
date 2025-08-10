package server

import (
	"errors"
	"github.com/NodiumHosting/VaultMapperSyncServer/models"
	"github.com/go-co-op/gocron/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

var DB *gorm.DB
var DBScheduler gocron.Scheduler

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("file:vaults.db?_journal_mode=WAL"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	err1 := DB.AutoMigrate(&models.Vault{}, &models.VaultCell{}, &models.PlayerVault{}, &models.Player{})
	if err1 != nil {
		return
	}
}

// CleanDB is used to clean the db of any old vaults, cells and players
// will be triggered every x using gocron
func CleanDB() {
	log.Println("Cleaning database")

	// add timing, display the time it took after commit done
	var startTime = time.Now()

	query := `
DELETE FROM vault_cells 
WHERE vault_id IN (
    SELECT vault_id FROM vaults 
    WHERE created_at < datetime('now', '-30 days')
);
DELETE FROM player_vaults 
WHERE vault_id IN (
    SELECT vault_id FROM vaults 
    WHERE created_at < datetime('now', '-30 days')
);
DELETE FROM vaults 
WHERE created_at < datetime('now', '-30 days');
`
	if err := DB.Exec(query).Error; err != nil {
		log.Println("Error cleaning database: ", err)

	}

	// print time it took to clean db
	log.Printf("Time taken to clean database: %s\n", time.Since(startTime))
}

func StartCron() {
	var err error
	DBScheduler, err = gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}

	_, err = DBScheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0))),
		gocron.NewTask(
			func() {
				CleanDB()
			},
		),
	)
	if err != nil {
		log.Fatalf("Failed to create job: %v", err)
	}
	DBScheduler.Start()
	log.Println("Cron started")
}

func StopCron() {
	_ = DBScheduler.Shutdown()
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
