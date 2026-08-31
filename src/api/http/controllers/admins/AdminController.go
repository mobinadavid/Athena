package admins

import (
	"athena/src/api/errs"
	"athena/src/api/http/requests/AdminRequests"
	"athena/src/api/http/requests/Users/UserRequests"
	"athena/src/api/http/response"
	"athena/src/config"
	"athena/src/database/scopes"
	"athena/src/models/consts"
	"athena/src/pkg/logger"
	"athena/src/pkg/utils"
	"athena/src/pkg/validator"
	"athena/src/services"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"strconv"
	"time"
)

type AdminController struct {
	AdminService services.IAdminService
}

func (controller *AdminController) GetList(c *gin.Context) {
	// get the query builder
	builder, exists := c.Get("query_parameters_builder")
	if !exists {
		logger.LogQueryBuilderError(c)
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).SetMessage(errs.SomeThingWentWrong.Error()).SetLog().Send()
		return
	}

	// fetch information to query builder
	builderModel := builder.(*scopes.BuilderModel)
	builderModel.GlobalSearchCols = append(builderModel.GlobalSearchCols, `"admins"."first_name" || ' ' || "admins"."last_name"`)
	if sortBy := c.Query("sort_by"); sortBy == "" {
		builderModel.SortBy = "admins.updated_at"
	}
	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)

	// get admin list
	admins, err := controller.AdminService.GetList(ctx, builderModel)
	if err != nil {
		logger.LogServiceV2(c, "failed to get list of admins", controller, err)
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"admins": admins,
		}).SetLog().Send()
}

func (controller *AdminController) GetByUuid(c *gin.Context) {
	// Get the UUID from the URL parameter and parse it
	uuidStr := c.Param("uuid")
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		logger.LogParseUUIDError(c, err)
		response.Api(c).SetMessage(errs.InvalidUuid.Error()).SetLog().Send()
		return
	}

	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)

	// get admin by uuid
	admin, err := controller.AdminService.GetByUuid(ctx, &id)
	if err != nil {
		logger.LogServiceV2(c, "failed to get admin", controller, err)
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// return response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"admin": admin,
		}).SetLog().Send()
}

func (controller *AdminController) GetProfile(c *gin.Context) {
	//first check the body for the refresh token
	var refreshToken string
	var req UserRequests.GetProfileHeaderRequest
	if err := c.ShouldBindJSON(&req); err == nil {
		refreshToken = req.RefreshToken
	}
	if refreshToken == "" {
		cookieToken, err := c.Cookie(config.GetInstance().Get("COOKIE_NAME"))
		if err == nil {
			refreshToken = cookieToken
		}
	}
	if refreshToken == "" {
		logger.LogServiceV2(c, "failed to get refresh token", controller, errs.SomeThingWentWrong)
		response.Api(c).SetStatusCode(http.StatusBadRequest).SetMessage(errs.SomeThingWentWrong.Error()).SetLog().Send()
		return
	}

	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)
	userId := auth.OwnerId
	// get admin by id
	admin, err := controller.AdminService.GetProfile(ctx, userId)
	if err != nil {
		logger.LogServiceV2(c, "failed to get profile", controller, err)
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	expirationSeconds, err := strconv.Atoi(config.GetInstance().Get("COOKIE_EXPIRATION"))
	if err != nil {
		logger.LogErrorWithFieldsV2(c, "failed to cast COOKIE_EXPIRATION to string", controller, err)
		response.Api(c).SetMessage(errs.ErrAuthenticationFailed.Error()).SetLog().Send()
		return
	}

	expiration := time.Now().Add(time.Duration(expirationSeconds) * time.Second)
	cookie := &http.Cookie{
		Name:     config.GetInstance().Get("COOKIE_NAME"),
		Value:    refreshToken,
		Path:     "/",
		Domain:   config.GetInstance().Get("COOKIE_HOST"), // Domain (leave empty for default)
		Secure:   true,                                    // Set to true if using HTTPS
		HttpOnly: true,                                    // To prevent JavaScript access
		SameSite: http.SameSiteNoneMode,                   // Adjust as needed
		Expires:  expiration,                              // Set expiration date
	}

	// Set the cookie using http.SetCookie (this should automatically add the cookie to headers)
	http.SetCookie(c.Writer, cookie)

	// return response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"admin": admin,
		}).SetLog().Send()
}

func (controller *AdminController) Create(c *gin.Context) {
	// bind request to json
	var req AdminRequests.CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogJSONBindError(c, err)
		response.Api(c).SetLog().Send()
		return
	}

	// Validate the payload.
	if err := validator.Validate(&req, c.GetString("locale")); err != nil {
		logger.LogValidationError(c, err)
		response.Api(c).SetErrors(err).SetLog().Send()
		return
	}
	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)
	creatorAdminID := auth.OwnerId
	// create admin
	admin, err := controller.AdminService.Create(ctx, &req, creatorAdminID)
	if err != nil {
		logger.LogServiceV2(c, "failed to create admin", controller, err)
		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// return response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusCreated).
		SetData(map[string]interface{}{
			"admin": admin,
		}).SetLog().Send()
}

func (controller *AdminController) Update(c *gin.Context) {
	// Get the UUID from the URL parameter and parse it
	uuidStr := c.Param("uuid")
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		logger.LogParseUUIDError(c, err)
		response.Api(c).SetMessage(errs.InvalidUuid.Error()).SetLog().Send()
		return
	}

	// bind request to json
	var req AdminRequests.UpdateAdminRequest
	if err = c.ShouldBindJSON(&req); err != nil {
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

	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)

	// update admin
	admin, err := controller.AdminService.UpdateByUuid(ctx, &id, &req)
	if err != nil {
		logger.LogServiceV2(c, "failed to update admin", controller, err)
		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// return response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"admin": admin,
		}).SetLog().Send()
}

func (controller *AdminController) Delete(c *gin.Context) {
	// Get the UUID from the URL parameter and parse it
	uuidStr := c.Param("uuid")
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		logger.LogParseUUIDError(c, err)
		response.Api(c).SetMessage(errs.InvalidUuid.Error()).SetLog().Send()
		return
	}

	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)

	// delete admin
	err = controller.AdminService.DeleteByUuid(ctx, &id)
	if err != nil {
		logger.LogServiceV2(c, "failed to delete admin", controller, err)
		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// send response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetLog().
		Send()
}

func (controller *AdminController) ChangePasswordVerifyOtp(c *gin.Context) {
	// bind request to json
	var req AdminRequests.ChangePasswordVerifyOtpRequest
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
	// get cookie name
	cookieName := config.GetInstance().Get("CHANGE_PASSWORD_COOKIE_KEY_NAME")
	if cookieName == "" {
		cookieName = "change_password_key"
	}

	// check the key is present or not
	if req.Key == "" {
		// Read the cookie set in Step 1
		key, err := c.Cookie(cookieName)
		if err != nil {
			// in this case cookie is expired
			logger.LogCookieDoesNotExist(c, err, cookieName)
			response.Api(c).SetMessage(errs.ErrChangePasswordTimeOut.Error()).SetErrorCode(errs.LoginTimeOutErrorCode).SetLog().Send()
			return
		}

		req.Key = key
	}

	auth := utils.GetAuthData(c)
	adminID := auth.OwnerId
	ctx := context.WithValue(context.Background(), "admin_id", adminID)
	ctx = context.WithValue(ctx, "req", req)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)
	// change password
	err := controller.AdminService.ChangePassword(ctx)
	if err != nil {
		logger.LogServiceV2(c, "failed to change password", controller, err)
		response.Api(c).SetMessage(err.Error()).Send()
		return
	}

	// send response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetLog().
		Send()
}

func (controller *AdminController) ChangePasswordSendOtp(c *gin.Context) {
	var req AdminRequests.ChangePasswordSendOtpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogJSONBindError(c, err)
		response.Api(c).SetLog().Send()
		return
	}
	// validate request
	if err := validator.Validate(&req, c.GetString("locale")); err != nil {
		logger.LogValidationError(c, err)
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).SetErrors(err).SetLog().Send()
		return
	}

	auth := utils.GetAuthData(c)
	adminID := auth.OwnerId
	ctx := context.WithValue(context.Background(), "admin_id", adminID)
	ctx = context.WithValue(ctx, "req", req)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)
	key, err := controller.AdminService.ChangePasswordSendOTP(ctx)
	fmt.Println("DEBUG OTP key:", key)
	if err != nil {
		logger.LogServiceV2(c, "failed to send otp for change admin password", controller, err)

		if utils.CheckError(err, errs.ErrAuthOTPExists) {
			response.Api(c).SetMessage(err.Error()).SetErrorCode(errs.OTPAlreadyExistErrorCode).SetLog().Send()
			return
		}

		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}
	// get expire time
	expiration, err := strconv.Atoi(config.GetInstance().Get("CHANGE_PASSWORD_LIFETIME"))
	if err != nil {
		logger.LogAToIError(c, err)
		expiration = 300
	}

	// get cookie name
	cookieName := config.GetInstance().Get("CHANGE_PASSWORD_COOKIE_KEY_NAME")
	if cookieName == "" {
		cookieName = "change_password_key"
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

func (controller *AdminController) ChangePasswordResendOTP(c *gin.Context) {
	// Bind check payload.
	var req AdminRequests.ChangePasswordResendOtpRequest
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
	cookieName := config.GetInstance().Get("CHANGE_PASSWORD_COOKIE_KEY_NAME")
	if cookieName == "" {
		cookieName = "change_password_key"
	}

	// check the key is present or not
	if req.Key == "" {
		// Read the cookie set in Step 1
		key, err := c.Cookie(cookieName)
		if err != nil {
			// in this case cookie is expired
			logger.LogCookieDoesNotExist(c, err, cookieName)
			response.Api(c).SetMessage(errs.ErrChangePasswordTimeOut.Error()).SetErrorCode(errs.RegisterTimeOutErrorCode).SetLog().Send()
			return
		}

		req.Key = key
	}

	// prepare data for service
	auth := utils.GetAuthData(c)
	ctx := context.WithValue(context.Background(), "req", &req)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)
	// register admin
	err := controller.AdminService.ChangePasswordResendOTP(ctx)
	if err != nil {
		logger.LogServiceV2(c, "failed to resend change password otp", controller, err)

		if utils.CheckError(err, errs.ErrAuthOTPExists) {
			response.Api(c).SetMessage(err.Error()).SetErrorCode(errs.OTPAlreadyExistErrorCode).SetLog().Send()
			return
		}
		if utils.CheckError(err, errs.ErrChangePasswordTimeOut) {
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
			"key": req.Key,
		}).SetLog().Send()
}
