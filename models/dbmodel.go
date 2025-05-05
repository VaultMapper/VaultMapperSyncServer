package models

import (
	"gorm.io/gorm"
)

type Vault struct {
	gorm.Model
	VaultID string      `gorm:"primaryKey"`
	Cells   []VaultCell `gorm:"foreignKey:VaultID"`
}

type VaultCell struct {
	gorm.Model
	VaultID   string `gorm:"index"`
	X         int32
	Z         int32
	CellType  int32
	RoomType  int32
	RoomName  string
	Explored  bool
	Inscribed bool
	Marked    bool
}

type Player struct {
	gorm.Model
	UUID string `gorm:"primaryKey"`
}

type PlayerVault struct {
	gorm.Model
	PlayerUUID string `gorm:"index"`
	VaultID    string `gorm:"index"`
}
