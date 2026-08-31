package AuthenticationRequests

type IssueTokenViaPasswordRequest struct {
	NationalIdentityCode string `json:"national_identity_code" validate:"required,iranian-national-identity-code"`
	Password             string `json:"password" validate:"required,max=255"`
}
