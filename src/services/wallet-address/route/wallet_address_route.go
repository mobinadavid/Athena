package route

import (
	"athena/src/api/http/middlewares"
	"athena/src/services"
	"github.com/gin-gonic/gin"
)

func WalletAddressRouter(router *gin.RouterGroup) {
	serviceContainer := services.GetContainer()

	walletAddress := router.Group("wallet_address")

	walletAddress.GET("", serviceContainer.WalletAddressController.GetList)
	walletAddress.GET("walletAddresses", serviceContainer.WalletAddressController.GetWalletAddress)
	walletAddress.GET("transactions/:uuid", serviceContainer.WalletAddressController.GetTransactions)
	walletAddress.GET(":uuid", middlewares.PaginationMiddleware, serviceContainer.WalletAddressController.GetByUuid)
	walletAddress.DELETE(":uuid", middlewares.PaginationMiddleware, serviceContainer.WalletAddressController.Delete)
	walletAddress.PUT(":uuid", middlewares.PaginationMiddleware, serviceContainer.WalletAddressController.Update)
	walletAddress.POST("", middlewares.PaginationMiddleware, serviceContainer.WalletAddressController.Create)

}
