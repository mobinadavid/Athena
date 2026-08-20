package Authentication

type RegisterRequest struct {
	NationalIdentityCode string `json:"national_identity_code" validate:"omitempty,iranian-national-identity-code"`
	Mobile               string `json:"mobile" validate:"required,iranian-mobile"`
	Password             string `json:"password" validate:"required,max=255,is-strong-password"`
	RePassword           string `json:"re_password" validate:"required,max=255,is-strong-password"`
}

type VerifyRegisterOTP struct {
	RegisterKey string `json:"register_key" validate:"omitempty"`
	OTP         string `json:"otp" validate:"required"`
}

type ResendRegisterOTP struct {
	RegisterKey string `json:"register_key" validate:"omitempty"`
}
