package authentication

import (
	"athena/src/api/http/requests/AuthenticationRequests/PasswordRecovery"
	"athena/src/api/http/response"
	"athena/src/models"
	"athena/src/models/consts"
	"athena/src/pkg/logger"
	"athena/src/pkg/validator"
	"athena/src/services/authentication"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PasswordRecoveryController struct {
	RecoveryPasswordService authentication.IRecoveryPasswordService
}

func (controller *PasswordRecoveryController) RequestOtp(c *gin.Context) {
	// bind request to json
	var req PasswordRecovery.RequestOtp
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

	nationalCode := req.NationalIdentityCode
	if req.NationalCompanyId != "" {
		nationalCode = req.NationalCompanyId
	}

	ctx := context.WithValue(c.Request.Context(), consts.NationalIdentityCode, nationalCode)
	ctx = context.WithValue(ctx, consts.OwnerType, models.AdminRole)
	ctx = context.WithValue(ctx, consts.RequestUuid, c.GetString(string(consts.RequestUuid)))
	// recover password
	err := controller.RecoveryPasswordService.RecoveryPasswordRequestOTP(ctx, req.Mobile, req.NationalIdentityCode, "admin")
	if err != nil {
		logger.LogServiceV2(c, "failed to request recovery password otp", controller, err)
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).
		SetMessage("auth-otp-sent").
		SetStatusCode(http.StatusOK).
		SetLog().
		Send()
}

func (controller *PasswordRecoveryController) RecoveryPasswordViaOtp(c *gin.Context) {
	// bind request to json
	var req PasswordRecovery.RecoverRequest
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

	nationalCode := req.NationalIdentityCode
	if req.NationalCompanyId != "" {
		nationalCode = req.NationalCompanyId
	}

	ctx := context.WithValue(c.Request.Context(), consts.NationalIdentityCode, nationalCode)
	ctx = context.WithValue(ctx, consts.OwnerType, models.AdminRole)
	ctx = context.WithValue(ctx, consts.RequestUuid, c.GetString(string(consts.RequestUuid)))

	// recover password via top
	err := controller.RecoveryPasswordService.RecoveryPasswordViaOTP(ctx, req.OTP, req.NewPassword, req.NewPasswordConfirmation, req.NationalIdentityCode, "admin")
	if err != nil {
		logger.LogService(c, err, "failed to recover password via otp", "PasswordRecoveryController.RecoveryPasswordService.RecoveryPasswordViaOTP()")
		response.Api(c).SetStatusCode(http.StatusServiceUnavailable).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetLog().
		Send()
}
