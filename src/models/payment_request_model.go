package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	PaymentStatusPending   = "pending"
	PaymentStatusConfirmed = "confirmed"
	PaymentStatusExpired   = "expired"
)

type PaymentRequest struct {
	ID              uint             `json:"id" gorm:"primaryKey"`
	UUID            uuid.UUID        `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();uniqueIndex"`
	UserID          uint             `json:"user_id"`
	User            *UserModel       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	BlockchainID    uint             `json:"blockchain_id"`
	Blockchain      *Blockchain      `json:"blockchain,omitempty" gorm:"foreignKey:BlockchainID"`
	RequestedCount  int              `json:"requested_count"`
	ExpectedAmount  *float64         `json:"expected_amount"`
	Status          string           `json:"status" gorm:"default:pending"`
	ExpiresAt       *time.Time       `json:"expires_at"`
	ConfirmedAt     *time.Time       `json:"confirmed_at"`
	WalletAddresses []*WalletAddress `json:"wallet_addresses,omitempty" gorm:"foreignKey:PaymentRequestID"`
	Deposits        []*Deposits      `json:"deposits,omitempty" gorm:"foreignKey:PaymentRequestID"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       gorm.DeletedAt   `json:"deleted_at,omitempty" gorm:"index"`
}

func (PaymentRequest) TableName() string {
	return "payment_requests"
}
