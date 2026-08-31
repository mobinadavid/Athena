package login

type LoginViaOtpRequest struct {
	NationalIdentityCode string `json:"national_identity_code" validate:"required,iranian-national-identity-code"`
	Otp                  string `json:"otp" validate:"required"`
}
