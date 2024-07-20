package route

import (
	"athena/src/api/http/middlewares"
	"athena/src/services"
	"github.com/gin-gonic/gin"
)

func BlockchainExplorerRouter(router *gin.RouterGroup) {

	serviceContainer := services.GetContainer()

	blockchain := router.Group("blockchain_explorer")

	blockchain.GET("", serviceContainer.BlockchainExplorerController.GetList)
	blockchain.GET(":uuid", middlewares.PaginationMiddleware, serviceContainer.BlockchainExplorerController.GetByUuid)
	blockchain.DELETE(":uuid", middlewares.PaginationMiddleware, serviceContainer.BlockchainExplorerController.Delete)
	blockchain.PUT(":uuid", middlewares.PaginationMiddleware, serviceContainer.BlockchainExplorerController.Update)
	blockchain.POST("", middlewares.PaginationMiddleware, serviceContainer.BlockchainExplorerController.Create)

}
