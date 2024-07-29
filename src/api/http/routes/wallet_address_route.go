package routes

import (
	"athena/src/api/http/middlewares"
	"athena/src/providers"
	"github.com/gin-gonic/gin"
)

func WalletAddressRouter(router *gin.RouterGroup) {
	serviceContainer := providers.GetContainer()

	walletAddress := router.Group("wallet-address")

	walletAddress.GET("", serviceContainer.WalletAddressController.GetList)
	walletAddress.POST("allocate", serviceContainer.WalletAddressController.AllocateWalletAddresses)
	walletAddress.GET("transactions/:uuid", serviceContainer.WalletAddressController.GetTransactions)
	walletAddress.GET(":uuid", middlewares.PaginationMiddleware, serviceContainer.WalletAddressController.GetByUuid)
	walletAddress.DELETE(":uuid", middlewares.PaginationMiddleware, serviceContainer.WalletAddressController.Delete)
	walletAddress.PUT(":uuid", middlewares.PaginationMiddleware, serviceContainer.WalletAddressController.Update)
	walletAddress.POST("", middlewares.PaginationMiddleware, serviceContainer.WalletAddressController.Create)

}
