package authentication

import (
	"athena/src/api/errs"
	adminLogin "athena/src/api/http/requests/Admin/authentication/login"
	"athena/src/api/http/requests/authentication"
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
	AdminService       services.IAdminService
	TwoFaService       ITwoFaService
}

type LoginResult struct {
	Tokens        *JwtDTO
	TwoFaRequired bool
	LoginKey      string
}

type twoFaLoginState struct {
	OwnerID   uint   `json:"owner_id"`
	OwnerType string `json:"owner_type"`
}

type ILoginService interface {
	CheckUserCredentials(ctx context.Context, nationalIdentityCode, password, ownerType string) (int, error)
	GetAccessTokenViaOtp(ctx context.Context, nationalIdentityCode, otp, ownerType string) (*JwtDTO, error)
	UserLoginViaPassword(ctx context.Context) (*LoginResult, error)
	UserLoginVerifyOTP(ctx context.Context) (*LoginResult, error)
	LoginViaOtpSendOtp(ctx context.Context) (string, error)
	ResendLoginOTP(ctx context.Context) error
	AdminLoginViaPassword(ctx context.Context) (*LoginResult, error)
	VerifyTwoFactorLogin(ctx context.Context) (*JwtDTO, error)
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

func (service *LoginService) UserLoginViaPassword(ctx context.Context) (*LoginResult, error) {
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

	if user.IsActive != nil && !*user.IsActive {
		logger.LogErrorWithFieldsV2(ctx, "user is deactivated", service, errs.ErrDeactivatedUser)
		return nil, errs.ErrDeactivatedUser
	}

	passwordOK, err := hash.VerifyStoredHash(user.Password, req.Password)
	if err != nil || !passwordOK {
		logger.LogErrorWithFieldsV2(ctx, "password mismatch", service, err)
		return nil, errs.ErrAuthenticationFailed
	}

	return service.completeLogin(ctx, user, models.UserRole, user.ID, user.TwoFaEnabled)
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

func (service *LoginService) UserLoginVerifyOTP(ctx context.Context) (*LoginResult, error) {
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
	user, err := service.UserService.GetByNationalIdentityCode(ctx, resp.NationalIdentityCode)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get user by national identity code", service, err)
		return nil, errs.ErrAuthenticationFailed
	}
	if user.IsActive != nil && !*user.IsActive {
		logger.LogErrorWithFieldsV2(ctx, "user is deactivated", service, errs.ErrDeactivatedUser)
		return nil, errs.ErrDeactivatedUser
	}

	otpIsValid, err := service.OTPService.VerifyOTP(ctx, user.Mobile, req.OTP)
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

	return service.completeLogin(ctx, user, models.UserRole, user.ID, user.TwoFaEnabled)
}

func (service *LoginService) AdminLoginViaPassword(ctx context.Context) (*LoginResult, error) {
	req, ok := ctx.Value("req").(*adminLogin.LoginViaPasswordRequest)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to read admin login request", service, nil)
		return nil, errs.ErrAuthenticationFailed
	}

	admin, err := service.AdminService.GetByUsername(ctx, req.Username)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "admin not found with username", service, err)
		return nil, errs.ErrAuthenticationFailed
	}

	if admin.IsActive == nil || !*admin.IsActive {
		logger.LogErrorWithFieldsV2(ctx, "admin is deactivated", service, errs.ErrDeactivatedAdmin)
		return nil, errs.ErrDeactivatedAdmin
	}

	passwordOK, err := hash.VerifyStoredHash(admin.Password, req.Password)
	if err != nil || !passwordOK {
		logger.LogErrorWithFieldsV2(ctx, "password mismatch", service, err)
		return nil, errs.ErrAuthenticationFailed
	}

	return service.completeLogin(ctx, admin, models.AdminRole, admin.ID, admin.TwoFaEnabled)
}

func (service *LoginService) VerifyTwoFactorLogin(ctx context.Context) (*JwtDTO, error) {
	req, ok := ctx.Value("req").(*Authentication.VerifyTwoFa)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to read two factor login request", service, nil)
		return nil, errs.ErrAuthenticationFailed
	}
	if req.TwoFaCode == "" && req.RecoveryCode == "" {
		return nil, errs.ErrTwoFactorChallengeMissing
	}

	res, err := cache.GetInstance().GetClient().Get(ctx, twoFaLoginCacheKey(req.LoginKey)).Result()
	if err != nil {
		if utils.CheckError(err, redis.Nil) {
			logger.LogErrorWithFieldsV2(ctx, "two factor login state expired", service, err)
			return nil, errs.ErrLoginTimeOut
		}
		logger.LogErrorWithFieldsV2(ctx, "failed to load two factor login state", service, err)
		return nil, errs.ErrAuthenticationFailed
	}

	var state twoFaLoginState
	if err = json.Unmarshal([]byte(res), &state); err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to unmarshal two factor login state", service, err)
		return nil, errs.SomeThingWentWrong
	}

	ctx = context.WithValue(ctx, consts.OwnerType, state.OwnerType)
	ctx = context.WithValue(ctx, consts.OwnerId, state.OwnerID)

	if req.TwoFaCode != "" {
		if err = service.TwoFaService.VerifyCode(ctx, state.OwnerID, req.TwoFaCode); err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to verify two factor code during login", service, err)
			return nil, err
		}
	} else if err = service.TwoFaService.VerifyRecoveryCode(ctx, state.OwnerID, req.RecoveryCode); err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify recovery code during login", service, err)
		return nil, err
	}

	owner, err := service.loadLoginOwner(ctx, state.OwnerType, state.OwnerID)
	if err != nil {
		return nil, err
	}

	tokens, err := service.issueTokens(ctx, owner, state.OwnerType)
	if err != nil {
		return nil, err
	}

	_ = cache.GetInstance().GetClient().Del(ctx, twoFaLoginCacheKey(req.LoginKey)).Err()
	return tokens, nil
}

func (service *LoginService) completeLogin(ctx context.Context, owner interface{}, ownerType string, ownerID uint, twoFaEnabled bool) (*LoginResult, error) {
	if twoFaEnabled {
		loginKey, err := service.saveTwoFaLoginState(ctx, ownerID, ownerType)
		if err != nil {
			return nil, err
		}
		return &LoginResult{TwoFaRequired: true, LoginKey: loginKey}, nil
	}

	tokens, err := service.issueTokens(ctx, owner, ownerType)
	if err != nil {
		return nil, err
	}
	return &LoginResult{Tokens: tokens}, nil
}

func (service *LoginService) issueTokens(ctx context.Context, owner interface{}, ownerType string) (*JwtDTO, error) {
	jwtDTO, err := service.JwtService.Generate(ctx, ownerType)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "jwt generation failed", service, err)
		return nil, errs.ErrAuthenticationFailed
	}

	ip, _ := ctx.Value("request-ip").(string)
	userAgent, _ := ctx.Value("request-user-agent").(string)

	_, err = service.AccessTokenService.Create(ctx, owner, jwtDTO, ip, userAgent)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "db token insert failed", service, err)
		return nil, errs.ErrAuthenticationFailed
	}
	return jwtDTO, nil
}

func (service *LoginService) saveTwoFaLoginState(ctx context.Context, ownerID uint, ownerType string) (string, error) {
	payload, err := json.Marshal(twoFaLoginState{OwnerID: ownerID, OwnerType: ownerType})
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to marshal two factor login state", service, err)
		return "", errs.SomeThingWentWrong
	}

	expiration, err := strconv.Atoi(config.GetInstance().Get("LOGIN_SAVE_STATE_LIFETIME"))
	if err != nil {
		expiration = 300
	}

	loginKey := uuid.Generate().String()
	if err = cache.GetInstance().GetClient().Set(ctx, twoFaLoginCacheKey(loginKey), payload, time.Duration(expiration)*time.Second).Err(); err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to save two factor login state", service, err)
		return "", errs.SomeThingWentWrong
	}
	return loginKey, nil
}

func (service *LoginService) loadLoginOwner(ctx context.Context, ownerType string, ownerID uint) (interface{}, error) {
	switch ownerType {
	case models.AdminRole:
		admin, err := service.AdminService.GetById(ctx, ownerID)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to get admin for two factor login", service, err)
			return nil, errs.ErrAuthenticationFailed
		}
		return admin, nil
	default:
		user, err := service.UserService.GetById(ctx, ownerID)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to get user for two factor login", service, err)
			return nil, errs.ErrAuthenticationFailed
		}
		return user, nil
	}
}

func twoFaLoginCacheKey(loginKey string) string {
	return "login-2fa-" + loginKey
}
