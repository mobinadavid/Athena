package Authentication

type LoginRequest struct {
	Username             string `json:"username" validate:"omitempty,max=255"`
	NationalIdentityCode string `json:"national_identity_code" validate:"omitempty,iranian-national-identity-code"`
	NationalCompanyId    string `json:"national_company_id" validate:"omitempty,iranian-company-national-id"`
	Password             string `json:"password" validate:"required,max=255"`
}

func (lr *LoginRequest) OnlyOneGiven() bool {
	count := 0
	if lr.Username != "" {
		count++
	}
	if lr.NationalIdentityCode != "" {
		count++
	}
	if lr.NationalCompanyId != "" {
		count++
	}

	if count != 1 {
		return false
	}
	return true
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

type ForceChangePasswordUserRequest struct {
	LoginKey           string `json:"login_key" validate:"omitempty"`
	NewPassword        string `json:"new_password" validate:"required,max=255,is-strong-password"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"required,max=255,eqfield=NewPassword"`
}

type LoginATMSendOtpRequest struct {
	NationalIdentityCode string `json:"national_identity_code" validate:"omitempty,iranian-national-identity-code"`
	NationalCompanyId    string `json:"national_company_id" validate:"omitempty,iranian-company-national-id"`
	Mobile               string `json:"mobile,required" validate:"required,iranian-mobile"`
}

type LoginScenario string

const (
	LoginScenarioRegister        LoginScenario = "register"
	LoginScenarioCompleteProfile LoginScenario = "complete_profile"
	LoginScenarioNormalLogin     LoginScenario = "normal_login"
)

type LoginATMState struct {
	Request  LoginATMSendOtpRequest `json:"request"`
	Scenario LoginScenario          `json:"scenario"`
}

type LoginATMRequest struct {
	Mobile               string `json:"mobile" validate:"required,iranian-mobile"`
	NationalIdentityCode string `json:"national_identity_code" validate:"omitempty,iranian-national-identity-code"`
	NationalCompanyId    string `json:"national_company_id" validate:"omitempty,iranian-company-national-id"`
}

// OnlyOneNationalIdGiven ensures exactly one of NationalIdentityCode / NationalCompanyId is set.
func (r *LoginATMRequest) OnlyOneNationalIdGiven() bool {
	if r.NationalIdentityCode != "" && r.NationalCompanyId != "" {
		return false
	}
	if r.NationalIdentityCode == "" && r.NationalCompanyId == "" {
		return false
	}
	return true
}

// NationalCode returns whichever national identifier was supplied (person or company).
func (r *LoginATMRequest) NationalCode() string {
	if r.NationalCompanyId != "" {
		return r.NationalCompanyId
	}
	return r.NationalIdentityCode
}
