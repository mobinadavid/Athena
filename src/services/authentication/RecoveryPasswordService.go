package authentication

import (
	"athena/src/api/errs"
	authenticationRequests "athena/src/api/http/requests/authentication"
	"athena/src/cache"
	"athena/src/config"
	"athena/src/hash"
	"athena/src/models/consts"
	"athena/src/pkg/logger"
	"athena/src/services"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/docker/distribution/uuid"
	"github.com/redis/go-redis/v9"
)

type IRecoveryPasswordService interface {
	RecoveryPasswordRequestOTP(ctx context.Context, mobile, nationalIdentityCode, owner string) error
	RecoveryPasswordViaOTP(ctx context.Context, otp, newPassword, newPasswordConfirmation, nationalIdentityCode, owner string) error
	SaveStateAndSendOTP(ctx context.Context) (string, error)
	VerifyRegisterOTPViaRedisKey(ctx context.Context) error
	SetPassword(ctx context.Context) error
	ResendRegisterOTP(ctx context.Context) error
}

type RecoveryPasswordService struct {
	UserService services.IUserService
	OTPService  services.IOTPService
}

func (service *RecoveryPasswordService) RecoveryPasswordRequestOTP(ctx context.Context, mobile, nationalIdentityCode, owner string) error {
	if owner == "user" {
		user, err := service.UserService.GetByNationalIdentityCode(ctx, nationalIdentityCode)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to get user by national identity code", service, err)
			return errs.ErrAuthenticationFailed
		}
		if user.Mobile != mobile {
			logger.LogErrorWithFieldsV2(ctx, "mobile is not valid", service, nil)
			return errs.ErrAuthenticationFailed
		}
		err = service.OTPService.RequestOTP(ctx, mobile)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to request otp", service, err,
				zap.String("mobile", mobile),
			)
			return err
		}
	}

	return nil
}

func (service *RecoveryPasswordService) RecoveryPasswordViaOTP(ctx context.Context, otp, newPassword, newPasswordConfirmation, nationalIdentityCode, owner string) error {
	if owner == "user" {
		user, err := service.UserService.GetByNationalIdentityCode(ctx, nationalIdentityCode)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to get user by national identity code", service, err)
			return errs.SomeThingWentWrong
		}

		otpIsValid, err := service.OTPService.VerifyOTP(ctx, user.Mobile, otp)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed verify otp", service, err,
				zap.String("mobile", user.Mobile),
				zap.String("otp", otp),
			)
			return errs.SomeThingWentWrong
		}

		if !otpIsValid {
			logger.LogErrorWithFieldsV2(ctx, "otp is not valid", service, nil,
				zap.String("mobile", user.Mobile),
				zap.String("otp", otp),
			)
			return errs.ErrOTPInvalid
		}

		if newPassword != newPasswordConfirmation {
			logger.LogErrorWithFieldsV2(ctx, "password confirmation does not match", service, nil)
			return errs.PasswordNotMatch
		}

		// Update User Password.
		user.Password = []byte(newPassword)
		_, err = service.UserService.Update(ctx, user)
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to update user password", service, err)
			return errs.SomeThingWentWrong
		}

	}

	return nil
}

func (service *RecoveryPasswordService) SaveStateAndSendOTP(ctx context.Context) (string, error) {
	// fetch data
	req, ok := ctx.Value("req").(*authenticationRequests.RecoverPasswordRequest)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to save the recover password state and send otp", service, nil)
		return "", errs.RecoverPasswordFailed
	}

	user, err := service.UserService.GetByNationalIdentityCode(ctx, req.NationalIdentityCode)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to save the recover password state and send otp", service, err)
		return "", errs.ErrAuthenticationFailed
	}

	if user.Mobile != req.Mobile {
		logger.LogErrorWithFieldsV2(ctx, "mobile is not valid", service, nil,
			zap.String("mobile", req.Mobile),
		)
		return "", errs.ErrAuthenticationFailed
	}

	// marshal the req to save in redis
	reqData, err := json.Marshal(req)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to save the recover password state and send otp", service, err)
		return "", errs.RecoverPasswordFailed
	}

	// get expire time
	expiration, err := strconv.Atoi(config.GetInstance().Get("RECOVER_PASSWORD_LIFETIME"))
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to save the recover password state and send otp", service, err)
		expiration = 1200
	}

	// Save the request data in Redis
	key := uuid.Generate().String()
	err = cache.GetInstance().GetClient().Set(ctx, key, reqData, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to save the recover password state and send otp", service, err)
		return "", errs.RecoverPasswordFailed
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

func (service *RecoveryPasswordService) VerifyRegisterOTPViaRedisKey(ctx context.Context) error {
	req, ok := ctx.Value("req").(*authenticationRequests.VerifyRecoverPasswordOTP)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify recover password otp", service, nil)
		return errs.RecoverPasswordFailed
	}

	res, err := cache.GetInstance().GetClient().Get(context.Background(), req.RecoverPasswordKey).Result()
	if err != nil {
		if err == redis.Nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to verify recover password otp", service, err)
			return errs.ErrRecoverPasswordTimeOut
		}

		logger.LogErrorWithFieldsV2(ctx, "failed to verify recover password otp", service, err)
		return errs.RecoverPasswordFailed
	}

	var resp authenticationRequests.RecoverPasswordRequest
	err = json.Unmarshal([]byte(res), &resp)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify recover password otp", service, err)
		return errs.RecoverPasswordFailed
	}
	nationalCode := resp.NationalIdentityCode
	ctx = context.WithValue(ctx, consts.NationalIdentityCode, nationalCode)

	var otpIsValid bool
	otpIsValid, err = service.OTPService.VerifyOTP(ctx, resp.Mobile, req.OTP)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify recover password otp", service, err,
			zap.String("mobile", resp.Mobile),
			zap.String("otp", req.OTP),
		)
		return errs.RecoverPasswordFailed
	}

	if !otpIsValid {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify recover password otp", service, nil,
			zap.String("mobile", resp.Mobile),
			zap.String("otp", req.OTP),
		)
		return errs.ErrOTPInvalid
	}

	return nil
}

func (service *RecoveryPasswordService) SetPassword(ctx context.Context) error {
	req, ok := ctx.Value("req").(*authenticationRequests.SetRecoverPasswordRequest)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to recover password", service, nil)
		return errs.RecoverPasswordFailed
	}

	res, err := cache.GetInstance().GetClient().Get(context.Background(), req.RecoverPasswordKey).Result()
	if err != nil {
		if err == redis.Nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to recover password", service, err)
			return errs.ErrRecoverPasswordTimeOut
		}

		logger.LogErrorWithFieldsV2(ctx, "failed to recover password", service, err)
		return errs.RecoverPasswordFailed
	}

	var resp authenticationRequests.RegisterRequest
	err = json.Unmarshal([]byte(res), &resp)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to recover password", service, err)
		return errs.RecoverPasswordFailed
	}
	nationalCode := resp.NationalIdentityCode
	ctx = context.WithValue(ctx, consts.NationalIdentityCode, nationalCode)

	user, err := service.UserService.GetByNationalIdentityCode(ctx, resp.NationalIdentityCode)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to save the recover password state and send otp", service, err)
		return errs.RecoverPasswordFailed
	}

	// check password
	if req.Password != req.RePassword {
		logger.LogErrorWithFieldsV2(ctx, "failed to recover password", service, nil)
		return errs.PasswordNotMatch
	}
	// Validate Given Current Password
	passwordHashCheck, err := hash.VerifyStoredHash(user.Password, req.Password)
	if err != nil || passwordHashCheck {
		logger.LogErrorWithFieldsV2(ctx, "same password", service, err)
		return errs.PasswordShouldBeNew
	}

	user.Password = []byte(req.Password)
	_, err = service.UserService.Update(ctx, user)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to recover password", service, err)
		return errs.RecoverPasswordFailed
	}

	return nil
}

func (service *RecoveryPasswordService) ResendRegisterOTP(ctx context.Context) error {
	req, ok := ctx.Value("req").(*authenticationRequests.ResendRecoverPasswordOTP)
	if !ok {
		logger.LogErrorWithFieldsV2(ctx, "failed to resend recover password otp", service, nil)
		return errs.FailedToSendOTP
	}

	res, err := cache.GetInstance().GetClient().Get(context.Background(), req.RecoverPasswordKey).Result()
	if err != nil {
		if err == redis.Nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to resend recover password otp", service, err)
			return errs.ErrRecoverPasswordTimeOut
		}

		logger.LogErrorWithFieldsV2(ctx, "failed to resend recover password otp", service, err)
		return errs.FailedToSendOTP
	}

	var resp authenticationRequests.RecoverPasswordRequest
	err = json.Unmarshal([]byte(res), &resp)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to resend recover password otp", service, err)
		return errs.FailedToSendOTP
	}
	nationalCode := resp.NationalIdentityCode
	ctx = context.WithValue(ctx, consts.NationalIdentityCode, nationalCode)

	_, err = service.UserService.GetByNationalIdentityCode(ctx, resp.NationalIdentityCode)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to resend recover password otp", service, err)
		return errs.FailedToSendOTP
	}

	err = service.OTPService.RequestOTP(ctx, resp.Mobile)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to resend recover password otp", service, err,
			zap.String("mobile", resp.Mobile),
		)
		return err
	}

	return nil
}
