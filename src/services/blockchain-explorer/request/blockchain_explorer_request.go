package request

type CreateBlockchainExplorerRequest struct {
	BaseUrl                string   `json:"base_url" validate:"required,max=255"`
	Apikey                 string   `json:"api_key" validate:"required,max=255"`
	BlockchainExplorerName string   ` json:"blockchain_explorer_name" validate:""`
	IsActive               *bool    `json:"is_active" validate:"required,boolean"`
	Blockchains            []string ` json:"blockchains" validate:""`
	IsDefault              *bool    `json:"is_default" validate:"boolean"`
}
