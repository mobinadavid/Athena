package Authentication

import "github.com/google/uuid"

type RegisterRequest struct {
	NationalIdentityCode string `json:"national_identity_code" validate:"omitempty,iranian-national-identity-code"`
	NationalCompanyId    string `json:"national_company_id" validate:"omitempty,iranian-company-national-id"`
	Mobile               string `json:"mobile" validate:"required,iranian-mobile"`
	Password             string `json:"password" validate:"required,max=255,is-strong-password"`
	RePassword           string `json:"re_password" validate:"required,max=255,is-strong-password"`
	Captcha              struct {
		Uuid uuid.UUID `json:"uuid" validate:"required,uuid"`
		Text string    `json:"text" validate:"max=6,min=1"`
	} `json:"captcha" validate:"required"`
}

type VerifyRegisterOTP struct {
	RegisterKey string `json:"register_key" validate:"omitempty"`
	OTP         string `json:"otp" validate:"required"`
}

type ResendRegisterOTP struct {
	RegisterKey string `json:"register_key" validate:"omitempty"`
}
