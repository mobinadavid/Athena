package model

import (
	"athena/src/services/blockchain/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type BlockchainExplorer struct {
	ID                     uint                `gorm:"primaryKey" json:"id"`
	UUID                   uuid.UUID           `gorm:"default:uuid_generate_v4()" json:"uuid"`
	BaseUrl                string              `gorm:"not null" json:"base_url"`
	BlockchainExplorerName string              `gorm:"not null;uniqueIndex" json:"blockchain-explorer-name"`
	IsActive               bool                `gorm:"type:bool;default:true" json:"is_active"`
	IsDefault              bool                `gorm:"type:bool;default:true" json:"is_default"`
	CreatedAt              time.Time           `json:"created_at"`
	UpdatedAt              time.Time           `json:"updated_at"`
	DeletedAt              gorm.DeletedAt      `gorm:"index" json:"deleted_at"`
	Blockchains            []*model.Blockchain `gorm:"many2many:blockchain_explorer_mappings;"`
}

// TableName sets the table name of the model
func (BlockchainExplorer) TableName() string {
	return "blockchain_explorers"
}
