package admins

import (
	"athena/src/api/http/middlewares"
	"athena/src/providers"

	"github.com/gin-gonic/gin"
)

func RegisterBlockchainExplorerRouter(router *gin.RouterGroup) {
	serviceContainer := providers.GetContainer()
	authenticationContainer := providers.GetAuthenticationContainer()
	authorizationContainer := providers.GetAuthorizationContainer()

	explorer := router.Group("blockchain-explorer")
	explorer.Use(authenticationContainer.AuthenticationMiddleware.Middleware("admin"))
	explorer.GET("",
		authorizationContainer.AuthorizationMiddleware.Middleware("blockchain-explorer-list"),
		middlewares.PaginationMiddleware,
		serviceContainer.BlockchainExplorerController.GetList,
	)
	explorer.GET(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("blockchain-explorer-show"),
		serviceContainer.BlockchainExplorerController.GetByUuid,
	)
	explorer.POST("",
		authorizationContainer.AuthorizationMiddleware.Middleware("blockchain-explorer-create"),
		serviceContainer.BlockchainExplorerController.Create,
	)
	explorer.PUT(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("blockchain-explorer-update"),
		serviceContainer.BlockchainExplorerController.Update,
	)
	explorer.DELETE(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("blockchain-explorer-delete"),
		serviceContainer.BlockchainExplorerController.Delete,
	)
}
