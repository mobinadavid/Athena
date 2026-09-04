package PasswordRecovery

type RequestOtp struct {
	Username string `json:"username" validate:"required,max=255"`
	Mobile   string `json:"mobile" validate:"required,iranian-mobile"`
}
