package authentication

import (
	"athena/src/api/errs"
	Authentication "athena/src/api/http/requests/authentication"
	"athena/src/cache"
	"athena/src/config"
	"athena/src/hash"
	"athena/src/models"
	"athena/src/models/consts"
	"athena/src/pkg/logger"
	"athena/src/pkg/utils"
	"athena/src/services"
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/docker/distribution/uuid"
	"github.com/redis/go-redis/v9"
)

type LoginService struct {
	UserService        services.IUserService
	JwtService         IJwtService
	AccessTokenService IAccessTokenService
	OTPService         services.IOTPService
	TwoFaService       ITwoFaService
	//ProfileService      services.IProfile
	//NotificationService *services.NotificationProcessService
}

type ILoginService interface {
	CheckUserCredentials(ctx context.Context, nationalIdentityCode, password, ownerType string) (int, error)
	GetAccessTokenViaOtp(ctx context.Context, nationalIdentityCode, otp, ownerType string) (*JwtDTO, error)
	UserLoginViaPassword(ctx context.Context) (*JwtDTO, error)
	UserLoginVerifyOTP(ctx context.Context) (*JwtDTO, error)
	LoginViaOtpSendOtp(ctx context.Context) (string, error)
	ResendLoginOTP(ctx context.Context) error
	//GetSecretKey(ctx context.Context) (string, string, error)
	//LoginVerifyTwoFaCodeFirstTime(ctx context.Context) (*JwtDTO, error)
	//LoginVerifyTwoFaCode(ctx context.Context) (interface{}, bool, error)
	///	ForceChangePassword(ctx context.Context) (*JwtDTO, error)
	//	GetTotp(ctx context.Context) (string, string, error)
	//	ForceChangePasswordForUser(ctx context.Context) (*JwtDTO, error)
}

func (service *LoginService) CheckUserCredentials(ctx context.Context, nationalIdentityCode, password, ownerType string) (int, error) {
	var owner interface{}
	var err error
	switch ownerType {
	case "user":
		owner, err = service.UserService.GetByNationalIdentityCode(ctx, nationalIdentityCode)

	default:
		logger.LogErrorWithFieldsV2(ctx, "unsupported owner type", service, errors.New("unsupported owner type"))
		return 0, errs.ErrAuthenticationFailed
	}
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to owner by national identity code", service, err)
		return 0, errs.ErrAuthenticationFailed
	}

	// Cast owner to the respective model type (UserModel or AdminModel).
	var passwordHashCheck bool
	var mobile string
	switch owner.(type) {
	case *models.UserModel:
		user := owner.(*models.UserModel)
		mobile = user.Mobile
		passwordHashCheck, _ = hash.VerifyStoredHash(user.Password, password)
	default:
		return 0, errs.ErrAuthenticationFailed
	}

	if !passwordHashCheck {
		logger.LogErrorWithFieldsV2(ctx, "password hash check failed", service, errors.New("the password is not correct"))
		return 0, errs.ErrAuthenticationFailed
	}

	//check for optional login
	switch owner.(type) {
	case *models.UserModel:
		//TODO:remove this on production
		if mobile == "09306637036" {
			return 10000, nil
		}
		err := service.OTPService.RequestOTP(ctx, mobile)
		if utils.CheckError(err, errs.ErrOTPRequired) || utils.CheckError(err, errs.ErrAuthOTPExists) {
			logger.LogErrorWithFieldsV2(ctx, "failed to request otp due the otp is already exist", service, err)
			return 10003, err
		}
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to request otp", service, err)
			return 0, err
		}
		return 10000, nil

	default:
		return 0, errs.ErrAuthenticationFailed
	}
}

func (service *LoginService) GetAccessTokenViaOtp(ctx context.Context, nationalIdentityCode, otp, ownerType string) (*JwtDTO, error) {
	var owner interface{}
	var err error
	var mobile string
	switch ownerType {
	case "user":
		owner, err = service.UserService.GetByNationalIdentityCode(ctx, nationalIdentityCode)
		user := owner.(*models.UserModel)
		mobile = user.Mobile
	default:
		logger.LogErrorWithFieldsV2(ctx, "unsupported owner type", service, errors.New("unsupported owner type"),
			zap.String("given-owner-type", ownerType),
		)
		return nil, errs.ErrAuthenticationFailed
	}
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to owner by national identity code", service, err)
		return nil, errs.ErrAuthenticationFailed
	}

	otpIsValid, err := service.OTPService.VerifyOTP(ctx, mobile, otp)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify otp", service, err,
			zap.String("otp", otp),
			zap.String("mobile", mobile),
		)
		return nil, errs.SomeThingWentWrong
	}
	if !otpIsValid {
		logger.LogErrorWithFieldsV2(ctx, "not a valid otp", service, errs.ErrOTPInvalid,
			zap.String("otp", otp),
			zap.String("mobile", mobile),
		)
		return nil, errs.ErrOTPInvalid
	}
	//generate token
	jwtDTO, err := service.JwtService.Generate(ctx, ownerType)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to generate jwt", service, err)
		return nil, errs.ErrAuthenticationFailed
	}

	// Store tokens in database
	ip := ctx.Value("request-ip").(string)
	userAgent := ctx.Value("request-user-agent").(string)

	_, err = service.AccessTokenService.Create(ctx, owner, jwtDTO, ip, userAgent)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to store token in db", service, err)
		return nil, errs.SomeThingWentWrong
	}

	// Return the generated tokens upon successful authentication and token generation.
	return jwtDTO, nil
}

func (service *LoginService) UserLoginViaPassword(ctx context.Context) (*JwtDTO, error) {
	req, ok := ctx.Value("req").(*Authentication.LoginRequest)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to read login request", service, nil)
		return nil, errs.ErrAuthenticationFailed
	}

	user, err := service.UserService.GetByNationalIdentityCode(ctx, req.NationalIdentityCode)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "user not found with NIC", service, err)
		return nil, errs.ErrAuthenticationFailed
	}

	passwordOK, err := hash.VerifyStoredHash(user.Password, req.Password)
	if err != nil || !passwordOK {
		logger.LogErrorWithFieldsV2(ctx, "password mismatch", service, err)
		return nil, errs.ErrAuthenticationFailed
	}

	ownerType := "user"
	jwtDTO, err := service.JwtService.Generate(ctx, ownerType)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "jwt generation failed", service, err)
		return nil, errs.ErrAuthenticationFailed
	}

	ip := ctx.Value("request-ip").(string)
	userAgent := ctx.Value("request-user-agent").(string)

	_, err = service.AccessTokenService.Create(ctx, user, jwtDTO, ip, userAgent)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "db token insert failed", service, err)
		return nil, errs.ErrAuthenticationFailed
	}

	return jwtDTO, nil
}

func (service *LoginService) LoginViaOtpSendOtp(ctx context.Context) (string, error) {
	// fetch data
	req, ok := ctx.Value("req").(*Authentication.LoginViaOTPSendOtpRequest)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to get the request from context", service, nil)
		return "", errs.SomeThingWentWrong
	}

	user, err := service.UserService.GetByNationalIdentityCode(ctx, req.NationalIdentityCode)
	if err != nil {
		logger.LogInfo(ctx, "user not found during login, returning dummy key to prevent enumeration", service)
		return "", errs.SomeThingWentWrong
	}
	// marshal the req to save in redis
	reqData, err := json.Marshal(req)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to marshal request", service, err)
		return "", errs.SomeThingWentWrong
	}

	// get expire time
	expiration, err := strconv.Atoi(config.GetInstance().Get("LOGIN_SAVE_STATE_LIFETIME"))
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to convert LOGIN_SAVE_STATE_LIFETIME to int", service, err)
		expiration = 300
	}

	// Save the request data in Redis
	key := uuid.Generate().String()
	err = cache.GetInstance().GetClient().Set(ctx, key, reqData, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to save it in redis", service, err)
		return "", errs.SomeThingWentWrong
	}

	err = service.OTPService.RequestOTP(ctx, user.Mobile)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to request otp", service, err,
			zap.String("mobile", user.Mobile),
		)
		return "", err
	}

	return key, nil
}

func (service *LoginService) ResendLoginOTP(ctx context.Context) error {
	req, ok := ctx.Value("req").(*Authentication.ResendLoginOTP)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to resend login otp", service, nil)
		return errs.SomeThingWentWrong
	}

	res, err := cache.GetInstance().GetClient().Get(context.Background(), req.LoginKey).Result()
	if err != nil {
		if utils.CheckError(err, redis.Nil) {
			logger.LogErrorWithFieldsV2(ctx, "failed to verify login otp", service, err)
			return errs.ErrLoginTimeOut
		}

		logger.LogErrorWithFieldsV2(ctx, "failed to resend login otp", service, err)
		return errs.SomeThingWentWrong
	}

	var resp Authentication.LoginViaOTPSendOtpRequest
	err = json.Unmarshal([]byte(res), &resp)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to resend login otp", service, err)
		return errs.SomeThingWentWrong
	}
	nationalCode := resp.NationalIdentityCode
	if resp.NationalCompanyId != "" {
		nationalCode = resp.NationalCompanyId
	}
	ctx = context.WithValue(ctx, consts.NationalIdentityCode, nationalCode)

	user, err := service.UserService.GetByNationalIdentityCode(ctx, resp.NationalIdentityCode)
	if err != nil && !utils.CheckError(err, errs.RecordNotFound) {
		logger.LogErrorWithFieldsV2(ctx, "failed to get user by national", service, err)
		return errs.ErrAuthenticationFailed
	}

	if user == nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get user by national", service, nil)

		return errs.ErrAuthenticationFailed
	}

	err = service.OTPService.RequestOTP(ctx, user.Mobile)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to resend register otp", service, err,
			zap.String("mobile", user.Mobile),
		)
		return err
	}

	return nil
}

func (service *LoginService) UserLoginVerifyOTP(ctx context.Context) (*JwtDTO, error) {
	req, ok := ctx.Value("req").(*Authentication.VerifyLoginOTP)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify login otp", service, nil)
		return nil, errs.SomeThingWentWrong
	}
	res, err := cache.GetInstance().GetClient().Get(context.Background(), req.LoginKey).Result()
	if err != nil {
		if utils.CheckError(err, redis.Nil) {
			logger.LogErrorWithFieldsV2(ctx, "failed to verify login otp", service, err)
			return nil, errs.ErrLoginTimeOut
		}

		logger.LogErrorWithFieldsV2(ctx, "failed to verify login otp", service, err)
		return nil, errs.ErrAuthenticationFailed
	}

	var resp Authentication.LoginViaOTPSendOtpRequest
	err = json.Unmarshal([]byte(res), &resp)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify login otp", service, err)
		return nil, errs.SomeThingWentWrong
	}
	// get user
	user, err := service.UserService.GetByNationalIdentityCode(ctx, resp.NationalIdentityCode)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get admin by national identity code", service, err)
		return nil, errs.ErrAuthenticationFailed
	}
	var otpIsValid bool
	otpIsValid, err = service.OTPService.VerifyOTP(ctx, user.Mobile, req.OTP)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify login otp", service, err,
			zap.String("mobile", user.Mobile),
			zap.String("otp", req.OTP),
		)
		return nil, errs.ErrAuthenticationFailed
	}

	if !otpIsValid {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify login otp", service, nil,
			zap.String("mobile", user.Mobile),
			zap.String("otp", req.OTP),
		)
		return nil, errs.ErrOTPInvalid
	}

	//generate token
	ownerType := "user"
	jwtDTO, err := service.JwtService.Generate(ctx, ownerType)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to generate jwt", service, err)
		return nil, errs.SomeThingWentWrong
	}

	// Store tokens in database
	ip := ctx.Value("request-ip").(string)
	userAgent := ctx.Value("request-user-agent").(string)

	_, err = service.AccessTokenService.Create(ctx, user, jwtDTO, ip, userAgent)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to store token in db", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return jwtDTO, nil
}

// GetSecretKey for setting up two-factor authentication
//func (service *LoginService) GetSecretKey(ctx context.Context) (string, string, error) {
//	url, secret, err := service.TwoFaService.GetSecretKey(ctx)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to enable two-factor secret", service, err)
//		return "", "", errs.SomeThingWentWrong
//	}
//
//	return url, secret, nil
//}
//
// LoginVerifyTwoFaCodeFirstTime for setting up two-factor authentication
//func (service *LoginService) LoginVerifyTwoFaCodeFirstTime(ctx context.Context) (*JwtDTO, error) {
//	req, ok := ctx.Value("req").(*Authentication.VerifyTwoFa)
//	if !ok {
//		logger.LogErrorWithFieldsV2(ctx, "failed to verify two-factor code", service, nil)
//		return nil, errs.ErrAuthenticationFailed
//	}
//
//	res, err := cache.GetInstance().GetClient().Get(context.Background(), req.LoginKey).Result()
//	if err != nil {
//		if utils.CheckError(err, redis.Nil) {
//			logger.LogErrorWithFieldsV2(ctx, "failed to verify two-factor code", service, err)
//			return nil, errs.ErrLoginTimeOut
//		}
//
//		logger.LogErrorWithFieldsV2(ctx, "failed to verify two-factor code", service, err)
//		return nil, errs.ErrAuthenticationFailed
//	}
//
//	var resp Authentication.LoginRequest
//	err = json.Unmarshal([]byte(res), &resp)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to verify two-factor code", service, err)
//		return nil, errs.ErrAuthenticationFailed
//	}
//	nationalCode := resp.NationalIdentityCode
//	if resp.NationalCompanyId != "" {
//		nationalCode = resp.NationalCompanyId
//	}
//	ctx = context.WithValue(ctx, consts.NationalIdentityCode, nationalCode)
//
//	// get user
//	admin, err := service.AdminService.GetByNationalIdentityCode(ctx, resp.NationalIdentityCode)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to get admin by national identity code", service, err)
//		return nil, errs.ErrAuthenticationFailed
//	}
//
//	err = service.TwoFaService.VerifyCodeFirstTime(ctx, admin.ID, req.TwoFaCode)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to verify two-factor code", service, err)
//		return nil, errs.ErrAuthenticationFailed
//	}
//
//	//generate token
//	ownerType := "admin"
//	jwtDTO, err := service.JwtService.Generate(ctx, ownerType)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to generate jwt", service, err)
//		return nil, errs.ErrAuthenticationFailed
//	}
//
//	// Store tokens in database
//	ip := ctx.Value("request-ip").(string)
//	userAgent := ctx.Value("request-user-agent").(string)
//	_, err = service.AccessTokenService.Create(ctx, admin, jwtDTO, ip, userAgent)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to verify two-factor code", service, err)
//		return nil, errs.ErrAuthenticationFailed
//	}
//
//	return jwtDTO, nil
//
//}
//
//func (service *LoginService) LoginVerifyTwoFaCode(ctx context.Context) (interface{}, bool, error) {
//	req, ok := ctx.Value("req").(*Authentication.VerifyTwoFa)
//	if !ok {
//		logger.LogErrorWithFieldsV2(ctx, "failed to verify two-factor code", service, nil)
//		return nil, false, errs.ErrAuthenticationFailed
//	}
//
//	res, err := cache.GetInstance().GetClient().Get(context.Background(), req.LoginKey).Result()
//	if err != nil {
//		if utils.CheckError(err, redis.Nil) {
//			logger.LogErrorWithFieldsV2(ctx, "failed to verify two-factor code", service, err)
//			return nil, false, errs.ErrLoginTimeOut
//		}
//
//		logger.LogErrorWithFieldsV2(ctx, "failed to verify two-factor code", service, err)
//		return nil, false, errs.ErrAuthenticationFailed
//	}
//
//	var resp Authentication.LoginRequest
//	err = json.Unmarshal([]byte(res), &resp)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to verify two-factor code", service, err)
//		return nil, false, errs.ErrAuthenticationFailed
//	}
//	nationalCode := resp.NationalIdentityCode
//	if resp.NationalCompanyId != "" {
//		nationalCode = resp.NationalCompanyId
//	}
//	ctx = context.WithValue(ctx, consts.NationalIdentityCode, nationalCode)
//
//	// get user
//	var admin *models.AdminModel
//	if resp.NationalIdentityCode != "" { // login by national id
//		adminIn, err := service.AdminService.GetByNationalIdentityCode(ctx, resp.NationalIdentityCode)
//		if err != nil {
//			logger.LogErrorWithFieldsV2(ctx, "failed to get admin by national identity code", service, err)
//			return nil, false, errs.ErrAuthenticationFailed
//		}
//		admin = adminIn
//	} else { // login by username
//		adminIn, err := service.AdminService.GetByUsername(ctx, resp.Username)
//		if err != nil {
//			logger.LogErrorWithFieldsV2(ctx, "failed to get admin by username", service, err)
//			return nil, false, errs.ErrAuthenticationFailed
//		}
//		admin = adminIn
//	}
//
//	if *admin.IsFirstLogin {
//		err = service.TwoFaService.VerifyCodeFirstTime(ctx, admin.ID, req.TwoFaCode)
//		if err != nil {
//			logger.LogErrorWithFieldsV2(ctx, "failed to verify two-factor code first time", service, err)
//			return nil, *admin.IsFirstLogin, err
//		}
//		return nil, *admin.IsFirstLogin, nil
//	}
//	err = service.TwoFaService.VerifyCode(ctx, admin.ID, req.TwoFaCode)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to verify two-factor code", service, err)
//		return nil, *admin.IsFirstLogin, errs.ErrAuthenticationFailed
//	}
//
//	//generate token
//	ownerType := "admin"
//	jwtDTO, err := service.JwtService.Generate(ctx, ownerType)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to generate jwt", service, err)
//		return nil, *admin.IsFirstLogin, errs.ErrAuthenticationFailed
//	}
//
//	// Store tokens in database
//	ip := ctx.Value("request-ip").(string)
//	userAgent := ctx.Value("request-user-agent").(string)
//	_, err = service.AccessTokenService.Create(ctx, admin, jwtDTO, ip, userAgent)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to verify two-factor code", service, err)
//		return nil, *admin.IsFirstLogin, errs.ErrAuthenticationFailed
//	}
//
//	return jwtDTO, *admin.IsFirstLogin, nil
//
//}
//
//func (service *LoginService) ForceChangePassword(ctx context.Context) (*JwtDTO, error) {
//	req, ok := ctx.Value("req").(Authentication.ChangePasswordRequest)
//	if !ok {
//		logger.LogErrorWithFieldsV2(ctx, "failed to save the change password state and send otp", service, nil)
//		return nil, errs.ChangePasswordFailed
//	}
//
//	redisData, err := cache.GetInstance().GetClient().Get(context.Background(), req.LoginKey).Result()
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to save the change password state and send otp", service, err)
//		return nil, errs.SomeThingWentWrong
//	}
//	var resp *Authentication.LoginRequest
//	err = json.Unmarshal([]byte(redisData), &resp)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to unmarshal data from redis", service, err)
//		return nil, errs.SomeThingWentWrong
//	}
//
//	nationalCode := resp.NationalIdentityCode
//	if resp.NationalCompanyId != "" {
//		nationalCode = resp.NationalCompanyId
//	}
//	if nationalCode != "" {
//		ctx = context.WithValue(ctx, consts.NationalIdentityCode, nationalCode)
//	}
//	if resp.Username != "" {
//		ctx = context.WithValue(ctx, consts.Username, resp.Username)
//	}
//
//	var admin *models.AdminModel
//	if resp.NationalIdentityCode != "" { // login by national id
//		adminIn, err := service.AdminService.GetByNationalIdentityCode(ctx, resp.NationalIdentityCode)
//		if err != nil {
//			logger.LogErrorWithFieldsV2(ctx, "failed to get admin by national identity code", service, err)
//			return nil, errs.ErrAuthenticationFailed
//		}
//		admin = adminIn
//	} else { // login by username
//		adminIn, err := service.AdminService.GetByUsername(ctx, resp.Username)
//		if err != nil {
//			logger.LogErrorWithFieldsV2(ctx, "failed to get admin by username", service, err)
//			return nil, errs.ErrAuthenticationFailed
//		}
//		admin = adminIn
//	}
//	passwordHashCheck, err := hash.VerifyStoredHash(admin.Password, req.NewPassword)
//	if err != nil || passwordHashCheck {
//		logger.LogErrorWithFieldsV2(ctx, "failed to verify password", service, nil)
//		return nil, errs.PasswordShouldBeNew
//	}
//
//	admin.Password = []byte(req.NewPassword)
//	isFirstLogin := false
//	admin.IsFirstLogin = &isFirstLogin
//	_, err = service.AdminService.Update(ctx, admin)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to update admin", service, err)
//		return nil, errs.SomeThingWentWrong
//	}
//
//	ownerType := "admin"
//	jwtDTO, err := service.JwtService.Generate(ctx, ownerType)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to generate jwt", service, err)
//		return nil, errs.ErrAuthenticationFailed
//	}
//
//	// Store tokens in database
//	ip := ctx.Value("request-ip").(string)
//	userAgent := ctx.Value("request-user-agent").(string)
//	_, err = service.AccessTokenService.Create(ctx, admin, jwtDTO, ip, userAgent)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to verify two-factor code", service, err)
//		return nil, errs.ErrAuthenticationFailed
//	}
//
//	return jwtDTO, nil
//
//}
//
//func (service *LoginService) GetTotp(ctx context.Context) (string, string, error) {
//	loginKey := ctx.Value("login_key").(string)
//	redisData, err := cache.GetInstance().GetClient().Get(context.Background(), loginKey).Result()
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to save the change password state and send otp", service, err)
//		return "", "", errs.SomeThingWentWrong
//	}
//	var resp *Authentication.LoginRequest
//	err = json.Unmarshal([]byte(redisData), &resp)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to unmarshal data from redis", service, err)
//		return "", "", errs.SomeThingWentWrong
//	}
//
//	nationalCode := resp.NationalIdentityCode
//	if resp.NationalCompanyId != "" {
//		nationalCode = resp.NationalCompanyId
//	}
//	if nationalCode != "" {
//		ctx = context.WithValue(ctx, consts.NationalIdentityCode, nationalCode)
//	}
//	if resp.Username != "" {
//		ctx = context.WithValue(ctx, consts.Username, resp.Username)
//	}
//
//	var admin *models.AdminModel
//	if resp.NationalIdentityCode != "" { // login by national id
//		adminIn, err := service.AdminService.GetByNationalIdentityCode(ctx, resp.NationalIdentityCode)
//		if err != nil {
//			logger.LogErrorWithFieldsV2(ctx, "failed to get admin by national identity code", service, err)
//			return "", "", errs.ErrAuthenticationFailed
//		}
//		admin = adminIn
//	} else { // login by username
//		adminIn, err := service.AdminService.GetByUsername(ctx, resp.Username)
//		if err != nil {
//			logger.LogErrorWithFieldsV2(ctx, "failed to get admin by username", service, err)
//			return "", "", errs.ErrAuthenticationFailed
//		}
//		admin = adminIn
//	}
//	decryptedSecret, err := encrypt.GetInstance().Decrypt(admin.TotpSecret)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to decrypt secret", service, err)
//		return "", "", errs.SomeThingWentWrong
//	}
//	decryptedSecretUrl, err := encrypt.GetInstance().Decrypt(admin.TotpSecretUrl)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to decrypt secret url", service, err)
//		return "", "", errs.SomeThingWentWrong
//	}
//
//	return string(decryptedSecret), string(decryptedSecretUrl), nil
//}
//
//func (service *LoginService) ForceChangePasswordForUser(ctx context.Context) (*JwtDTO, error) {
//
//	loginKey, ok := ctx.Value("login_key").(string)
//	if !ok || loginKey == "" {
//		logger.LogErrorWithFieldsV2(ctx, "login_key not found in context", service, nil)
//		return nil, errs.ChangePasswordFailed
//	}
//
//	redisData, err := cache.GetInstance().
//		GetClient().
//		Get(context.Background(), loginKey).
//		Result()
//
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to get login_key from redis", service, err)
//		return nil, errs.SomeThingWentWrong
//	}
//
//	var loginReq Authentication.LoginRequest
//	err = json.Unmarshal([]byte(redisData), &loginReq)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to unmarshal redis data", service, err)
//		return nil, errs.SomeThingWentWrong
//	}
//
//	user, err := service.UserService.GetByNationalIdentityCode(ctx, loginReq.NationalIdentityCode)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to get user by national identity code", service, err)
//		return nil, errs.ErrAuthenticationFailed
//	}
//
//	req, ok := ctx.Value("req").(*Authentication.ForceChangePasswordUserRequest)
//	if !ok {
//		logger.LogErrorWithFieldsV2(ctx, "invalid force change password request", service, nil)
//		return nil, errs.ChangePasswordFailed
//	}
//
//	passwordHashCheck, err := hash.VerifyStoredHash(user.Password, req.NewPassword)
//	if err != nil || passwordHashCheck {
//		logger.LogErrorWithFieldsV2(ctx, "new password must be different from old", service, nil)
//		return nil, errs.PasswordShouldBeNew
//	}
//
//	user.Password = []byte(req.NewPassword)
//
//	isFirstLogin := false
//	user.IsFirstLogin = &isFirstLogin
//
//	_, err = service.UserService.Update(ctx, user)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to update user", service, err)
//		return nil, errs.SomeThingWentWrong
//	}
//
//	jwtDTO, err := service.JwtService.Generate(ctx, "user")
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to generate jwt", service, err)
//		return nil, errs.ErrAuthenticationFailed
//	}
//
//	ip := ctx.Value("request-ip").(string)
//	userAgent := ctx.Value("request-user-agent").(string)
//
//	_, err = service.AccessTokenService.Create(ctx, user, jwtDTO, ip, userAgent)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to store token", service, err)
//		return nil, errs.ErrAuthenticationFailed
//	}
//
//	cache.GetInstance().GetClient().Del(context.Background(), loginKey)
//
//	return jwtDTO, nil
//}
