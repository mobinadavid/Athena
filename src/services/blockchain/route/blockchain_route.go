package route

import (
	"athena/src/api/http/middlewares"
	"athena/src/services"
	"github.com/gin-gonic/gin"
)

func BlockchainRouter(router *gin.RouterGroup) {

	serviceContainer := services.GetContainer()

	blockchain := router.Group("blockchain")

	blockchain.GET("", serviceContainer.BlockchainController.GetList)
	blockchain.GET(":uuid", middlewares.PaginationMiddleware, serviceContainer.BlockchainController.GetByUuid)
	blockchain.DELETE(":uuid", middlewares.PaginationMiddleware, serviceContainer.BlockchainController.Delete)
	blockchain.PUT(":uuid", middlewares.PaginationMiddleware, serviceContainer.BlockchainController.Update)
	blockchain.POST("", middlewares.PaginationMiddleware, serviceContainer.BlockchainController.Create)

}
