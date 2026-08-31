package AdminRequests

type ChangePasswordSendOtpRequest struct {
	CurrentPassword         string `json:"current_password" validate:"required,max=255"`
	NewPassword             string `json:"new_password" validate:"required,max=255,is-strong-password"`
	NewPasswordConfirmation string `json:"new_password_confirmation" validate:"required,max=255,eqfield=NewPassword"`
	Mobile                  string `json:"mobile"`
}
type ChangePasswordVerifyOtpRequest struct {
	OTP string `json:"otp" validate:"required,max=255"`
	Key string `json:"key" validate:"omitempty,max=255"`
}

type ChangePasswordResendOtpRequest struct {
	Key string `json:"key" validate:"omitempty,max=255"`
}
