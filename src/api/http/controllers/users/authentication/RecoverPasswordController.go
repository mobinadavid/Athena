package authentication

import (
	"athena/src/api/errs"
	authentication2 "athena/src/api/http/requests/authentication"
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

type RecoverPasswordController struct {
	RecoverPasswordService authentication.IRecoveryPasswordService
}

func (controller *RecoverPasswordController) RecoverPassword(c *gin.Context) {
	// Bind check payload.
	var req authentication2.RecoverPasswordRequest
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
	key, err := controller.RecoverPasswordService.SaveStateAndSendOTP(ctx)
	if err != nil {
		logger.LogServiceV2(c, "failed to recover user password", controller, err)

		if utils.CheckError(err, errs.ErrAuthOTPExists) {
			response.Api(c).SetMessage(err.Error()).SetErrorCode(errs.OTPAlreadyExistErrorCode).SetLog().Send()
			return
		}

		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// get expire time
	expiration, err := strconv.Atoi(config.GetInstance().Get("RECOVER_PASSWORD_LIFETIME"))
	if err != nil {
		logger.LogAToIError(c, err)
		expiration = 120
	}

	// get cookie name
	cookieName := config.GetInstance().Get("RECOVER_PASSWORD_COOKIE_KEY_NAME")
	if cookieName == "" {
		cookieName = "recover_password_key"
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
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]any{
			"key":      key,
			"time-out": expiration,
		}).SetLog().Send()
}

func (controller *RecoverPasswordController) VerifyOtpRecoverPassword(c *gin.Context) {
	// Bind check payload.
	var req authentication2.VerifyRecoverPasswordOTP
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
	cookieName := config.GetInstance().Get("RECOVER_PASSWORD_COOKIE_KEY_NAME")
	if cookieName == "" {
		cookieName = "recover_password_key"
	}

	// check the key is present or not
	if req.RecoverPasswordKey == "" {
		// Read the cookie set in Step 1
		key, err := c.Cookie(cookieName)
		if err != nil {
			// in this case cookie is expired
			logger.LogCookieDoesNotExist(c, err, cookieName)
			response.Api(c).SetMessage(errs.RecoverPasswordFailed.Error()).SetErrorCode(errs.RecoverPasswordTimeOutErrorCode).SetLog().Send()
			return
		}

		req.RecoverPasswordKey = key
	}

	// prepare data for service
	ctx := context.WithValue(context.Background(), "req", &req)
	ctx = context.WithValue(ctx, consts.OwnerType, models.UserRole)
	ctx = context.WithValue(ctx, consts.RequestUuid, c.GetString("request-uuid"))

	err := controller.RecoverPasswordService.VerifyRegisterOTPViaRedisKey(ctx)
	if err != nil {
		logger.LogServiceV2(c, "failed to register user", controller, err)

		if utils.CheckError(err, errs.ErrRecoverPasswordTimeOut) {
			response.Api(c).SetMessage(err.Error()).SetErrorCode(errs.RecoverPasswordTimeOutErrorCode).SetLog().Send()
			return
		}

		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]any{
			"key": req.RecoverPasswordKey,
		}).SetLog().Send()
}

func (controller *RecoverPasswordController) SetPassword(c *gin.Context) {
	// Bind check payload.
	var req authentication2.SetRecoverPasswordRequest
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
	cookieName := config.GetInstance().Get("RECOVER_PASSWORD_COOKIE_KEY_NAME")
	if cookieName == "" {
		cookieName = "recover_password_key"
	}

	// check the key is present or not
	if req.RecoverPasswordKey == "" {
		// Read the cookie set in Step 1
		key, err := c.Cookie(cookieName)
		if err != nil {
			// in this case cookie is expired
			logger.LogCookieDoesNotExist(c, err, cookieName)
			response.Api(c).SetMessage(errs.RecoverPasswordFailed.Error()).SetErrorCode(errs.RecoverPasswordTimeOutErrorCode).SetLog().Send()
			return
		}

		req.RecoverPasswordKey = key
	}

	// prepare data for service
	ctx := context.WithValue(context.Background(), "req", &req)
	ctx = context.WithValue(ctx, consts.OwnerType, models.UserRole)
	ctx = context.WithValue(ctx, consts.RequestUuid, c.GetString("request-uuid"))

	// register user
	err := controller.RecoverPasswordService.SetPassword(ctx)
	if err != nil {
		logger.LogServiceV2(c, "failed to register user", controller, err)

		if utils.CheckError(err, errs.ErrRecoverPasswordTimeOut) {
			response.Api(c).SetMessage(err.Error()).SetErrorCode(errs.RecoverPasswordTimeOutErrorCode).SetLog().Send()
			return
		}

		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// expire the cookie
	expirationTime := time.Unix(0, 0)
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    "",
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
	response.Api(c).SetMessage("recovery-password-request-successful").
		SetStatusCode(http.StatusOK).
		SetLog().
		Send()
}

func (controller *RecoverPasswordController) ResendOTP(c *gin.Context) {
	// Bind check payload.
	var req authentication2.ResendRecoverPasswordOTP
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
	cookieName := config.GetInstance().Get("RECOVER_PASSWORD_COOKIE_KEY_NAME")
	if cookieName == "" {
		cookieName = "recover_password_key"
	}

	// check the key is present or not
	if req.RecoverPasswordKey == "" {
		// Read the cookie set in Step 1
		key, err := c.Cookie(cookieName)
		if err != nil {
			// in this case cookie is expired
			logger.LogCookieDoesNotExist(c, err, cookieName)
			response.Api(c).SetMessage(errs.RecoverPasswordFailed.Error()).SetErrorCode(errs.RecoverPasswordTimeOutErrorCode).SetLog().Send()
			return
		}

		req.RecoverPasswordKey = key
	}

	// prepare data for service
	ctx := context.WithValue(context.Background(), "req", &req)
	ctx = context.WithValue(ctx, consts.OwnerType, models.UserRole)
	ctx = context.WithValue(ctx, consts.RequestUuid, c.GetString("request-uuid"))

	// register user
	err := controller.RecoverPasswordService.ResendRegisterOTP(ctx)
	if err != nil {
		logger.LogServiceV2(c, "failed to resend register otp", controller, err)

		if utils.CheckError(err, errs.ErrAuthOTPExists) {
			response.Api(c).SetMessage(err.Error()).SetErrorCode(errs.OTPAlreadyExistErrorCode).SetLog().Send()
			return
		}

		if utils.CheckError(err, errs.ErrRecoverPasswordTimeOut) {
			response.Api(c).SetMessage(err.Error()).SetErrorCode(errs.RecoverPasswordTimeOutErrorCode).SetLog().Send()
			return
		}

		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]any{
			"key": req.RecoverPasswordKey,
		}).SetLog().Send()
}
