package admins

import (
	"athena/src/api/errs"
	"athena/src/api/http/response"
	"athena/src/config"
	"athena/src/database/scopes"
	"athena/src/models"
	"athena/src/models/consts"
	"athena/src/pkg/logger"
	"athena/src/pkg/utils"
	"athena/src/services/authentication"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"strconv"
	"time"
)

type AccessTokenController struct {
	AccessTokenService authentication.IAccessTokenService
}

func (controller *AccessTokenController) GetList(c *gin.Context) {
	// get the query builder
	builder, exists := c.Get("query_parameters_builder")
	if !exists {
		logger.LogQueryBuilderError(c)
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).SetMessage(errs.SomeThingWentWrong.Error()).SetLog().Send()
		return
	}

	// fetch information to query builder
	auth := utils.GetAuthData(c)
	builderModel := builder.(*scopes.BuilderModel)
	builderModel.Filters["owner_type"] = auth.OwnerType
	builderModel.Filters["owner_id"] = auth.OwnerId
	ctx := context.WithValue(c.Request.Context(), consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)

	// get list of access token
	data, err := controller.AccessTokenService.GetList(ctx, builderModel)
	if err != nil {
		logger.LogServiceV2(c, "failed to get list of token", controller, err)
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"access_tokens": data,
		}).SetLog().Send()
}

func (controller *AccessTokenController) GetActiveTokens(c *gin.Context) {
	// get the query builder
	builder, exists := c.Get("query_parameters_builder")
	if !exists {
		logger.LogQueryBuilderError(c)
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).SetMessage(errs.SomeThingWentWrong.Error()).SetLog().Send()
		return
	}

	// fetch information to query builder
	auth := utils.GetAuthData(c)
	builderModel := builder.(*scopes.BuilderModel)
	builderModel.Filters["owner_type"] = auth.OwnerType
	builderModel.Filters["owner_id"] = auth.OwnerId
	ctx := context.WithValue(c.Request.Context(), consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)

	// get list of access token
	data, err := controller.AccessTokenService.GetActiveTokens(ctx, builderModel)
	if err != nil {
		logger.LogServiceV2(c, "failed to get list of token", controller, err)
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// return response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"access_tokens": data,
		}).SetLog().Send()
}

func (controller *AccessTokenController) GetByUuid(c *gin.Context) {
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

	// get a token by uuid
	data, err := controller.AccessTokenService.GetByUuid(ctx, &id)
	if err != nil {
		logger.LogServiceV2(c, "failed to get access token", controller, err)
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// send response
	response.Api(c).
		SetMessage("request-successful").SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"access_token": data,
		}).SetLog().Send()
}

func (controller *AccessTokenController) RefreshAccessToken(c *gin.Context) {
	// Get refresh token from cookies
	refreshToken, cookieErr := c.Cookie(config.GetInstance().Get("COOKIE_NAME"))
	if cookieErr != nil {
		logger.LogServiceV2(c, "failed to get refresh token by cookie", controller, cookieErr)
		response.Api(c).SetStatusCode(http.StatusUnauthorized).SetMessage(errs.RefreshTokenMissing.Error()).SetLog().Send()
		return
	}

	ctx := context.WithValue(c.Request.Context(), consts.RequestUuid, c.GetString("request-uuid"))
	ctx = context.WithValue(ctx, consts.OwnerType, models.AdminRole)

	jwt, err := controller.AccessTokenService.RefreshAccessTokens(ctx, refreshToken, "admin")
	if err != nil {
		logger.LogServiceV2(c, "failed to refresh token", controller, err)
		response.Api(c).SetMessage(err.Error()).SetStatusCode(http.StatusUnauthorized).SetLog().Send()
		return
	}

	expirationSeconds, err := strconv.Atoi(config.GetInstance().Get("COOKIE_EXPIRATION"))
	if err != nil {
		logger.LogServiceV2(c, "failed to get cookie expiration by cookie", controller, err)
		response.Api(c).SetStatusCode(http.StatusUnauthorized).SetMessage(errs.SomeThingWentWrong.Error()).SetLog().Send()
		return
	}

	expiration := time.Now().Add(time.Duration(expirationSeconds) * time.Second)
	cookie := &http.Cookie{
		Name:     config.GetInstance().Get("COOKIE_NAME"),
		Value:    jwt.RefreshTokenString,
		Path:     "/",
		Domain:   config.GetInstance().Get("COOKIE_HOST"), // Domain (leave empty for default)
		Secure:   true,                                    // Set to true if using HTTPS
		HttpOnly: true,                                    // To prevent JavaScript access
		SameSite: http.SameSiteNoneMode,                   // Adjust as needed
		Expires:  expiration,                              // Set expiration date
	}

	// Set the cookie using http.SetCookie (this should automatically add the cookie to headers)
	http.SetCookie(c.Writer, cookie)

	// Return response.
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"access_tokens": jwt,
		}).SetLog().Send()
}

func (controller *AccessTokenController) RevokeTokens(c *gin.Context) {
	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)

	// revoke tokens
	err := controller.AccessTokenService.RevokeTokens(ctx, auth)
	if err != nil {
		logger.LogServiceV2(c, "failed to revoke token", controller, err)
		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetLog().
		Send()
}

func (controller *AccessTokenController) RevokeTokenByUuid(c *gin.Context) {
	// Get the UUID from the URL parameter and parse it
	uuidStr := c.Param("uuid")
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		logger.LogParseUUIDError(c, err)
		response.Api(c).SetMessage(errs.InvalidUuid.Error()).SetLog().Send()
		return
	}

	// data preparation for service
	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)

	// revoke the token by it's uuid
	err = controller.AccessTokenService.RevokeTokenByUuid(ctx, &id, auth.OwnerId, auth.OwnerType)
	if err != nil {
		logger.LogServiceV2(c, "failed to revoke token by uuid", controller, err)
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetLog().
		Send()
}

func (controller *AccessTokenController) RevokeCurrentToken(c *gin.Context) {
	auth := utils.GetAuthData(c)
	id, err := uuid.Parse(auth.TokenUUID)
	if err != nil {
		logger.LogParseUUIDError(c, err)
		response.Api(c).SetMessage(errs.InvalidUuid.Error()).SetLog().Send()
		return
	}

	// data preparation for service
	ctx := context.WithValue(c.Request.Context(), consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)

	// revoke the token by it's uuid
	err = controller.AccessTokenService.RevokeTokenByUuid(ctx, &id, auth.OwnerId, auth.OwnerType)
	if err != nil {
		logger.LogServiceV2(c, "failed to revoke current token", controller, err)
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// send response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetLog().
		Send()
}
