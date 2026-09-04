package requests

type AllocateWalletAddress struct {
	Blockchain     string   `json:"blockchain" validate:"required,max=255"`
	Count          int      `json:"count" validate:"required,min=1,max=20"`
	ExpectedAmount *float64 `json:"expected_amount" validate:"omitempty,gt=0"`
}
