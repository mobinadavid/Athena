package models

import (
	"gorm.io/gorm"
	"time"
)

type Deposits struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	TransactionHash string         `gorm:"not null" json:"transaction_hash"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName sets the table name of the model
func (Deposits) TableName() string {
	return "deposits"
}
