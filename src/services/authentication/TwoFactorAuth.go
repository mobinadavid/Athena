package authentication

import (
	"bytes"
	"context"
	"encoding/base64"
	"image/png"
	"strconv"

	"athena/src/api/errs"
	"athena/src/config"
	"athena/src/encrypt"
	"athena/src/hash"
	"athena/src/models"
	"athena/src/models/consts"
	"athena/src/pkg/logger"
	"athena/src/pkg/utils"
	"athena/src/services"

	"github.com/lib/pq"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type TwoFaService struct {
	UserService  services.IUserService
	AdminService services.IAdminService
	OTPService   services.IOTPService
}

type ITwoFaService interface {
	Enable(ctx context.Context, ownerId uint) (secretURL, secretKey, qrCode string, err error)
	VerifyCodeFirstTime(ctx context.Context, ownerId uint, code string) ([]string, error)
	Disable(ctx context.Context, ownerId uint, totpCode, recoveryCode string) error
	VerifyCode(ctx context.Context, ownerId uint, code string) error
	CheckTwoFAEnable(ctx context.Context, ownerId uint) (bool, error)
	VerifyRecoveryCode(ctx context.Context, ownerId uint, recoveryCode string) error
}

type twoFaOwner struct {
	ID            uint
	AccountName   string
	TotpSecret    []byte
	TotpSecretUrl []byte
	TwoFaEnabled  bool
	RecoveryCodes pq.StringArray
}

func (service *TwoFaService) ownerType(ctx context.Context) string {
	if v, ok := ctx.Value(consts.OwnerType).(string); ok && v != "" {
		return v
	}
	return models.UserRole
}

func (service *TwoFaService) loadOwner(ctx context.Context, ownerId uint) (*twoFaOwner, error) {
	switch service.ownerType(ctx) {
	case models.AdminRole:
		admin, err := service.AdminService.GetById(ctx, ownerId)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to get admin by id", service, err)
			return nil, errs.SomeThingWentWrong
		}
		accountName := admin.Username
		if accountName == "" {
			accountName = admin.Mobile
		}
		return &twoFaOwner{
			ID:            admin.ID,
			AccountName:   accountName,
			TotpSecret:    admin.TotpSecret,
			TotpSecretUrl: admin.TotpSecretUrl,
			TwoFaEnabled:  admin.TwoFaEnabled,
			RecoveryCodes: admin.RecoveryCodes,
		}, nil
	default:
		user, err := service.UserService.GetById(ctx, ownerId)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to get user by id", service, err)
			return nil, errs.SomeThingWentWrong
		}
		accountName := user.Mobile
		if accountName == "" {
			accountName = user.NationalIdentityCode
		}
		return &twoFaOwner{
			ID:            user.ID,
			AccountName:   accountName,
			TotpSecret:    user.TotpSecret,
			TotpSecretUrl: user.TotpSecretUrl,
			TwoFaEnabled:  user.TwoFaEnabled,
			RecoveryCodes: user.RecoveryCodes,
		}, nil
	}
}

func (service *TwoFaService) persistOwner(ctx context.Context, ownerId uint, owner *twoFaOwner) error {
	switch service.ownerType(ctx) {
	case models.AdminRole:
		admin, err := service.AdminService.GetById(ctx, ownerId)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to get admin by id", service, err)
			return errs.SomeThingWentWrong
		}
		admin.TotpSecret = owner.TotpSecret
		admin.TotpSecretUrl = owner.TotpSecretUrl
		admin.TwoFaEnabled = owner.TwoFaEnabled
		admin.RecoveryCodes = owner.RecoveryCodes
		if _, err = service.AdminService.Save(ctx, admin); err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to save admin two factor data", service, err)
			return errs.SomeThingWentWrong
		}
		return nil
	default:
		user, err := service.UserService.GetById(ctx, ownerId)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to get user by id", service, err)
			return errs.SomeThingWentWrong
		}
		user.TotpSecret = owner.TotpSecret
		user.TotpSecretUrl = owner.TotpSecretUrl
		user.TwoFaEnabled = owner.TwoFaEnabled
		user.RecoveryCodes = owner.RecoveryCodes
		if _, err = service.UserService.Save(ctx, user); err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to save user two factor data", service, err)
			return errs.SomeThingWentWrong
		}
		return nil
	}
}

func (service *TwoFaService) Enable(ctx context.Context, ownerId uint) (string, string, string, error) {
	owner, err := service.loadOwner(ctx, ownerId)
	if err != nil {
		return "", "", "", err
	}

	if owner.TwoFaEnabled {
		logger.LogErrorWithFieldsV2(ctx, "two factor authentication is already enabled", service, errs.ErrTwoFactorAlreadyEnabled)
		return "", "", "", errs.ErrTwoFactorAlreadyEnabled
	}

	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      config.GetInstance().Get("APP_NAME"),
		AccountName: owner.AccountName,
		Period:      30,
	})
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to generate totp secret", service, err)
		return "", "", "", errs.SomeThingWentWrong
	}

	encryptedSecret, err := encrypt.GetInstance().Encrypt([]byte(secret.Secret()))
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to encrypt totp secret", service, err)
		return "", "", "", errs.SomeThingWentWrong
	}
	encryptedSecretUrl, err := encrypt.GetInstance().Encrypt([]byte(secret.URL()))
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to encrypt totp secret url", service, err)
		return "", "", "", errs.SomeThingWentWrong
	}

	owner.TotpSecret = encryptedSecret
	owner.TotpSecretUrl = encryptedSecretUrl
	owner.TwoFaEnabled = false
	owner.RecoveryCodes = nil
	if err = service.persistOwner(ctx, ownerId, owner); err != nil {
		return "", "", "", err
	}

	qrCode, err := encodeTotpQR(secret)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to encode totp qr code", service, err)
		return "", "", "", errs.SomeThingWentWrong
	}

	return secret.URL(), secret.Secret(), qrCode, nil
}

func (service *TwoFaService) VerifyCodeFirstTime(ctx context.Context, ownerId uint, code string) ([]string, error) {
	owner, err := service.loadOwner(ctx, ownerId)
	if err != nil {
		return nil, err
	}

	if owner.TwoFaEnabled {
		logger.LogErrorWithFieldsV2(ctx, "two factor authentication is already enabled", service, errs.ErrTwoFactorAlreadyEnabled)
		return nil, errs.ErrTwoFactorAlreadyEnabled
	}
	if len(owner.TotpSecret) == 0 {
		logger.LogErrorWithFieldsV2(ctx, "two factor secret is missing", service, errs.ErrTwoFactorSecretMissing)
		return nil, errs.ErrTwoFactorSecretMissing
	}

	if err = validateTotp(owner.TotpSecret, code); err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify totp for first-time setup", service, err)
		return nil, err
	}

	recoveryCodeNumber, err := strconv.Atoi(config.GetInstance().Get("RECOVERY_CODES_NUMBER"))
	if err != nil || recoveryCodeNumber <= 0 {
		recoveryCodeNumber = 8
	}

	recoveryCodes, err := utils.GenerateRandomCodes(recoveryCodeNumber)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to generate recovery codes", service, err)
		return nil, errs.SomeThingWentWrong
	}

	hashedCodes := make(pq.StringArray, 0, len(recoveryCodes))
	for _, recoveryCode := range recoveryCodes {
		hashedCode, hashErr := hash.GetInstance().Generate([]byte(recoveryCode))
		if hashErr != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to hash recovery code", service, hashErr)
			return nil, errs.SomeThingWentWrong
		}
		hashedCodes = append(hashedCodes, string(hashedCode))
	}

	owner.RecoveryCodes = hashedCodes
	owner.TwoFaEnabled = true
	if err = service.persistOwner(ctx, ownerId, owner); err != nil {
		return nil, err
	}

	return recoveryCodes, nil
}

func (service *TwoFaService) Disable(ctx context.Context, ownerId uint, totpCode, recoveryCode string) error {
	owner, err := service.loadOwner(ctx, ownerId)
	if err != nil {
		return err
	}

	if !owner.TwoFaEnabled {
		logger.LogErrorWithFieldsV2(ctx, "two factor authentication is not active", service, errs.NotActiveTwoFactor)
		return errs.NotActiveTwoFactor
	}

	if totpCode == "" && recoveryCode == "" {
		return errs.ErrTwoFactorChallengeMissing
	}

	if totpCode != "" {
		if err = validateTotp(owner.TotpSecret, totpCode); err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to verify totp while disabling two factor", service, err)
			return err
		}
	} else if err = service.consumeRecoveryCode(owner, recoveryCode); err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify recovery code while disabling two factor", service, err)
		return err
	}

	owner.TwoFaEnabled = false
	owner.TotpSecret = nil
	owner.TotpSecretUrl = nil
	owner.RecoveryCodes = nil
	return service.persistOwner(ctx, ownerId, owner)
}

func (service *TwoFaService) VerifyCode(ctx context.Context, ownerId uint, code string) error {
	owner, err := service.loadOwner(ctx, ownerId)
	if err != nil {
		return err
	}

	if !owner.TwoFaEnabled {
		logger.LogErrorWithFieldsV2(ctx, "two factor authentication is not active", service, errs.NotActiveTwoFactor)
		return errs.NotActiveTwoFactor
	}

	if err = validateTotp(owner.TotpSecret, code); err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify totp code", service, err)
		return err
	}
	return nil
}

func (service *TwoFaService) CheckTwoFAEnable(ctx context.Context, ownerId uint) (bool, error) {
	owner, err := service.loadOwner(ctx, ownerId)
	if err != nil {
		return false, err
	}
	return owner.TwoFaEnabled, nil
}

func (service *TwoFaService) VerifyRecoveryCode(ctx context.Context, ownerId uint, recoveryCode string) error {
	owner, err := service.loadOwner(ctx, ownerId)
	if err != nil {
		return err
	}

	if !owner.TwoFaEnabled {
		logger.LogErrorWithFieldsV2(ctx, "two factor authentication is not active", service, errs.NotActiveTwoFactor)
		return errs.NotActiveTwoFactor
	}

	if err = service.consumeRecoveryCode(owner, recoveryCode); err != nil {
		return err
	}

	return service.persistOwner(ctx, ownerId, owner)
}

func (service *TwoFaService) consumeRecoveryCode(owner *twoFaOwner, recoveryCode string) error {
	for index, storedCode := range owner.RecoveryCodes {
		if isValid, _ := hash.VerifyStoredHash([]byte(storedCode), recoveryCode); isValid {
			owner.RecoveryCodes = append(owner.RecoveryCodes[:index], owner.RecoveryCodes[index+1:]...)
			return nil
		}
	}
	return errs.InvalidRecoveryCode
}

func validateTotp(encryptedSecret []byte, code string) error {
	if len(encryptedSecret) == 0 {
		return errs.ErrTwoFactorSecretMissing
	}
	decryptedSecret, err := encrypt.GetInstance().Decrypt(encryptedSecret)
	if err != nil {
		return errs.SomeThingWentWrong
	}
	if !totp.Validate(code, string(decryptedSecret)) {
		return errs.Invalid2FaCode
	}
	return nil
}

func encodeTotpQR(secret *otp.Key) (string, error) {
	img, err := secret.Image(256, 256)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err = png.Encode(&buf, img); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
