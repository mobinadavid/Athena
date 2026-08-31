package authentication

import (
	"athena/src/api/errs"
	authRequests "athena/src/api/http/requests/authentication"
	"athena/src/api/http/response"
	"athena/src/config"
	"athena/src/models"
	"athena/src/models/consts"
	"athena/src/pkg/logger"
	"athena/src/pkg/utils"
	"athena/src/pkg/validator"
	"athena/src/services/authentication"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RegisterController struct {
	RegisterService authentication.IRegisterService
}

func (controller *RegisterController) Register(c *gin.Context) {
	// Bind check payload.
	var req authRequests.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogJSONBindError(c, err)
		response.Api(c).SetLog().Send()
		return
	}

	// validate the payload.
	if err := validator.Validate(&req, c.GetString("locale")); err != nil {
		logger.LogValidationError(c, err)
		response.Api(c).SetErrors(err).SetLog().Send()
		return
	}

	ctx := context.WithValue(context.Background(), "req", &req)
	nationalCode := req.NationalIdentityCode
	ctx = context.WithValue(ctx, consts.NationalIdentityCode, nationalCode)
	ctx = context.WithValue(ctx, consts.OwnerType, models.UserRole)
	ctx = context.WithValue(ctx, consts.RequestUuid, c.GetString("request-uuid"))

	key, err := controller.RegisterService.SaveStateAndSendOTP(ctx)
	if err != nil {
		logger.LogServiceV2(c, "failed to register user", controller, err)

		if utils.CheckError(err, errs.ErrAuthOTPExists) {
			response.Api(c).SetMessage(err.Error()).SetErrorCode(errs.OTPAlreadyExistErrorCode).SetLog().Send()
			return
		}

		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// get expire time
	expiration, err := strconv.Atoi(config.GetInstance().Get("REGISTER_SAVE_STATE_LIFETIME"))
	if err != nil {
		logger.LogAToIError(c, err)
		expiration = 120
	}

	// get cookie name
	cookieName := config.GetInstance().Get("REGISTER_COOKIE_KEY_NAME")
	if cookieName == "" {
		cookieName = "register_key"
	}

	// init cookie
	expirationTime := time.Now().Add(time.Duration(expiration) * time.Second)
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    key,
		Path:     "/",
		Domain:   config.GetInstance().Get("COOKIE_HOST"),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  expirationTime,
	}

	// set cookie
	http.SetCookie(c.Writer, cookie)
	c.Header("Set-Cookie", cookie.String())

	// Return response.
	response.Api(c).SetMessage("register-request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]any{
			"key":      key,
			"time-out": expiration,
		}).SetLog().Send()
}

func (controller *RegisterController) VerifyRegister(c *gin.Context) {
	// Bind check payload.
	var req authRequests.VerifyRegisterOTP
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogJSONBindError(c, err)
		response.Api(c).SetLog().Send()
		return
	}

	// validate the payload.
	if err := validator.Validate(&req, c.GetString("locale")); err != nil {
		logger.LogValidationError(c, err)
		response.Api(c).SetErrors(err).SetLog().Send()
		return
	}
	// get key from cookie
	if req.RegisterKey == "" {
		key, cookieErr := c.Cookie(config.GetInstance().Get("REGISTER_COOKIE_KEY_NAME"))
		if cookieErr != nil {
			logger.LogServiceV2(c, "failed to get register_key from cookie", controller, cookieErr)
			response.Api(c).SetStatusCode(http.StatusUnauthorized).SetMessage(errs.SomeThingWentWrong.Error()).SetLog().Send()
			return
		}
		req.RegisterKey = key
	}
	// register user
	ctx := context.WithValue(context.Background(), "req", &req)
	ctx = context.WithValue(ctx, consts.OwnerType, models.UserRole)
	ctx = context.WithValue(ctx, consts.RequestUuid, c.GetString("request-uuid"))

	err := controller.RegisterService.VerifyRegisterOTPViaRedisKey(ctx)
	if err != nil {
		logger.LogServiceV2(c, "failed to register user", controller, err)

		if utils.CheckError(err, errs.ErrAuthOTPExists) {
			response.Api(c).SetMessage(err.Error()).SetErrorCode(errs.OTPAlreadyExistErrorCode).SetLog().Send()
			return
		}

		if utils.CheckError(err, errs.RegisterFailed) {
			response.Api(c).SetMessage(err.Error()).SetErrorCode(errs.RegisterTimeOutErrorCode).SetLog().Send()
			return
		}

		resp := response.Api(c).SetMessage(err.Error())
		resp.SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).SetMessage("verify-register-request-successful").
		SetStatusCode(http.StatusOK).
		SetLog().
		Send()
}

func (controller *RegisterController) ResendOTP(c *gin.Context) {
	// Bind check payload.
	var req authRequests.ResendRegisterOTP
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogJSONBindError(c, err)
		response.Api(c).SetLog().Send()
		return
	}

	// validate the payload.
	if err := validator.Validate(&req, c.GetString("locale")); err != nil {
		logger.LogValidationError(c, err)
		response.Api(c).SetErrors(err).SetLog().Send()
		return
	}

	// get cookie name
	cookieName := config.GetInstance().Get("REGISTER_COOKIE_KEY_NAME")
	if cookieName == "" {
		cookieName = "register_key"
	}

	// check the key is present or not
	if req.RegisterKey == "" {
		// Read the cookie set in Step 1
		key, err := c.Cookie(cookieName)
		if err != nil {
			// in this case cookie is expired
			logger.LogCookieDoesNotExist(c, err, cookieName)
			response.Api(c).SetMessage(errs.RegisterFailed.Error()).SetErrorCode(errs.RegisterTimeOutErrorCode).SetLog().Send()
			return
		}

		req.RegisterKey = key
	}

	// prepare data for service
	ctx := context.WithValue(context.Background(), "req", &req)
	ctx = context.WithValue(ctx, consts.OwnerType, models.UserRole)
	ctx = context.WithValue(ctx, consts.RequestUuid, c.GetString("request-uuid"))

	// register user
	err := controller.RegisterService.ResendRegisterOTP(ctx)
	if err != nil {
		logger.LogServiceV2(c, "failed to resend register otp", controller, err)

		if utils.CheckError(err, errs.ErrAuthOTPExists) {
			response.Api(c).SetMessage(err.Error()).SetErrorCode(errs.OTPAlreadyExistErrorCode).SetLog().Send()
			return
		}

		if utils.CheckError(err, errs.RegisterFailed) {
			response.Api(c).SetMessage(err.Error()).SetErrorCode(errs.RegisterTimeOutErrorCode).SetLog().Send()
			return
		}

		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]any{
			"key": req.RegisterKey,
		}).SetLog().Send()
}
