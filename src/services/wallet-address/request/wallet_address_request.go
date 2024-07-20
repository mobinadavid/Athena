package request

type CreateWalletAddressRequest struct {
	WalletAddress     string `json:"wallet_address" validate:"required,max=255"`
	WalletAddressName string ` json:"wallet_address_name" validate:""`
	IsActive          *bool  `json:"is_active" validate:"required,boolean"`
	BlockchainName    string `json:"blockchain_name" validate:"required,max=255"`
}
