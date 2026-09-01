package authentication

import (
	"athena/src/api/errs"
	"athena/src/api/http/requests/TwoFactorAuthentication"
	"athena/src/api/http/response"
	"athena/src/models/consts"
	"athena/src/pkg/logger"
	"athena/src/pkg/utils"
	"athena/src/pkg/validator"
	"athena/src/services/authentication"
	"context"

	"net/http"

	"github.com/gin-gonic/gin"
)

type TwoFAController struct {
	TwoFaService authentication.ITwoFaService
}

// SendEnable2FaOTP godoc
// @Summary      send enable 2fa otp
// @Description  Send an otp required to start enabling two factor authentication for the authenticated user
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Success      200      {object}  response.Response
// @Security      BearerAuth
// @Security      CookieAuth
// @Router       /profile/2fa/enable/send-otp [post]
func (controller *TwoFAController) SendEnable2FaOTP(c *gin.Context) {
	userID := c.GetUint("authenticated-user-id")

	// getting user by id
	ctx := context.WithValue(context.Background(), "user_id", userID)
	auth := utils.GetAuthData(c)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)
	err := controller.TwoFaService.SendEnable2FAOTP(ctx)
	if err != nil {
		logger.LogServiceV2(c, "failed to send otp", controller, err)

		if utils.CheckError(err, errs.ErrAuthOTPExists) {
			response.Api(c).SetMessage(err.Error()).SetErrorCode(errs.OTPAlreadyExistErrorCode).SetLog().Send()
			return
		}

		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// send response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetLog().
		Send()
}

// GetSecretKey godoc
// @Summary      get 2fa secret key
// @Description  Verify the enable otp and return the TOTP secret key and url to set up two factor authentication
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        otp  path  string  true  "One-time password sent to enable 2fa"
// @Success      200      {object}  response.Response
// @Security      BearerAuth
// @Security      CookieAuth
// @Router       /profile/2fa/enable/{otp} [get]
func (controller *TwoFAController) GetSecretKey(c *gin.Context) {
	userId := c.GetUint("authenticated-user-id")
	otp := c.Param("otp")
	if otp == "" {
		logger.LogValueNotExist(c, "otp")
		response.Api(c).SetMessage(errs.RecordNotFound.Error()).SetLog().Send()
		return
	}
	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)
	secretUrl, secretKey, err := controller.TwoFaService.GetSecretKey(ctx, userId, otp)
	if err != nil {
		logger.LogServiceV2(c, "failed to get secret key for two factor authentication", controller, err)
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"secret_url": secretUrl,
			"secret_key": secretKey,
		}).
		SetLog().Send()
}

// VerifyCode godoc
// @Summary      verify 2fa code
// @Description  Verify the TOTP code for the first time to enable two factor authentication and return recovery codes
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request  body      TwoFactorAuthentication.VerifyTotpRequest  true  "verify totp request body"
// @Success      200      {object}  response.Response
// @Security      BearerAuth
// @Security      CookieAuth
// @Router       /profile/2fa/enable/verify [post]
func (controller *TwoFAController) VerifyCode(c *gin.Context) {
	var req TwoFactorAuthentication.VerifyTotpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogJSONBindError(c, err)
		response.Api(c).SetLog().Send()
		return
	}

	if err := validator.Validate(&req, c.GetString("locale")); err != nil {
		logger.LogValidationError(c, err)
		response.Api(c).SetErrors(err).SetLog().Send()
		return
	}

	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)
	recoverCodes, err := controller.TwoFaService.VerifyCodeFirstTime(ctx, auth.OwnerId, req.Totp)
	if err != nil {
		logger.LogServiceV2(c, "failed to verify code", controller, err)
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"recovery_codes": recoverCodes,
		}).SetLog().Send()
}

// Disable godoc
// @Summary      disable 2fa
// @Description  Disable two factor authentication for the authenticated user
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Success      200      {object}  response.Response
// @Security      BearerAuth
// @Security      CookieAuth
// @Router       /profile/2fa/disable [post]
func (controller *TwoFAController) Disable(c *gin.Context) {
	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)
	err := controller.TwoFaService.Disable(ctx, auth.OwnerId)
	if err != nil {
		logger.LogServiceV2(c, "failed to disable two factor authentication", controller, err)
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetLog().
		Send()
}
