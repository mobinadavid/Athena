package admins

import (
	"athena/src/api/errs"
	"athena/src/api/http/requests/AdminRequests"
	"athena/src/api/http/response"
	"athena/src/database/scopes"
	"athena/src/models/consts"
	"athena/src/pkg/logger"
	"athena/src/pkg/utils"
	"athena/src/pkg/validator"
	"athena/src/services"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	UserService services.IUserService
}

func (controller *UserController) GetList(c *gin.Context) {
	builder, exists := c.Get("query_parameters_builder")
	if !exists {
		logger.LogQueryBuilderError(c)
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).SetMessage(errs.SomeThingWentWrong.Error()).SetLog().Send()
		return
	}

	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)

	// prepare data for service
	builderModel := builder.(*scopes.BuilderModel)
	//add global search cols
	builderModel.GlobalSearchCols = append(builderModel.GlobalSearchCols, `"users"."first_name" || ' ' || "users"."last_name"`)
	// get user list
	data, err := controller.UserService.GetList(ctx, builderModel)
	if err != nil {
		logger.LogServiceV2(c, "failed to get user list", controller, err)
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// send response
	response.Api(c).
		SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"users": data,
		}).SetLog().Send()
}

func (controller *UserController) GetByUuid(c *gin.Context) {
	// Get the UUID from the URL parameter
	uuidStr := c.Param("uuid")
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		logger.LogParseUUIDError(c, err)
		response.Api(c).SetMessage(errs.InvalidUuid.Error()).SetLog().Send()
		return
	}

	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)

	// Use the UserRepository to find the user by UUID
	user, err := controller.UserService.GetByUuid(ctx, &id)
	if err != nil {
		logger.LogServiceV2(c, "failed to get user", controller, err)
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	response.Api(c).
		SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"user": user,
		}).SetLog().Send()
}
func (controller *UserController) Create(c *gin.Context) {
	var req AdminRequests.CreateUpdateUserByAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogJSONBindError(c, err)
		response.Api(c).SetStatusCode(http.StatusBadRequest).SetMessage(errs.SomeThingWentWrong.Error()).SetLog().Send()
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

	user, err := controller.UserService.CreateUserByAdmin(ctx, &req)
	if err != nil {
		logger.LogServiceV2(c, "failed to create user by admin", controller, err)
		response.Api(c).SetStatusCode(http.StatusBadRequest).SetMessage(err.Error()).SetLog().Send()
		return
	}

	response.Api(c).
		SetMessage("request-successful").
		SetStatusCode(http.StatusCreated).
		SetData(map[string]interface{}{
			"user": user,
		}).SetLog().Send()
}

func (controller *UserController) Update(c *gin.Context) {
	uuidStr := c.Param("uuid")
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		logger.LogParseUUIDError(c, err)
		response.Api(c).SetStatusCode(http.StatusBadRequest).SetMessage(errs.InvalidUuid.Error()).SetLog().Send()
		return
	}

	var req AdminRequests.CreateUpdateUserByAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogJSONBindError(c, err)
		response.Api(c).SetStatusCode(http.StatusBadRequest).SetMessage(errs.SomeThingWentWrong.Error()).SetLog().Send()
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

	user, err := controller.UserService.UpdateUserByAdmin(ctx, id, &req)
	if err != nil {
		logger.LogServiceV2(c, "failed to update user by admin", controller, err)
		response.Api(c).SetStatusCode(http.StatusBadRequest).SetMessage(err.Error()).SetLog().Send()
		return
	}

	response.Api(c).
		SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"user": user,
		}).SetLog().Send()
}

func (controller *UserController) GetUserProfile(c *gin.Context) {
	// Get the UUID from the URL parameter
	uuidStr := c.Param("uuid")
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		logger.LogParseUUIDError(c, err)
		response.Api(c).SetMessage(errs.InvalidUuid.Error()).SetLog().Send()
		return
	}

	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)
	userObj, err := controller.UserService.GetByUuid(ctx, &id)
	if err != nil {
		logger.LogServiceV2(c, "failed to get user by uuid", controller, err)
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	user, err := controller.UserService.GetProfile(ctx, userObj.ID,
		"CompanyDocuments", "ShareholderDocuments", "BankAccounts", "BankAccounts.Bank", "BoardMemberDocuments")
	if err != nil {
		logger.LogServiceV2(c, "failed to get user by id", controller, err)
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// send response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"user": user,
		}).SetLog().Send()
}
