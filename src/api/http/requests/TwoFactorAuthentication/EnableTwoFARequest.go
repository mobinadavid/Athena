package TwoFactorAuthentication

type VerifyTotpRequest struct {
	Totp string `json:"totp" validate:"required"`
}

type DisableTwoFARequest struct {
	Totp         string `json:"totp" validate:"omitempty"`
	RecoveryCode string `json:"recovery_code" validate:"omitempty"`
}
