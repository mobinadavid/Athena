package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Blockchain struct {
	ID                  uint                  `gorm:"primaryKey" json:"id"`
	UUID                *uuid.UUID            `gorm:"default:uuid_generate_v4()" json:"uuid"`
	NativeAsset         string                `gorm:"not null" json:"native_asset"`
	Title               datatypes.JSON        `gorm:"type:json" json:"title"`
	Name                string                `gorm:"not null;uniqueIndex" json:"name"`
	IsActive            *bool                 `gorm:"type:bool;default:true" json:"is_active"`
	CreatedAt           time.Time             `json:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at"`
	DeletedAt           gorm.DeletedAt        `gorm:"index" json:"deleted_at"`
	WalletAddresses     []WalletAddress       `gorm:"foreignKey:BlockchainID" json:"wallet_addresses"`
	BlockchainExplorers []*BlockchainExplorer `gorm:"many2many:blockchain_explorer_mappings;" json:"blockchain_explorers"`
}

// TableName sets the table name of the model
func (Blockchain) TableName() string {
	return "blockchains"
}
