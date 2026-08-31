package AdminRequests

type CreateUpdateUserByAdminRequest struct {
	NationalIdentityCode string `json:"national_identity_code" validate:"omitempty,iranian-national-identity-code"`
	NationalCompanyId    string `json:"national_company_id" validate:"omitempty,iranian-company-national-id"`
	Mobile               string `json:"mobile" validate:"required,iranian-mobile"`
	Password             string `json:"password" validate:"required,max=255,is-strong-password"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,max=255,is-strong-password,eqfield=Password"`
}
