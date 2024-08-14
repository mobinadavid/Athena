package requests

type CreatePaymentRequest struct {
	Amount      float64 `json:"amount" validate:"required,min=0"`
	Mobile      string  `json:"mobile" validate:"omitempty,e164"`
	Driver      string  `json:"driver"`
	CallBackUrl string  `json:"callback_url" validate:"required"`
}
