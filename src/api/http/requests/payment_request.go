package requests

type PaymentRequest struct {
	Amount float64 `json:"amount" validate:"required,min=0"`
	Mobile string  `json:"mobile" validate:"required,e164"`
}
