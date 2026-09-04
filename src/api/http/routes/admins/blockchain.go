package admins

import (
	"athena/src/api/http/middlewares"
	"athena/src/providers"

	"github.com/gin-gonic/gin"
)

func RegisterBlockchainRouter(router *gin.RouterGroup) {
	serviceContainer := providers.GetContainer()
	authenticationContainer := providers.GetAuthenticationContainer()
	authorizationContainer := providers.GetAuthorizationContainer()

	blockchain := router.Group("blockchain")
	blockchain.Use(authenticationContainer.AuthenticationMiddleware.Middleware("admin"))
	blockchain.GET("",
		authorizationContainer.AuthorizationMiddleware.Middleware("blockchain-list"),
		middlewares.PaginationMiddleware,
		serviceContainer.BlockchainController.GetList,
	)
	blockchain.GET(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("blockchain-show"),
		serviceContainer.BlockchainController.GetByUuid,
	)
	blockchain.POST("",
		authorizationContainer.AuthorizationMiddleware.Middleware("blockchain-create"),
		serviceContainer.BlockchainController.Create,
	)
	blockchain.PUT(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("blockchain-update"),
		serviceContainer.BlockchainController.Update,
	)
	blockchain.DELETE(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("blockchain-delete"),
		serviceContainer.BlockchainController.Delete,
	)
}
