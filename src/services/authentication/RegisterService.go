package authentication

import (
	"athena/src/api/errs"
	"athena/src/api/http/requests/Users/UserRequests"
	authenticationrequests "athena/src/api/http/requests/authentication"
	"athena/src/cache"
	"athena/src/config"
	"athena/src/models/consts"
	"athena/src/pkg/logger"
	"athena/src/pkg/utils"
	"athena/src/services"
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"strconv"
	"time"

	"github.com/docker/distribution/uuid"
	"github.com/redis/go-redis/v9"
)

type RegisterService struct {
	UserService services.IUserService
	OTPService  services.IOTPService
	//NotificationService *services.NotificationProcessService
}

type IRegisterService interface {
	SaveStateAndSendOTP(ctx context.Context) (string, error)
	VerifyRegisterOTPViaRedisKey(ctx context.Context) error
	ResendRegisterOTP(ctx context.Context) error
}

func (service *RegisterService) SaveStateAndSendOTP(ctx context.Context) (string, error) {
	// fetch data
	req, ok := ctx.Value("req").(*authenticationrequests.RegisterRequest)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to save the register state and send otp", service, nil)
		return "", errs.RegisterFailed
	}

	if req.Password != req.RePassword {
		logger.LogErrorWithFieldsV2(ctx, "password is not equal to re-password", service, nil)
		return "", errs.PasswordNotMatch
	}

	user, err := service.UserService.GetByNationalIdentityCode(ctx, req.NationalIdentityCode)
	if err != nil && !utils.CheckError(err, errs.RecordNotFound) {
		logger.LogErrorWithFieldsV2(ctx, "failed to get user by national", service, err)
		return "", errs.RegisterFailed
	}
	if user != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to save the register state and send otp", service, nil)
		return "", errs.RegisterFailed
	}

	// marshal the req to save in redis
	reqData, err := json.Marshal(req)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to save the register state and send otp", service, err)
		return "", errs.RegisterFailed
	}

	// get expire time
	expiration, err := strconv.Atoi(config.GetInstance().Get("REGISTER_SAVE_STATE_LIFETIME"))
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to save the register state and send otp", service, err)
		expiration = 120
	}

	// Save the request data in Redis
	key := uuid.Generate().String()
	err = cache.GetInstance().GetClient().Set(ctx, key, reqData, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to save the register state and send otp", service, err)
		return "", errs.RegisterFailed
	}

	err = service.OTPService.RequestOTP(ctx, req.Mobile)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to request otp", service, err,
			zap.String("mobile", req.Mobile),
		)
		return "", err
	}

	return key, nil
}

func (service *RegisterService) VerifyRegisterOTPViaRedisKey(ctx context.Context) error {
	req, ok := ctx.Value("req").(*authenticationrequests.VerifyRegisterOTP)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify register otp", service, nil)
		return errs.RegisterFailed
	}

	res, err := cache.GetInstance().GetClient().Get(context.Background(), req.RegisterKey).Result()
	if err != nil {
		if err == redis.Nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to verify register otp", service, err)
			return errs.RegisterFailed
		}

		logger.LogErrorWithFieldsV2(ctx, "failed to verify register otp", service, err)
		return errs.RegisterFailed
	}

	var resp authenticationrequests.RegisterRequest
	err = json.Unmarshal([]byte(res), &resp)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify register otp", service, err)
		return errs.RegisterFailed
	}
	nationalCode := resp.NationalIdentityCode
	ctx = context.WithValue(ctx, consts.NationalIdentityCode, nationalCode)

	var otpIsValid bool
	otpIsValid, err = service.OTPService.VerifyOTP(ctx, resp.Mobile, req.OTP)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify register otp", service, err,
			zap.String("mobile", resp.Mobile),
			zap.String("otp", req.OTP),
		)
		return errs.RegisterFailed
	}

	if !otpIsValid {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify register otp", service, nil,
			zap.String("mobile", resp.Mobile),
			zap.String("otp", req.OTP),
		)
		return errs.ErrOTPInvalid
	}

	// check password
	if resp.Password != resp.RePassword {
		logger.LogErrorWithFieldsV2(ctx, "failed to set password", service, nil)
		return errs.PasswordNotMatch
	}

	_, err = service.UserService.Create(ctx, &UserRequests.CreateRequest{
		NationalIdentityCode: resp.NationalIdentityCode,
		Mobile:               resp.Mobile,
		Password:             resp.Password,
		IsActive:             true,
	})
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to set password", service, err)
		return errs.SomeThingWentWrong
	}
	return nil
}

func (service *RegisterService) ResendRegisterOTP(ctx context.Context) error {
	req, ok := ctx.Value("req").(*authenticationrequests.ResendRegisterOTP)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to resend register otp", service, nil)
		return errs.RegisterFailed
	}

	res, err := cache.GetInstance().GetClient().Get(context.Background(), req.RegisterKey).Result()
	if err != nil {
		if err == redis.Nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to verify register otp", service, err)
			return errs.RegisterFailed
		}

		logger.LogErrorWithFieldsV2(ctx, "failed to resend register otp", service, err)
		return errs.RegisterFailed
	}

	var resp authenticationrequests.RegisterRequest
	err = json.Unmarshal([]byte(res), &resp)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to resend register otp", service, err)
		return errs.RegisterFailed
	}
	nationalCode := resp.NationalIdentityCode
	ctx = context.WithValue(ctx, consts.NationalIdentityCode, nationalCode)

	user, err := service.UserService.GetByNationalIdentityCode(ctx, resp.NationalIdentityCode)
	if err != nil && !utils.CheckError(err, errs.RecordNotFound) {
		logger.LogErrorWithFieldsV2(ctx, "failed to get user by national", service, err)
		return errs.RegisterFailed
	}

	if user != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get user by national", service, nil)
		return errs.ErrAuthenticationFailed
	}

	err = service.OTPService.RequestOTP(ctx, resp.Mobile)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to resend register otp", service, err,
			zap.String("mobile", resp.Mobile),
		)
		return err
	}

	return nil
}
