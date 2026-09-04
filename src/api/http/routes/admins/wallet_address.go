package admins

import (
	"athena/src/api/http/middlewares"
	"athena/src/providers"

	"github.com/gin-gonic/gin"
)

func RegisterWalletAddressRouter(router *gin.RouterGroup) {
	serviceContainer := providers.GetContainer()
	authenticationContainer := providers.GetAuthenticationContainer()
	authorizationContainer := providers.GetAuthorizationContainer()

	router.POST("/webhook/:hash", serviceContainer.WalletAddressController.Webhook)

	walletAddress := router.Group("wallet-address")
	walletAddress.Use(authenticationContainer.AuthenticationMiddleware.Middleware("admin"))
	walletAddress.GET("",
		authorizationContainer.AuthorizationMiddleware.Middleware("wallet-address-list"),
		middlewares.PaginationMiddleware,
		serviceContainer.WalletAddressController.GetList,
	)
	walletAddress.POST("allocate",
		authorizationContainer.AuthorizationMiddleware.Middleware("wallet-address-allocate"),
		serviceContainer.WalletAddressController.AllocateWalletAddresses,
	)
	walletAddress.GET("transactions/:uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("wallet-address-transactions"),
		middlewares.PaginationMiddleware,
		serviceContainer.WalletAddressController.GetTransactions,
	)
	walletAddress.GET(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("wallet-address-show"),
		serviceContainer.WalletAddressController.GetByUuid,
	)
	walletAddress.POST("",
		authorizationContainer.AuthorizationMiddleware.Middleware("wallet-address-create"),
		serviceContainer.WalletAddressController.Create,
	)
	walletAddress.PUT(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("wallet-address-update"),
		serviceContainer.WalletAddressController.Update,
	)
	walletAddress.DELETE(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("wallet-address-delete"),
		serviceContainer.WalletAddressController.Delete,
	)
}
