package authentication

import (
	"context"
	"net/http"

	"athena/src/api/errs"
	"athena/src/api/http/requests/TwoFactorAuthentication"
	"athena/src/api/http/response"
	"athena/src/models"
	"athena/src/models/consts"
	"athena/src/pkg/logger"
	"athena/src/pkg/utils"
	"athena/src/pkg/validator"
	"athena/src/services/authentication"

	"github.com/gin-gonic/gin"
)

type TwoFAController struct {
	TwoFaService authentication.ITwoFaService
}

func twoFaContext(c *gin.Context) (context.Context, *utils.Auth) {
	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerType, models.AdminRole)
	ctx = context.WithValue(ctx, consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)
	return ctx, auth
}

func (controller *TwoFAController) Status(c *gin.Context) {
	ctx, auth := twoFaContext(c)
	enabled, err := controller.TwoFaService.CheckTwoFAEnable(ctx, auth.OwnerId)
	if err != nil {
		logger.LogServiceV2(c, "failed to get two factor status", controller, err)
		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"two_fa_enabled": enabled,
		}).
		SetLog().
		Send()
}

func (controller *TwoFAController) Enable(c *gin.Context) {
	ctx, auth := twoFaContext(c)
	secretURL, secretKey, qrCode, err := controller.TwoFaService.Enable(ctx, auth.OwnerId)
	if err != nil {
		logger.LogServiceV2(c, "failed to enable two factor authentication", controller, err)
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).SetMessage(err.Error()).SetLog().Send()
		return
	}

	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"secret_url": secretURL,
			"secret_key": secretKey,
			"qr_code":    qrCode,
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

	ctx, auth := twoFaContext(c)
	recoverCodes, err := controller.TwoFaService.VerifyCodeFirstTime(ctx, auth.OwnerId, req.Totp)
	if err != nil {
		logger.LogServiceV2(c, "failed to verify two factor code", controller, err)
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).SetMessage(err.Error()).SetLog().Send()
		return
	}

	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"recovery_codes": recoverCodes,
		}).SetLog().Send()
}

func (controller *TwoFAController) Disable(c *gin.Context) {
	var req TwoFactorAuthentication.DisableTwoFARequest
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

	if req.Totp == "" && req.RecoveryCode == "" {
		response.Api(c).SetMessage(errs.ErrTwoFactorChallengeMissing.Error()).SetLog().Send()
		return
	}

	ctx, auth := twoFaContext(c)
	err := controller.TwoFaService.Disable(ctx, auth.OwnerId, req.Totp, req.RecoveryCode)
	if err != nil {
		logger.LogServiceV2(c, "failed to disable two factor authentication", controller, err)
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).SetMessage(err.Error()).SetLog().Send()
		return
	}

	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetLog().
		Send()
}
