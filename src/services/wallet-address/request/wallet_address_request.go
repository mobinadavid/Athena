package request

import (
	"gorm.io/datatypes"
)

type CreateWalletAddressRequest struct {
	WalletAddress string         `json:"wallet_address" validate:"required,max=255"`
	Title         datatypes.JSON ` json:"title" validate:""`
	IsActive      *bool          `json:"is_active" validate:"required,boolean"`
	BlockchainId  uint           `json:"blockchain_id" validate:"required,number"`
}
