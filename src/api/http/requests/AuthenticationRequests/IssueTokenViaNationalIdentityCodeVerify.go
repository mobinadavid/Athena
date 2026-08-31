package AuthenticationRequests

type IssueTokenViaNationalIdentityCodeVerify struct {
	NationalIdentityCode string `json:"national_identity_code" validate:"required,iranian-national-identity-code"`
	OTP                  string `json:"otp" validate:"required,number,len=5"`
}
