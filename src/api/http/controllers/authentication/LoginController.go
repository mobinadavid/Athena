package authentication

import (
	"athena/src/api/errs"
	authenticationrequests "athena/src/api/http/requests/authentication"
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

type LoginController struct {
	LoginService authentication.ILoginService
}

func (controller *LoginController) LoginViaPassword(c *gin.Context) {
	// bind the incoming request to json
	var req authenticationrequests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogJSONBindError(c, err)
		response.Api(c).SetLog().Send()
		return
	}

	// validate request
	if err := validator.Validate(&req, c.GetString("locale")); err != nil {
		logger.LogValidationError(c, err)
		response.Api(c).SetErrors(err).SetLog().Send()
		return
	}

	ctx := context.WithValue(context.Background(), "req", &req)
	ctx = context.WithValue(ctx, "request-ip", c.GetString("request-ip"))
	ctx = context.WithValue(ctx, "request-user-agent", c.GetHeader("User-Agent"))
	nationalCode := req.NationalIdentityCode
	if nationalCode != "" {
		ctx = context.WithValue(ctx, consts.NationalIdentityCode, nationalCode)
	}
	ctx = context.WithValue(ctx, consts.RequestUuid, c.GetString("request-uuid"))
	ctx = context.WithValue(ctx, consts.OwnerType, models.UserRole)

	jwtDTO, err := controller.LoginService.UserLoginViaPassword(ctx)
	if err != nil {
		logger.LogServiceV2(c, "failed to login via password (user)", controller, err)
		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	expSeconds, err := strconv.Atoi(config.GetInstance().Get("COOKIE_EXPIRATION"))
	if err != nil {
		logger.LogAToIError(c, err)
		response.Api(c).SetMessage(errs.SomeThingWentWrong.Error()).SetLog().Send()
		return
	}

	expiration := time.Now().Add(time.Duration(expSeconds) * time.Second)

	cookie := &http.Cookie{
		Name:     config.GetInstance().Get("COOKIE_NAME"),
		Value:    jwtDTO.RefreshTokenString,
		Path:     "/",
		Domain:   config.GetInstance().Get("COOKIE_HOST"),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  expiration,
	}

	http.SetCookie(c.Writer, cookie)
	c.Header("Set-Cookie", cookie.String())

	response.Api(c).
		SetMessage("welcome").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"access_tokens": jwtDTO,
		}).
		SetLog().Send()
}

func (controller *LoginController) LoginViaOtpSendOtp(c *gin.Context) {
	// Bind check payload.
	var req authenticationrequests.LoginViaOTPSendOtpRequest
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
	ctx = context.WithValue(ctx, consts.RequestUuid, c.GetString("request-uuid"))
	ctx = context.WithValue(ctx, consts.OwnerType, models.UserRole)
	key, err := controller.LoginService.LoginViaOtpSendOtp(ctx)
	if err != nil {
		logger.LogServiceV2(c, "failed to login via otp", controller, err)

		if utils.CheckError(err, errs.ErrAuthOTPExists) {
			response.Api(c).SetMessage(err.Error()).SetErrorCode(errs.OTPAlreadyExistErrorCode).SetLog().Send()
			return
		}

		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}
	// get expire time
	expiration, err := strconv.Atoi(config.GetInstance().Get("LOGIN_SAVE_STATE_LIFETIME"))
	if err != nil {
		logger.LogAToIError(c, err)
		expiration = 300
	}

	// get cookie name
	cookieName := config.GetInstance().Get("LOGIN_COOKIE_KEY_NAME")
	if cookieName == "" {
		cookieName = "login_key"
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

func (controller *LoginController) ResendOTP(c *gin.Context) {
	// Bind check payload.
	var req authenticationrequests.ResendLoginOTP
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
	cookieName := config.GetInstance().Get("LOGIN_COOKIE_KEY_NAME")
	if cookieName == "" {
		cookieName = "login_key"
	}

	// check the key is present or not
	if req.LoginKey == "" {
		// Read the cookie set in Step 1
		key, err := c.Cookie(cookieName)
		if err != nil {
			// in this case cookie is expired
			logger.LogCookieDoesNotExist(c, err, cookieName)
			response.Api(c).SetMessage(errs.ErrLoginTimeOut.Error()).SetErrorCode(errs.RegisterTimeOutErrorCode).SetLog().Send()
			return
		}

		req.LoginKey = key
	}

	// prepare data for service
	ctx := context.WithValue(context.Background(), "req", &req)
	ctx = context.WithValue(ctx, consts.OwnerType, models.UserRole)
	ctx = context.WithValue(ctx, consts.RequestUuid, c.GetString("request-uuid"))

	// register user
	err := controller.LoginService.ResendLoginOTP(ctx)
	if err != nil {
		logger.LogServiceV2(c, "failed to resend register otp", controller, err)

		if utils.CheckError(err, errs.ErrAuthOTPExists) {
			response.Api(c).SetMessage(err.Error()).SetErrorCode(errs.OTPAlreadyExistErrorCode).SetLog().Send()
			return
		}

		if utils.CheckError(err, errs.ErrAuthenticationFailed) {
			response.Api(c).SetMessage(err.Error()).SetErrorCode(errs.LoginTimeOutErrorCode).SetLog().Send()
			return
		}

		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]any{
			"key": req.LoginKey,
		}).SetLog().Send()
}

func (controller *LoginController) VerifyOTP(c *gin.Context) {
	var req authenticationrequests.VerifyLoginOTP
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogJSONBindError(c, err)
		response.Api(c).SetLog().Send()
		return
	}

	// validate request
	if err := validator.Validate(&req, c.GetString("locale")); err != nil {
		logger.LogValidationError(c, err)
		response.Api(c).SetErrors(err).Send()
		return
	}

	// get cookie name
	cookieName := config.GetInstance().Get("LOGIN_COOKIE_KEY_NAME")
	if cookieName == "" {
		cookieName = "login_key"
	}

	// check the key is present or not
	if req.LoginKey == "" {
		// Read the cookie set in Step 1
		key, err := c.Cookie(cookieName)
		if err != nil {
			// in this case cookie is expired
			logger.LogCookieDoesNotExist(c, err, cookieName)
			response.Api(c).SetMessage(errs.ErrLoginTimeOut.Error()).SetErrorCode(errs.LoginTimeOutErrorCode).SetLog().Send()
			return
		}

		req.LoginKey = key
	}

	// prepare data for service
	ctx := context.WithValue(context.Background(), "req", &req)
	ctx = context.WithValue(ctx, "request-ip", c.GetString("request-ip"))
	ctx = context.WithValue(ctx, "request-user-agent", c.GetHeader("User-Agent"))
	ctx = context.WithValue(ctx, consts.OwnerType, models.UserRole)
	ctx = context.WithValue(ctx, consts.RequestUuid, c.GetString("request-uuid"))

	jwt, err := controller.LoginService.UserLoginVerifyOTP(ctx)
	if err != nil {
		logger.LogServiceV2(c, "failed to get access token via otp", controller, err)

		if utils.CheckError(err, errs.ErrLoginTimeOut) {
			response.Api(c).SetMessage(err.Error()).SetErrorCode(errs.LoginTimeOutErrorCode).SetLog().Send()
			return
		}

		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	expirationSeconds, err := strconv.Atoi(config.GetInstance().Get("COOKIE_EXPIRATION"))
	if err != nil {
		logger.LogAToIError(c, err)
		response.Api(c).SetMessage(errs.SomeThingWentWrong.Error()).SetLog().Send()
		return
	}

	expiration := time.Now().Add(time.Duration(expirationSeconds) * time.Second)
	cookie := &http.Cookie{
		Name:     config.GetInstance().Get("COOKIE_NAME"),
		Value:    jwt.RefreshTokenString,
		Path:     "/",
		Domain:   config.GetInstance().Get("COOKIE_HOST"),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  expiration,
	}

	// Set the cookie using http.SetCookie (this should automatically add the cookie to headers)
	http.SetCookie(c.Writer, cookie)
	c.Header("Set-Cookie", cookie.String())

	// Return response.
	response.Api(c).SetMessage("welcome").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"access_tokens": jwt,
		}).SetLog().Send()
}
