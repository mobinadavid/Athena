package errs

import (
	"errors"
)

type SeverityLevel uint16

const (
	SeverityLevelInfo SeverityLevel = iota
	SeverityLevelWarn
	SeverityLevelError
)

type AppError struct {
	Err           error
	SeverityLevel SeverityLevel
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(err error, level SeverityLevel) *AppError {
	return &AppError{
		Err:           err,
		SeverityLevel: level,
	}
}

func (e *AppError) WithLevel(level SeverityLevel) *AppError {
	return &AppError{
		Err:           e.Err,
		SeverityLevel: level,
	}
}

// common
var (
	SomeThingWentWrong           = errors.New("something-went-wrong")
	InvalidUuid                  = errors.New("uuid-is-invalid")
	RecordNotFound               = errors.New("record-not-found")
	CantChangeSuperAdmin         = errors.New("you-can-not-change-super-admin-title")
	InvalidMethodForNotification = errors.New("not-supported-method-to-send-notification")
	BankNotExists                = errors.New("bank-not-exists")
	NotImplementedYet            = errors.New("feature-not-implemented-yet")
	InvalidInputForThisRequest   = errors.New("invalid-input-for-this-request")
)

// file
var (
	InvalidFormData               = NewAppError(errors.New("invalid-form-data"), SeverityLevelError)
	CannotOpenFile                = NewAppError(errors.New("cannot-open-file"), SeverityLevelError)
	CDNServiceIsTemporaryDown     = NewAppError(errors.New("cdn-service-is-temporary-down"), SeverityLevelError)
	ErrTooManySimultaneousUploads = NewAppError(errors.New("too-many-simultaneous-uploads"), SeverityLevelWarn)
	ErrFileSizeTooLarge           = NewAppError(errors.New("file-size-to-large"), SeverityLevelWarn)
	ErrInvalidFileExtension       = NewAppError(errors.New("invalid-file-extension"), SeverityLevelWarn)
)

// authenticate
var (
	ErrAuthenticationFailed   = errors.New("auth-failed")
	RegisterFailed            = errors.New("register-failed")
	RecoverPasswordFailed     = errors.New("recover-password-failed")
	ChangePasswordFailed      = errors.New("change-password-failed")
	PasswordNotMatch          = errors.New("invalid-password-match")
	PasswordShouldBeNew       = errors.New("new-password-is-equal-to-current-password")
	ErrDeactivatedAdmin       = errors.New("admin-is-not-active")
	ErrDeactivatedUser        = errors.New("user-is-not-active")
	ErrRegisterTimeOut        = errors.New("register-time-out")
	ErrRecoverPasswordTimeOut = errors.New("recover-password-time-out")
	ErrChangePasswordTimeOut  = errors.New("change-password-time-out")
	ErrLoginTimeOut           = errors.New("login-time-out")
	ErrAdminHasTwoFactorAuth     = errors.New("user-has-two-factor-auth")
	ErrTwoFactorRequired         = errors.New("two-factor-authentication-required")
	ErrTwoFactorAlreadyEnabled   = errors.New("two-factor-authentication-already-enabled")
	ErrTwoFactorSecretMissing    = errors.New("two-factor-secret-is-missing")
	ErrSameUserActivation        = errors.New("user-has-the-same-activation")
	Invalid2FaCode               = errors.New("invalid-two-factor-code")
	NotActiveTwoFactor           = errors.New("two-factor-authentication-is-not-active")
	InvalidRecoveryCode          = errors.New("invalid-recovery-code")
	ErrTwoFactorChallengeMissing = errors.New("two-factor-code-or-recovery-code-is-required")
	OTPIsNotValid             = errors.New("otp-is-not-valid")
	DuplicateUsername         = errors.New("duplicate-user")
	InvalidRevokeTokenFlow    = errors.New("invalid_revoke_token_flow")
)

// authorization
var (
	DuplicatePermissionGroupName = NewAppError(errors.New("duplicate-permission-group-name"), SeverityLevelWarn)
	DuplicateRoleName            = NewAppError(errors.New("duplicate-role-name"), SeverityLevelWarn)
	CantDeletePermissionGroup    = NewAppError(errors.New("cant-delete-permission-group"), SeverityLevelInfo)
	CantDeleteThisRole           = NewAppError(errors.New("cant-delete-this-role"), SeverityLevelInfo)
)

// token
var (
	RefreshTokenMissing     = errors.New("refresh-token-is-missing")
	ErrInvalidRefreshToken  = errors.New("invalid-refresh-token")
	ErrTokenExpired         = errors.New("token-expired")
	ErrInvalidToken         = errors.New("invalid-token")
	ErrInvalidSigningMethod = errors.New("unexpected-signing-method")
)

// otp
var (
	ErrOTPRequired   = errors.New("auth-otp-sent")
	ErrOTPInvalid    = errors.New("auth-otp-invalid")
	ErrAuthOTPExists = errors.New("auth-otp-exists")
	FailedToSendOTP  = errors.New("failed-to-send-top")
)

// captcha
var (
	CaptchaServiceIsTemporaryDown = errors.New("captcha-service-is-temporary-down")
	CaptchaIsNotValid             = errors.New("captcha-is-not-valid")
)

// rate limiter
var (
	TooManyRequest = errors.New("too-many-request")
)

// user
var (
	UserNotFound = errors.New("user-not-found")
)
