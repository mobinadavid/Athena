package routes

import (
	"athena/src/api/http/middlewares"
	"athena/src/providers"
	"github.com/gin-gonic/gin"
)

func WalletAddressRouter(router *gin.RouterGroup) {
	serviceContainer := providers.GetContainer()

	walletAddress := router.Group("wallet-address")

	walletAddress.GET("", middlewares.PaginationMiddleware, serviceContainer.WalletAddressController.GetList)
	walletAddress.POST("allocate", serviceContainer.WalletAddressController.AllocateWalletAddresses)
	walletAddress.GET("transactions/:uuid", middlewares.PaginationMiddleware, serviceContainer.WalletAddressController.GetTransactions)
	walletAddress.GET(":uuid", serviceContainer.WalletAddressController.GetByUuid)
	walletAddress.DELETE(":uuid", serviceContainer.WalletAddressController.Delete)
	walletAddress.PUT(":uuid", serviceContainer.WalletAddressController.Update)
	walletAddress.POST("", serviceContainer.WalletAddressController.Create)

}
