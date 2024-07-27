package request

type GetWalletAddress struct {
	Blockchain string `json:"blockchain" validate:"required,max=255"`
	Number     int    `json:"number" validate:"required,min=1"`
}
