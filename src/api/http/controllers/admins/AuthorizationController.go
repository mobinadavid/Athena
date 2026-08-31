package admins

import (
	"athena/src/api/errs"
	"athena/src/api/http/requests/AuthorizationRequests"
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

type AuthorizationController struct {
	AuthorizationService services.IAuthorizationService
}

func (controller *AuthorizationController) GetPermissions(c *gin.Context) {
	// get the query builder
	builder, exists := c.Get("query_parameters_builder")
	if !exists {
		logger.LogQueryBuilderError(c)
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).SetMessage(errs.SomeThingWentWrong.Error()).SetLog().Send()
		return
	}

	// fetch information to query builder
	builderModel := builder.(*scopes.BuilderModel)
	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)

	// get permissions
	data, err := controller.AuthorizationService.GetPermissionsList(ctx, builderModel)
	if err != nil {
		logger.LogServiceV2(c, "failed to get list of permission", controller, err)
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// send response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"permissions": data,
		}).SetLog().Send()
}

func (controller *AuthorizationController) GetPermissionGroups(c *gin.Context) {
	// get the query builder
	builder, exists := c.Get("query_parameters_builder")
	if !exists {
		logger.LogQueryBuilderError(c)
		response.Api(c).SetStatusCode(http.StatusUnprocessableEntity).SetMessage(errs.SomeThingWentWrong.Error()).SetLog().Send()
		return
	}

	// fetch information to query builder
	builderModel := builder.(*scopes.BuilderModel)
	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)

	// get permissions
	data, err := controller.AuthorizationService.GetPermissionGroupsList(ctx, builderModel)
	if err != nil {
		logger.LogServiceV2(c, "failed to get list of permission groups", controller, err)
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// send response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"permission_groups": data,
		}).SetLog().Send()
}

func (controller *AuthorizationController) GetPermissionGroupByUuid(c *gin.Context) {
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

	// get PermissionGroup
	permissionGroup, err := controller.AuthorizationService.GetPermissionGroupByUuid(ctx, &id, "Permissions")
	if err != nil {
		logger.LogServiceV2(c, "failed to get PermissionGroup", controller, err)
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// send response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"permission_group": permissionGroup,
		}).SetLog().Send()
}

func (controller *AuthorizationController) CreatePermissionGroup(c *gin.Context) {
	// Bind check payload.
	var req AuthorizationRequests.CreatePermissionGroupRequest
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

	// create permission group
	group, err := controller.AuthorizationService.CreatePermissionGroup(ctx, &req)
	if err != nil {
		logger.LogServiceV2(c, "failed to create permission group", controller, err)
		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusCreated).
		SetData(map[string]interface{}{
			"permission_group": group,
		}).SetLog().Send()
}

func (controller *AuthorizationController) UpdatePermissionGroup(c *gin.Context) {
	// Bind check payload.
	var req AuthorizationRequests.UpdatePermissionGroupRequest
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

	// create permission group
	group, err := controller.AuthorizationService.UpdatePermissionGroupByUuid(ctx, &id, &req)
	if err != nil {
		logger.LogServiceV2(c, "failed to create permission group", controller, err)
		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusCreated).
		SetData(map[string]interface{}{
			"permission_group": group,
		}).SetLog().Send()
}

func (controller *AuthorizationController) DeletePermissionGroup(c *gin.Context) {
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

	// create permission group
	err = controller.AuthorizationService.DeletePermissionGroupByUuid(ctx, &id)
	if err != nil {
		logger.LogServiceV2(c, "failed to delete permission group", controller, err)
		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusCreated).
		SetLog().Send()
}

func (controller *AuthorizationController) GetRolesList(c *gin.Context) {
	// get the query builder
	builder, exists := c.Get("query_parameters_builder")
	if !exists {
		logger.LogQueryBuilderError(c)
		response.Api(c).SetMessage(errs.SomeThingWentWrong.Error()).SetStatusCode(http.StatusUnprocessableEntity).SetLog().Send()
		return
	}

	// fetch information to query builder
	builderModel := builder.(*scopes.BuilderModel)
	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)

	// get role list
	roles, err := controller.AuthorizationService.GetRolesList(ctx, builderModel)
	if err != nil {
		logger.LogServiceV2(c, "failed to get list of roles", controller, err)
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// send response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"roles": roles,
		}).SetLog().Send()
}

func (controller *AuthorizationController) GetRoleByUuid(c *gin.Context) {
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

	// get role
	role, err := controller.AuthorizationService.GetRoleByUuid(ctx, &id)
	if err != nil {
		logger.LogServiceV2(c, "failed to get role", controller, err)
		response.Api(c).SetStatusCode(http.StatusNotFound).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// send response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"role": role,
		}).SetLog().Send()
}

func (controller *AuthorizationController) CreateRole(c *gin.Context) {
	// Bind check payload.
	var req AuthorizationRequests.CreateRoleRequest
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

	// create role
	role, err := controller.AuthorizationService.CreateRole(ctx, &req)
	if err != nil {
		logger.LogServiceV2(c, "failed to create role", controller, err)
		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusCreated).
		SetData(map[string]interface{}{
			"role": role,
		}).SetLog().Send()
}

func (controller *AuthorizationController) UpdateRole(c *gin.Context) {
	// Get the UUID from the URL parameter
	uuidStr := c.Param("uuid")
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		logger.LogParseUUIDError(c, err)
		response.Api(c).SetMessage(errs.InvalidUuid.Error()).SetLog().Send()
		return
	}

	// Bind check payload.
	var req AuthorizationRequests.UpdateRoleRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.LogJSONBindError(c, err)
		response.Api(c).SetLog().Send()
		return
	}

	// Validate the payload.
	if err := validator.Validate(&req, c.GetString("locale")); err != nil {
		logger.LogValidationError(c, err)
		response.Api(c).SetErrors(err).Send()
		return
	}

	auth := utils.GetAuthData(c)
	ctx := context.WithValue(c.Request.Context(), consts.OwnerId, auth.OwnerId)
	ctx = context.WithValue(ctx, consts.OwnerType, auth.OwnerType)
	ctx = context.WithValue(ctx, consts.RequestUuid, auth.RequestUuid)

	// update role
	role, err := controller.AuthorizationService.UpdateRoleByUuid(ctx, &id, &req)
	if err != nil {
		logger.LogServiceV2(c, "failed to update role", controller, err)
		response.Api(c).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// Return response.
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetData(map[string]interface{}{
			"role": role,
		}).SetLog().Send()
}

func (controller *AuthorizationController) DeleteRole(c *gin.Context) {
	// Get the UUID from the URL parameter
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

	// delete role
	err = controller.AuthorizationService.DeleteRoleByUuid(ctx, &id)
	if err != nil {
		logger.LogServiceV2(c, "failed to delete role", controller, err)
		response.Api(c).SetStatusCode(http.StatusBadRequest).SetMessage(err.Error()).SetLog().Send()
		return
	}

	// return response
	response.Api(c).SetMessage("request-successful").
		SetStatusCode(http.StatusOK).
		SetLog().
		Send()
}
