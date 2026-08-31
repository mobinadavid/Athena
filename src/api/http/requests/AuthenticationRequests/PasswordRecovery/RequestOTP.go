package PasswordRecovery

type RequestOtp struct {
	NationalIdentityCode string `json:"national_identity_code" validate:"omitempty,iranian-national-identity-code"`
	NationalCompanyId    string `json:"national_company_id" validate:"omitempty,iranian-company-national-id"`
	Mobile               string `json:"mobile,required" validate:"required,iranian-mobile"`
}
