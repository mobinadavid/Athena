package routes

import (
	"athena/src/api/http/middlewares"
	"athena/src/providers"
	"github.com/gin-gonic/gin"
)

func BlockchainExplorerRouter(router *gin.RouterGroup) {
	serviceContainer := providers.GetContainer()

	blockchain := router.Group("blockchain-explorer")

	blockchain.GET("", middlewares.PaginationMiddleware, serviceContainer.BlockchainExplorerController.GetList)
	blockchain.GET(":uuid", serviceContainer.BlockchainExplorerController.GetByUuid)
	blockchain.DELETE(":uuid", serviceContainer.BlockchainExplorerController.Delete)
	blockchain.PUT(":uuid", serviceContainer.BlockchainExplorerController.Update)
	blockchain.POST("", serviceContainer.BlockchainExplorerController.Create)

}
