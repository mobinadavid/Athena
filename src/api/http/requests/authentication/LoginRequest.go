package Authentication

type LoginRequest struct {
	Username             string `json:"username" validate:"omitempty,max=255"`
	NationalIdentityCode string `json:"national_identity_code" validate:"omitempty,iranian-national-identity-code"`
	Password             string `json:"password" validate:"required,max=255"`
}

type VerifyLoginOTP struct {
	LoginKey string `json:"login_key" validate:"omitempty"`
	OTP      string `json:"otp" validate:"required"`
}

type ResendLoginOTP struct {
	LoginKey string `json:"login_key" validate:"omitempty"`
}
type LoginViaOTPSendOtpRequest struct {
	NationalIdentityCode string `json:"national_identity_code" validate:"omitempty,iranian-national-identity-code"`
	NationalCompanyId    string `json:"national_company_id" validate:"omitempty,iranian-company-national-id"`
}

type VerifyTwoFa struct {
	LoginKey  string `json:"login_key" validate:"omitempty"`
	TwoFaCode string `json:"two_fa_code" validate:"required"`
}

type ChangePasswordRequest struct {
	NewPassword             string `json:"new_password" validate:"required,max=255,is-strong-password"`
	NewPasswordConfirmation string `json:"new_password_confirmation" validate:"required,max=255,eqfield=NewPassword"`
	LoginKey                string `json:"login_key" validate:"omitempty"`
}
