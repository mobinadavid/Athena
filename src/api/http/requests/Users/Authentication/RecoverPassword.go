package Authentication

type RecoverPasswordRequest struct {
	NationalIdentityCode string `json:"national_identity_code" validate:"omitempty,iranian-national-identity-code"`
	NationalCompanyId    string `json:"national_company_id" validate:"omitempty,iranian-company-national-id"`
	Mobile               string `json:"mobile" validate:"required,iranian-mobile"`
}

type VerifyRecoverPasswordOTP struct {
	RecoverPasswordKey string `json:"recover_password_key" validate:"omitempty"`
	OTP                string `json:"otp" validate:"required"`
}

type ResendRecoverPasswordOTP struct {
	RecoverPasswordKey string `json:"recover_password_key" validate:"omitempty"`
}

type SetRecoverPasswordRequest struct {
	RecoverPasswordKey string `json:"recover_password_key" validate:"omitempty"`
	Password           string `json:"password" validate:"required,max=255,is-strong-password"`
	RePassword         string `json:"re_password" validate:"required,max=255,is-strong-password"`
}
