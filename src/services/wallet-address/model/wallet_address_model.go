package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type WalletAddress struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UUID          uuid.UUID      `gorm:"default:uuid_generate_v4()" json:"uuid"`
	Title         datatypes.JSON `gorm:"type:json" json:"title"`
	WalletAddress string         `gorm:"unique" json:"wallet_address"`
	WebhookURL    string         `gorm:"unique" json:"webhook_url"`
	IsActive      *bool          `gorm:"type:bool;default:true" json:"is_active"`
	AllocatedAt   time.Time      `gorm:"default:null" json:"allocated_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	BlockchainID  uint           `gorm:"not null" json:"blockchain_id"`
}

// TableName sets the table name of the model
func (WalletAddress) TableName() string {
	return "wallet_addresses"
}
