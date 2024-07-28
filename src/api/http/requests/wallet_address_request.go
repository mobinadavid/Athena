package requests

type CreateWalletAddressRequest struct {
	WalletAddress string `json:"wallet_address" validate:"required,is-wallet-address"`
	Name          string ` json:"name" validate:""`
	IsActive      *bool  `json:"is_active" validate:"required,boolean"`
	Blockchain    string `json:"blockchain_name" validate:"required,max=255"`
}
