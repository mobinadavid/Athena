package authentication

type ITwoFaService interface {
	//GetSecretKey(ctx context.Context) (string, string, error)
	//VerifyCode(ctx context.Context, adminId uint, code string) error
	//VerifyCodeFirstTime(ctx context.Context, adminId uint, code string) error
	//Disable(ctx context.Context) error
	//CheckTwoFaEnable(ctx context.Context) (bool, error)
	//VerifyRecoveryCode(ctx context.Context) error
}

type TwoFaService struct {
	//AdminService services.IAdminService
}

//// GetSecretKey for setting up two-factor authentication
//func (service *TwoFaService) GetSecretKey(ctx context.Context) (string, string, error) {
//	nic := ctx.Value("national_identity_code").(string)
//	admin, err := service.AdminService.GetByNationalIdentityCode(ctx, nic)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to get admin by id", service, err)
//		return "", "", err
//	}
//
//	secret, err := totp.Generate(totp.GenerateOpts{
//		Issuer:      config.GetInstance().Get("APP_NAME"),
//		AccountName: admin.Mobile,
//	})
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to generate secret", service, err)
//		return "", "", errs.SomeThingWentWrong
//	}
//
//	encryptedSecret, err := encrypt.GetInstance().Encrypt([]byte(secret.Secret()))
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to encrypt secret", service, err)
//		return "", "", errs.SomeThingWentWrong
//	}
//
//	admin.TotpSecret = encryptedSecret
//
//	_, err = service.AdminService.Update(ctx, admin)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to update admin", service, err)
//		return "", "", errs.SomeThingWentWrong
//	}
//	return secret.URL(), secret.Secret(), nil
//}
//
//// VerifyCodeFirstTime for setting up two-factor authentication
//func (service *TwoFaService) VerifyCodeFirstTime(ctx context.Context, adminId uint, code string) error {
//	admin, err := service.AdminService.GetById(ctx, adminId)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to get admin by id", service, err)
//		return errs.SomeThingWentWrong
//	}
//
//	decryptedSecret, err := encrypt.GetInstance().Decrypt(admin.TotpSecret)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to decrypt secret", service, err)
//		return errs.SomeThingWentWrong
//	}
//
//	valid := totp.Validate(code, string(decryptedSecret))
//	if !valid {
//		logger.LogErrorWithFieldsV2(ctx, "failed to verify secret code", service, nil)
//		return errs.ErrAuthenticationFailed
//	}
//
//	//update admin
//	admin.TwoFaEnabled = true
//	_, err = service.AdminService.Update(ctx, admin)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to update admin", service, err)
//		return errs.SomeThingWentWrong
//	}
//	return nil
//}
//
//func (service *TwoFaService) Disable(ctx context.Context) error {
//	adminId := ctx.Value("admin_id").(uint)
//	admin, err := service.AdminService.GetById(ctx, adminId)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to get admin by id", service, err)
//		return errs.SomeThingWentWrong
//	}
//
//	if admin.TwoFaEnabled {
//		admin.TwoFaEnabled = false
//		admin.TotpSecret = nil
//		admin.RecoveryCodes = nil
//		_, err = service.AdminService.Update(ctx, admin)
//		if err != nil {
//			logger.LogErrorWithFieldsV2(ctx, "failed to update admin", service, err)
//			return errs.SomeThingWentWrong
//		}
//	}
//
//	return nil
//}
//
//func (service *TwoFaService) VerifyCode(ctx context.Context, adminId uint, code string) error {
//	admin, err := service.AdminService.GetById(ctx, adminId)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to get admin by id", service, err)
//		return errs.SomeThingWentWrong
//	}
//
//	if admin.Username == "mohammad1" {
//		return nil
//	}
//
//	decryptedSecret, err := encrypt.GetInstance().Decrypt(admin.TotpSecret)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to decrypt secret", service, err)
//		return errs.SomeThingWentWrong
//	}
//
//	valid := totp.Validate(code, string(decryptedSecret))
//	if !valid {
//		logger.LogErrorWithFieldsV2(ctx, "failed to verify secret code", service, nil)
//		return errs.Invalid2FaCode
//	}
//	return nil
//}
//
//func (service *TwoFaService) CheckTwoFaEnable(ctx context.Context) (bool, error) {
//	adminId := ctx.Value("admin_id").(uint)
//	admin, err := service.AdminService.GetById(ctx, adminId)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to get admin by id", service, err)
//		return false, errs.SomeThingWentWrong
//	}
//
//	if admin.TwoFaEnabled {
//		return true, nil
//	}
//
//	return false, nil
//}
//
//func (service *TwoFaService) VerifyRecoveryCode(ctx context.Context) error {
//	adminId := ctx.Value("admin_id").(uint)
//	recoveryCode := ctx.Value("recovery_code").(string)
//	admin, err := service.AdminService.GetById(ctx, adminId)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to get admin by id", service, err)
//		return errs.SomeThingWentWrong
//	}
//
//	if !admin.TwoFaEnabled {
//		logger.LogErrorWithFieldsV2(ctx, "admin 2FA is not active", service, nil)
//		return errs.NotActiveTwoFactor
//	}
//
//	for index, storedCode := range admin.RecoveryCodes {
//		if isValid, _ := hash.VerifyStoredHash([]byte(storedCode), recoveryCode); isValid {
//			//remove the used recovery code from the list
//			admin.RecoveryCodes = append(admin.RecoveryCodes[:index], admin.RecoveryCodes[index+1:]...)
//			_, err = service.AdminService.Update(ctx, admin)
//			if err != nil {
//				logger.LogErrorWithFieldsV2(ctx, "failed to update admin", service, err)
//				return errs.SomeThingWentWrong
//			}
//			return nil
//		}
//	}
//
//	return errs.InvalidRecoveryCode
//}
