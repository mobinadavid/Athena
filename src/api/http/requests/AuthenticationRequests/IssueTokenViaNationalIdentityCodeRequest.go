package AuthenticationRequests

type IssueTokenViaNationalIdentityCodeRequest struct {
	NationalIdentityCode string `json:"national_identity_code" validate:"required,iranian-national-identity-code"`
}
