package authentication

import (
	"athena/src/api/errs"
	"athena/src/config"
	"athena/src/encrypt"
	"athena/src/hash"
	"athena/src/models"
	"athena/src/pkg/logger"
	"athena/src/pkg/utils"
	"athena/src/services"
	"context"
	"strconv"

	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
)

type TwoFaService struct {
	UserService services.IUserService
	OTPService  services.IOTPService
}

type ITwoFaService interface {
	GetSecretKey(ctx context.Context, userId uint, otp string) (string, string, error)
	VerifyCode(ctx context.Context, userId uint, code string) error
	VerifyCodeFirstTime(ctx context.Context, userId uint, code string) ([]string, error)
	Disable(ctx context.Context, userId uint) error
	CheckTwoFAEnable(ctx context.Context, userId uint) (bool, error)
	VerifyRecoveryCode(ctx context.Context, userId uint, recoveryCode string) error
	SendEnable2FAOTP(ctx context.Context) error
}

func (service *TwoFaService) SendEnable2FAOTP(ctx context.Context) error {
	userId, ok := ctx.Value("user_id").(uint)
	if !ok {
		logger.LogErrorWithFields("the user id is not set in context", service, models.UserRole, 0, errs.SomeThingWentWrong)
		return errs.SomeThingWentWrong
	}

	user, err := service.UserService.GetById(ctx, userId)
	if err != nil {
		logger.LogErrorWithFields("failed to get user by id", service, models.UserRole, userId, err)
		return errs.SomeThingWentWrong
	}

	err = service.OTPService.RequestOTP(ctx, user.Mobile)
	if err != nil {
		logger.LogErrorWithFields("failed to send create order otp", service, models.UserRole, userId, err)
		return err
	}

	return nil
}

func (service *TwoFaService) GetSecretKey(ctx context.Context, userId uint, otp string) (string, string, error) {
	user, err := service.UserService.GetById(ctx, userId)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get user by id", service, err)
		return "", "", err
	}

	otpIsValid, err := service.OTPService.VerifyOTP(ctx, user.Mobile, otp)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify otp", service, err)
		return "", "", errs.SomeThingWentWrong
	}
	if !otpIsValid {
		logger.LogErrorWithFieldsV2(ctx, "otp is not valid", service, nil,
			zap.String("otp", otp),
		)
		return "", "", errs.ErrOTPInvalid
	}

	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      config.GetInstance().Get("APP_NAME"),
		AccountName: user.Mobile,
	})
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to generate secret", service, err)
		return "", "", errs.SomeThingWentWrong
	}

	encryptedSecret, err := encrypt.GetInstance().Encrypt([]byte(secret.Secret()))
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to encrypt secret", service, err)
		return "", "", errs.SomeThingWentWrong
	}

	user.TotpSecret = encryptedSecret

	_, err = service.UserService.Update(ctx, user)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to update user", service, err)
		return "", "", errs.SomeThingWentWrong
	}
	return secret.URL(), secret.Secret(), nil
}

func (service *TwoFaService) VerifyCodeFirstTime(ctx context.Context, userId uint, code string) ([]string, error) {
	user, err := service.UserService.GetById(ctx, userId)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get user by id", service, err)
		return nil, errs.SomeThingWentWrong
	}

	decryptedSecret, err := encrypt.GetInstance().Decrypt(user.TotpSecret)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to decrypt secret", service, err)
		return nil, errs.SomeThingWentWrong
	}

	valid := totp.Validate(code, string(decryptedSecret))
	if !valid {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify secret code", service, errs.Invalid2FaCode)
		return nil, errs.Invalid2FaCode
	}

	//generate recovery code
	recoveryCodeNumber, err := strconv.Atoi(config.GetInstance().Get("RECOVERY_CODES_NUMBER"))
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to covert RECOVERY_CODES_NUMBER to int", service, err)
		recoveryCodeNumber = 3
	}

	recoveryCodes, err := utils.GenerateRandomCodes(recoveryCodeNumber)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to generate random code", service, err)
		return nil, errs.SomeThingWentWrong
	}

	var codes []string
	for _, recoveryCode := range recoveryCodes {
		hashedCode, err := hash.GetInstance().Generate([]byte(recoveryCode))
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to hash code", service, err)
			return nil, errs.SomeThingWentWrong
		}
		codes = append(codes, string(hashedCode))
	}

	//update user
	user.RecoveryCodes = codes
	user.TwoFaEnabled = true
	_, err = service.UserService.Update(ctx, user)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to update user", service, err)
		return nil, errs.SomeThingWentWrong
	}
	return recoveryCodes, nil
}

func (service *TwoFaService) Disable(ctx context.Context, userId uint) error {
	user, err := service.UserService.GetById(ctx, userId)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get user by id", service, err)
		return errs.SomeThingWentWrong
	}

	if user.TwoFaEnabled {
		user.TwoFaEnabled = false
		user.TotpSecret = nil
		user.RecoveryCodes = nil
		_, err = service.UserService.Update(ctx, user)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to update user", service, err)
			return errs.SomeThingWentWrong
		}
	}

	return nil
}

func (service *TwoFaService) VerifyCode(ctx context.Context, userId uint, code string) error {
	user, err := service.UserService.GetById(ctx, userId)
	if err != nil {
		logger.LogErrorWithFields("failed to get user by id", service, models.UserRole, userId, err)
		return errs.SomeThingWentWrong
	}

	decryptedSecret, err := encrypt.GetInstance().Decrypt(user.TotpSecret)
	if err != nil {
		logger.LogErrorWithFields("failed to decrypt secret", service, models.UserRole, userId, err)
		return errs.SomeThingWentWrong
	}

	valid := totp.Validate(code, string(decryptedSecret))
	if !valid {
		logger.LogErrorWithFields("failed to verify secret code", service, models.UserRole, userId, errs.Invalid2FaCode)
		return errs.Invalid2FaCode
	}
	return nil
}

func (service *TwoFaService) CheckTwoFAEnable(ctx context.Context, userId uint) (bool, error) {
	user, err := service.UserService.GetById(ctx, userId)
	if err != nil {
		logger.LogErrorWithFields("failed to get user by id", service, models.UserRole, userId, err)
		return false, errs.SomeThingWentWrong
	}

	if user.TwoFaEnabled {
		return true, nil
	}

	return false, nil
}

func (service *TwoFaService) VerifyRecoveryCode(ctx context.Context, userId uint, recoveryCode string) error {
	user, err := service.UserService.GetById(ctx, userId)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get user by id", service, err)
		return errs.SomeThingWentWrong
	}

	if !user.TwoFaEnabled {
		logger.LogErrorWithFieldsV2(ctx, "user 2FA is not active", service, errs.NotActiveTwoFactor)
		return errs.NotActiveTwoFactor
	}

	for index, storedCode := range user.RecoveryCodes {
		if isValid, _ := hash.VerifyStoredHash([]byte(storedCode), recoveryCode); isValid {
			//remove the used recovery code from the list
			user.RecoveryCodes = append(user.RecoveryCodes[:index], user.RecoveryCodes[index+1:]...)
			_, err = service.UserService.Update(ctx, user)
			if err != nil {
				logger.LogErrorWithFieldsV2(ctx, "failed to update user", service, err)
				return errs.SomeThingWentWrong
			}
			return nil
		}
	}

	return errs.InvalidRecoveryCode
}
