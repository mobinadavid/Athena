package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	DepositStatusDetected  = "detected"
	DepositStatusConfirmed = "confirmed"
)

type Deposits struct {
	ID               uint            `json:"id" gorm:"primaryKey"`
	UUID             uuid.UUID       `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();uniqueIndex"`
	TransactionHash  string          `json:"transaction_hash"`
	UserID           *uint           `json:"user_id"`
	User             *UserModel      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	WalletAddressID  *uint           `json:"wallet_address_id"`
	WalletAddress    *WalletAddress  `json:"wallet_address,omitempty" gorm:"foreignKey:WalletAddressID"`
	PaymentRequestID *uint           `json:"payment_request_id"`
	PaymentRequest   *PaymentRequest `json:"payment_request,omitempty" gorm:"foreignKey:PaymentRequestID"`
	BlockchainID     *uint           `json:"blockchain_id"`
	Blockchain       *Blockchain     `json:"blockchain,omitempty" gorm:"foreignKey:BlockchainID"`
	FromAddress      string          `json:"from_address"`
	ToAddress        string          `json:"to_address"`
	Amount           float64         `json:"amount"`
	Fee              float64         `json:"fee"`
	Confirmations    int             `json:"confirmations"`
	Status           string          `json:"status" gorm:"default:confirmed"`
	NotifiedAt       *time.Time      `json:"notified_at"`
	BlockNumber      int64           `json:"block_number"`
	PaidAt           *time.Time      `json:"paid_at"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"index"`
}

func (Deposits) TableName() string {
	return "deposits"
}
