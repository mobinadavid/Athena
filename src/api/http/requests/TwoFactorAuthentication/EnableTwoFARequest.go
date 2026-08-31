package TwoFactorAuthentication

type VerifyTotpRequest struct {
	Totp string `json:"totp" validate:"required"`
}
