package PasswordRecovery

type RecoverRequest struct {
	OTP                     string `json:"otp" validate:"required,number,len=5"`
	NationalIdentityCode    string `json:"national_identity_code" validate:"omitempty,iranian-national-identity-code"`
	NationalCompanyId       string `json:"national_company_id" validate:"omitempty,iranian-company-national-id"`
	NewPassword             string `json:"new_password" validate:"required,max=255,is-strong-password"`
	NewPasswordConfirmation string `json:"new_password_confirmation" validate:"required,max=255,is-strong-password,eqfield=NewPassword"`
}
