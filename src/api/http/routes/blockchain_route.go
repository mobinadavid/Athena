package routes

import (
	"athena/src/api/http/middlewares"
	"athena/src/providers"

	"github.com/gin-gonic/gin"
)

func BlockchainRouter(router *gin.RouterGroup) {

	serviceContainer := providers.GetContainer()
	authenticationContainer := providers.GetAuthenticationContainer()

	blockchain := router.Group("blockchain")
	blockchain.Use(authenticationContainer.AuthenticationMiddleware.Middleware("user"))
	blockchain.GET("", middlewares.PaginationMiddleware, serviceContainer.BlockchainController.GetList)
	blockchain.GET(":uuid", serviceContainer.BlockchainController.GetByUuid)
	blockchain.DELETE(":uuid", serviceContainer.BlockchainController.Delete)
	blockchain.PUT(":uuid", serviceContainer.BlockchainController.Update)
	blockchain.POST("", serviceContainer.BlockchainController.Create)

}
