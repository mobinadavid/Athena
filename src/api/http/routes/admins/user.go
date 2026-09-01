package admins

import (
	"athena/src/api/http/middlewares"
	"athena/src/models"
	"athena/src/providers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(router *gin.RouterGroup) {
	authenticationContainer := providers.GetAuthenticationContainer()
	authorizationContainer := providers.GetAuthorizationContainer()
	userContainer := providers.GetUserContainer()

	users := router.Group("users")
	users.Use(authenticationContainer.AuthenticationMiddleware.Middleware("admin"))

	// GET /users
	users.GET("",
		authorizationContainer.AuthorizationMiddleware.Middleware("user-list"),
		middlewares.QueryParametersBuilderMiddleware(models.UserModel{}),
		userContainer.AdminUserController.GetList,
	)

	users.POST("",
		authorizationContainer.AuthorizationMiddleware.Middleware("user-create"),
		userContainer.AdminUserController.Create,
	)

	users.PUT(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("user-update"),
		userContainer.AdminUserController.Update,
	)

	// GET /users/:uuid
	users.GET(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("user-show"),
		userContainer.AdminUserController.GetByUuid,
	)

	// Get /users/profile/:uuid
	users.GET("profile/:uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("user-show"),
		userContainer.AdminUserController.GetUserProfile,
	)

}
