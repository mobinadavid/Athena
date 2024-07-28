package request

type GetWalletAddress struct {
	Blockchain string `json:"blockchain" validate:"required,max=255"`
	Count      int    `json:"count" validate:"required,min=1"`
}
