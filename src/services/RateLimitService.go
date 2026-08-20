package services

import (
	"athena/src/api/http/response"
	"athena/src/cache"
	"athena/src/config"
	"athena/src/pkg/logger"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
	"go.uber.org/zap"
)

type IRateLimitService interface {
	GetLimiter() *limiter.Limiter
	GetExcludedKey() func(string) bool
	KeyFunc(*gin.Context) string
	ExcludedKeyFunc(string) bool
	SetLimiter(limiter *limiter.Limiter) IRateLimitService
	SetKey(keyGetter func(*gin.Context) string) IRateLimitService
	SetExcludedKey(excludedKey func(string) bool) IRateLimitService
	GetOnExceedHandler(c *gin.Context, resetIn string) bool
	PostOnExceedHandler(c *gin.Context, resetIn string) bool
}

type RateLimitService struct {
	Limiter     *limiter.Limiter
	KeyGetter   func(*gin.Context) string
	ExcludedKey func(string) bool
}

func (s *RateLimitService) KeyFunc(ctx *gin.Context) string {
	return s.KeyGetter(ctx)
}

func (s *RateLimitService) ExcludedKeyFunc(key string) bool {
	return s.ExcludedKey(key)
}

func (s *RateLimitService) GetExcludedKey() func(string) bool {
	return s.ExcludedKey
}

func (s *RateLimitService) GetLimiter() *limiter.Limiter {
	return s.Limiter
}

func (s *RateLimitService) SetLimiter(limiter *limiter.Limiter) IRateLimitService {
	s.Limiter = limiter
	return s
}

func (s *RateLimitService) SetKey(keyGetter func(*gin.Context) string) IRateLimitService {
	s.KeyGetter = keyGetter
	return s
}

func (s *RateLimitService) SetExcludedKey(excludedKey func(string) bool) IRateLimitService {
	s.ExcludedKey = excludedKey
	return s
}

func (s *RateLimitService) GetOnExceedHandler(c *gin.Context, resetIn string) bool {
	response.Api(c).SetStatusCode(http.StatusTooManyRequests).
		SetMessage(fmt.Sprintf("لطفا بعد از گذشت %s دقیقه مجددا تلاش کنید", resetIn)).SetLog().Send()
	return false
}

func (s *RateLimitService) PostOnExceedHandler(c *gin.Context, resetIn string) bool {
	// read the body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.GetInstance().Error("failed to read request body",
			zap.String("service", reflect.TypeOf(s).Elem().Name()),
			zap.String("reason", "RateLimitService() -> io.ReadAll() return an error"),
			zap.Error(err),
			zap.Time("timestamp", time.Now()),
		)
		response.Api(c).SetStatusCode(http.StatusInternalServerError).SetLog().Send()
		return false
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	response.Api(c).SetStatusCode(http.StatusTooManyRequests).
		SetMessage(fmt.Sprintf("لطفا بعد از گذشت %s دقیقه مجددا تلاش کنید", resetIn)).SetLog().Send()
	return false
}

func DefaultLimiter() *limiter.Limiter {
	strRate := config.GetInstance().Get("RATE_LIMITER_DEFAULT_LIMIT")
	strPeriod := config.GetInstance().Get("RATE_LIMITER_DEFAULT_PERIOD_PER_SECOND")

	if strRate == "" {
		strRate = "60"
	}

	if strPeriod == "" {
		strPeriod = "60"
	}

	rate, err := strconv.Atoi(strRate)
	if err != nil {
		logger.GetInstance().Error("failed to cast string to int",
			zap.String("service", "CriticalLimiter"),
			zap.String("reason", "RATE_LIMITER_DEFAULT_LIMIT should be numeric"),
			zap.Error(err),
			zap.Time("timestamp", time.Now()),
		)
		log.Fatal(err)
	}

	period, err := strconv.Atoi(strPeriod)
	if err != nil {
		logger.GetInstance().Error("failed to cast string to int",
			zap.String("service", "CriticalLimiter"),
			zap.String("reason", "RATE_LIMITER_DEFAULT_PERIOD_PER_SECOND should be numeric"),
			zap.Error(err),
			zap.Time("timestamp", time.Now()),
		)
		log.Fatal(err)
	}

	limiterRate := limiter.Rate{
		Period: time.Duration(period) * time.Second,
		Limit:  int64(rate),
	}

	store, err := sredis.NewStore(cache.GetInstance().GetClient())
	if err != nil {
		logger.GetInstance().Error("failed to create default Limiter",
			zap.String("service", "CriticalLimiter"),
			zap.String("reason", "DefaultLimiter() -> sredis.NewStore() return an error"),
			zap.Error(err),
			zap.Time("timestamp", time.Now()),
		)
		log.Fatal(err)
	}

	return limiter.New(store, limiterRate)
}

func DefaultKeyGetter(c *gin.Context) string {
	return fmt.Sprintf("default-%s", c.GetString("request-ip"))
}

func CriticalLimiter() *limiter.Limiter {
	strRate := config.GetInstance().Get("RATE_LIMITER_CRITICAL_LIMIT")
	strPeriod := config.GetInstance().Get("RATE_LIMITER_CRITICAL_PERIOD_PER_SECOND")

	if strRate == "" {
		strRate = "3"
	}

	if strPeriod == "" {
		strPeriod = "60"
	}

	rate, err := strconv.Atoi(strRate)
	if err != nil {
		logger.GetInstance().Error("failed to cast string to int",
			zap.String("service", "CriticalLimiter"),
			zap.String("reason", "RATE_LIMITER_CRITICAL_LIMIT should be numeric"),
			zap.Error(err),
			zap.Time("timestamp", time.Now()),
		)
		log.Fatal(err)
	}

	period, err := strconv.Atoi(strPeriod)
	if err != nil {
		logger.GetInstance().Error("failed to cast string to int",
			zap.String("service", "CriticalLimiter"),
			zap.String("reason", "RATE_LIMITER_CRITICAL_PERIOD_PER_SECOND should be numeric"),
			zap.Error(err),
			zap.Time("timestamp", time.Now()),
		)
		log.Fatal(err)
	}

	limiterRate := limiter.Rate{
		Period: time.Duration(period) * time.Second,
		Limit:  int64(rate),
	}

	store, err := sredis.NewStore(cache.GetInstance().GetClient())
	if err != nil {
		logger.GetInstance().Error("failed to create critical Limiter",
			zap.String("service", "CriticalLimiter"),
			zap.String("reason", "CriticalLimiter() -> sredis.NewStore() return an error"),
			zap.Error(err),
			zap.Time("timestamp", time.Now()),
		)
		log.Fatal(err)
	}

	return limiter.New(store, limiterRate)
}

func RegisterCriticalLimiter() *limiter.Limiter {
	strRate := config.GetInstance().Get("RATE_LIMITER_REGISTER_CRITICAL_LIMIT")
	strPeriod := config.GetInstance().Get("RATE_LIMITER_REGISTER_CRITICAL_PERIOD_PER_SECOND")
	return SetLimiter(strRate, strPeriod)
}

func LoginCriticalLimiter() *limiter.Limiter {
	strRate := config.GetInstance().Get("RATE_LIMITER_LOGIN_CRITICAL_LIMIT")
	strPeriod := config.GetInstance().Get("RATE_LIMITER_LOGIN_CRITICAL_PERIOD_PER_SECOND")
	return SetLimiter(strRate, strPeriod)
}

func AdminLoginCriticalLimiter() *limiter.Limiter {
	strRate := config.GetInstance().Get("RATE_LIMITER_ADMIN_LOGIN_CRITICAL_LIMIT")
	strPeriod := config.GetInstance().Get("RATE_LIMITER_ADMIN_LOGIN_CRITICAL_PERIOD_PER_SECOND")
	return SetLimiter(strRate, strPeriod)
}

func SetLimiter(strRate, strPeriod string) *limiter.Limiter {
	if strRate == "" {
		strRate = "3"
	}

	if strPeriod == "" {
		strPeriod = "60"
	}

	rate, err := strconv.Atoi(strRate)
	if err != nil {
		logger.GetInstance().Error("failed to cast string to int",
			zap.String("service", "RegisterCriticalLimiter"),
			zap.String("reason", "RATE_LIMITER_REGISTER_CRITICAL_LIMIT should be numeric"),
			zap.Error(err),
			zap.Time("timestamp", time.Now()),
		)
		log.Fatal(err)
	}

	period, err := strconv.Atoi(strPeriod)
	if err != nil {
		logger.GetInstance().Error("failed to cast string to int",
			zap.String("service", "RegisterCriticalLimiter"),
			zap.String("reason", "RATE_LIMITER_REGISTER_CRITICAL_PERIOD_PER_SECOND should be numeric"),
			zap.Error(err),
			zap.Time("timestamp", time.Now()),
		)
		log.Fatal(err)
	}

	limiterRate := limiter.Rate{
		Period: time.Duration(period) * time.Second,
		Limit:  int64(rate),
	}

	store, err := sredis.NewStore(cache.GetInstance().GetClient())
	if err != nil {
		logger.GetInstance().Error("failed to create critical Limiter",
			zap.String("service", "RegisterCriticalLimiter"),
			zap.String("reason", "RegisterCriticalLimiter() -> sredis.NewStore() return an error"),
			zap.Error(err),
			zap.Time("timestamp", time.Now()),
		)
		log.Fatal(err)
	}

	return limiter.New(store, limiterRate)
}

func CriticalKeyGetter(c *gin.Context) string {
	return fmt.Sprintf("critical-%s", c.GetString("request-ip"))
}

func CriticalInvestorShownKeySetter(c *gin.Context) string {
	return fmt.Sprintf("critical-investor-shown-%s", c.GetString("request-ip"))
}

func CriticalChangePasswordKeySetter(c *gin.Context) string {
	return fmt.Sprintf("critical-change-password-%s", c.GetString("request-ip"))
}

func CriticalChangePasswordSendOtpKeySetter(c *gin.Context) string {
	return fmt.Sprintf("critical-change-password-send-otp-%s", c.GetString("request-ip"))
}

func CriticalChangePasswordResendOtpKeySetter(c *gin.Context) string {
	return fmt.Sprintf("critical-change-password-resend-otp-%s", c.GetString("request-ip"))
}

func CriticalStorageKeySetter(c *gin.Context) string {
	return fmt.Sprintf("critical-storage-%s", c.GetString("request-ip"))
}

func CriticalVerifyOtpKeySetter(c *gin.Context) string {
	return fmt.Sprintf("critical-verify-login-%s", c.GetString("request-ip"))
}

func CriticalRecoveryPasswordKEYSetter(c *gin.Context) string {
	return fmt.Sprintf("critical-recovery-password-%s", c.GetString("request-ip"))
}

func CriticalVerifyRecoveryPasswordKEYSetter(c *gin.Context) string {
	return fmt.Sprintf("critical-verify-recovery-password-%s", c.GetString("request-ip"))
}

func GenericCriticalKeyGetter(key string) func(*gin.Context) string {
	return func(c *gin.Context) string {
		return fmt.Sprintf("%s-%s", key, c.GetString("request-ip"))
	}
}
