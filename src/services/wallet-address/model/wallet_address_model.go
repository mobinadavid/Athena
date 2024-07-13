package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Blockchain struct {
	ID   uint      `gorm:"primaryKey" json:"id"`
	UUID uuid.UUID `gorm:"default:uuid_generate_v4()" json:"uuid"`
	//	Walet
	IsActive  *bool          `gorm:"type:bool;default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName sets the table name of the model
func (Blockchain) TableName() string {
	return "blockchains"
}
