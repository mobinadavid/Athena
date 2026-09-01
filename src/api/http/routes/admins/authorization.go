package admins

import (
	"athena/src/api/http/middlewares"
	"athena/src/models"
	"athena/src/providers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthorizationRouter(router *gin.RouterGroup) {
	authenticationContainer := providers.GetAuthenticationContainer()
	authorizationContainer := providers.GetAuthorizationContainer()

	authorization := router.Group("authorization")
	authorization.Use(authenticationContainer.AuthenticationMiddleware.Middleware("admin"))

	// GET /authorization/permissions
	permissions := authorization.Group("permissions")
	permissions.GET("",
		authorizationContainer.AuthorizationMiddleware.Middleware("permission-list"),
		middlewares.QueryParametersBuilderMiddleware(models.PermissionModel{}),
		authorizationContainer.AdminAuthorizationController.GetPermissions,
	)

	// GET /authorization/permissionGroup
	permissionGroups := authorization.Group("permission-groups")
	permissionGroups.GET("",
		authorizationContainer.AuthorizationMiddleware.Middleware("permission-group-list"),
		middlewares.QueryParametersBuilderMiddleware(models.PermissionGroupModel{}),
		authorizationContainer.AdminAuthorizationController.GetPermissionGroups,
	)
	permissionGroups.GET(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("permission-group-show"),
		authorizationContainer.AdminAuthorizationController.GetPermissionGroupByUuid,
	)
	permissionGroups.POST("",
		authorizationContainer.AuthorizationMiddleware.Middleware("permission-group-create"),
		authorizationContainer.AdminAuthorizationController.CreatePermissionGroup,
	)
	permissionGroups.PUT(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("permission-group-update"),
		authorizationContainer.AdminAuthorizationController.UpdatePermissionGroup,
	)
	permissionGroups.DELETE(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("permission-group-delete"),
		authorizationContainer.AdminAuthorizationController.DeletePermissionGroup,
	)

	// GET /authorization/roles
	roles := authorization.Group("roles")
	roles.GET("",
		authorizationContainer.AuthorizationMiddleware.Middleware("role-list"),
		middlewares.QueryParametersBuilderMiddleware(models.RoleModel{}),
		authorizationContainer.AdminAuthorizationController.GetRolesList,
	)

	roles.GET(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("role-show"),
		authorizationContainer.AdminAuthorizationController.GetRoleByUuid,
	)

	roles.POST("",
		authorizationContainer.AuthorizationMiddleware.Middleware("role-create"),
		authorizationContainer.AdminAuthorizationController.CreateRole,
	)

	roles.PUT(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("role-update"),
		authorizationContainer.AdminAuthorizationController.UpdateRole,
	)

	roles.DELETE(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("role-delete"),
		authorizationContainer.AdminAuthorizationController.DeleteRole,
	)
}
