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
