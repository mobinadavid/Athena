package AuthenticationRequests

type VerifyOTP struct {
	Mobile string `json:"mobile" validate:"required,iranian-mobile"`
	OTP    string `json:"otp" validate:"required,number,len=5"`
}
