package requests

type CreateBlockchainExplorerRequest struct {
	BaseUrl     string   `json:"base_url" validate:"required,max=255"`
	Name        string   ` json:"name" validate:""`
	IsActive    *bool    `json:"is_active" validate:"required,boolean"`
	Blockchains []string ` json:"blockchains" validate:""`
	IsDefault   bool     `json:"is_default" validate:"boolean"`
}
