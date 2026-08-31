package services

import (
	"athena/src/api/errs"
	"athena/src/cache"
	"athena/src/config"
	"athena/src/pkg/logger"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type IOTPService interface {
	RequestOTP(ctx context.Context, mobile string) error
	VerifyOTP(ctx context.Context, mobile, otp string) (bool, error)
	generateOTP(length int) (string, error)
}

type OTPService struct {
	//NotificationService *NotificationProcessService
}

func (service *OTPService) RequestOTP(ctx context.Context, mobile string) error {
	whiteListMobiles := strings.Split(config.GetInstance().Get("WHITE_LIST_MOBILES"), ",")
	for _, validMobile := range whiteListMobiles {
		if mobile == validMobile {
			return nil
		}
	}

	// get a key for the otp
	key := getRedisKey(mobile)

	// get if otp exist don't let new otp be create
	exists, err := cache.GetInstance().GetClient().Exists(ctx, key).Result()
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get key from cache", service, err,
			zap.String("mobile", mobile),
		)
		return errs.SomeThingWentWrong
	}

	if exists != 0 {
		logger.LogErrorWithFieldsV2(ctx, "otp already exists", service, nil,
			zap.String("mobile", mobile),
		)
		return errs.ErrAuthOTPExists
	}

	otpLength, err := strconv.Atoi(config.GetInstance().Get("OTP_LENGTH"))
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to convert OTP string to int", service, err)
		return errs.SomeThingWentWrong
	}

	// generate a otp
	otpExpiration, err := time.ParseDuration(config.GetInstance().Get("OTP_EXPIRATION") + "s")
	otp, err := service.generateOTP(otpLength)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "generate OTP failed", service, err,
			zap.String("mobile", mobile),
		)
		return errs.SomeThingWentWrong
	}
	//Convert the format to international
	//InternationalMobile := "+98" + mobile[1:]
	//// Send otp
	//notification := drivers.NewSmsNotification().
	//	SetDestinations([]string{InternationalMobile}).
	//	SetDataMethod("otp").
	//	SetTemplateId("8").
	//	AddParameter("otp", otp)
	//
	//// send otp to user
	//err = service.NotificationService.Send(ctx, notification)
	//if err != nil {
	//	logger.LogErrorWithFieldsV2(ctx, "failed to send otp", service, err,
	//		zap.String("mobile", mobile),
	//	)
	//	return errs.SomeThingWentWrong
	//}

	// Send otp
	fmt.Println("the otp is:", otp)

	// set key:otp in redis
	err = cache.GetInstance().GetClient().Set(ctx, key, otp, otpExpiration).Err()
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to set otp in cache", service, err,
			zap.String("mobile", mobile),
		)
		return errs.SomeThingWentWrong
	}

	return nil
}

func (service *OTPService) VerifyOTP(ctx context.Context, mobile, otp string) (bool, error) {

	// get a key for the otp
	key := getRedisKey(mobile)

	// get the value from redis
	storedOTP, err := cache.GetInstance().GetClient().Get(ctx, key).Result()

	whiteListMobiles := strings.Split(config.GetInstance().Get("WHITE_LIST_MOBILES"), ",")
	FixedOtpValue := config.GetInstance().Get("FIXED_OTP")
	for _, validMobile := range whiteListMobiles {
		if mobile == validMobile {
			if FixedOtpValue == otp {
				return true, nil
			}
		}
	}
	// check otp exist
	if errors.Is(err, redis.Nil) {
		logger.LogErrorWithFieldsV2(ctx, "otp not exist", service, err)
		return false, nil
	} else if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get otp from cache", service, err)
		return false, errs.SomeThingWentWrong
	}

	// check otp value is valid
	if storedOTP != otp {
		logger.LogErrorWithFieldsV2(ctx, "failed to verify otp", service, nil,
			zap.String("mobile", mobile),
			zap.String("otp", otp),
			zap.String("stored-otp", storedOTP),
		)
		return false, nil
	}

	// remove otp in redis if it's ok
	err = cache.GetInstance().GetClient().Del(ctx, key).Err()
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to delete otp from cache", service, err,
			zap.String("mobile", mobile),
		)
		return false, errs.SomeThingWentWrong
	}

	return true, nil
}

func (service *OTPService) generateOTP(length int) (string, error) {
	const charset = "0123456789"
	otp := make([]byte, length)
	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		otp[i] = charset[num.Int64()]
	}
	return string(otp), nil
}

func getRedisKey(mobile string) string {
	return fmt.Sprintf("otp-%s", mobile)
}
