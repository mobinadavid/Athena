package errs

const (
	UserSuccessfullyLoginErrorCode       = 10000
	ToManyRequestErrorCode               = 10002
	OTPAlreadyExistErrorCode             = 10003
	RegisterTimeOutErrorCode             = 10004
	LoginTimeOutErrorCode                = 10004
	RecoverPasswordTimeOutErrorCode      = 10004
	UserHasTwoFactorAuthErrorCode        = 10006
	AdminHasTwoFactorAuthErrorCode       = 10006
	AdminMustSetupTwoFactorAuthErrorCode = 10008
	DuplicateKey                         = "23505"
)
