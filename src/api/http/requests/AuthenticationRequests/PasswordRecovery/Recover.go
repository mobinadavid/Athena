package PasswordRecovery

type RecoverRequest struct {
	Username                string `json:"username" validate:"required,max=255"`
	OTP                     string `json:"otp" validate:"required,number,len=5"`
	NewPassword             string `json:"new_password" validate:"required,max=255,is-strong-password"`
	NewPasswordConfirmation string `json:"new_password_confirmation" validate:"required,max=255,is-strong-password,eqfield=NewPassword"`
}
