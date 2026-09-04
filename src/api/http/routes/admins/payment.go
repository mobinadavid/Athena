package admins

import (
	"athena/src/api/http/middlewares"
	"athena/src/providers"

	"github.com/gin-gonic/gin"
)

func RegisterPaymentRouter(router *gin.RouterGroup) {
	serviceContainer := providers.GetContainer()
	authenticationContainer := providers.GetAuthenticationContainer()
	authorizationContainer := providers.GetAuthorizationContainer()

	auth := authenticationContainer.AuthenticationMiddleware.Middleware("admin")

	payments := router.Group("payments")
	payments.Use(auth)
	payments.GET("",
		authorizationContainer.AuthorizationMiddleware.Middleware("payment-list"),
		middlewares.PaginationMiddleware,
		serviceContainer.PaymentController.GetList,
	)
	payments.GET(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("payment-show"),
		serviceContainer.PaymentController.GetByUuid,
	)
	payments.GET(":uuid/transactions",
		authorizationContainer.AuthorizationMiddleware.Middleware("payment-show"),
		middlewares.PaginationMiddleware,
		serviceContainer.PaymentController.GetTransactions,
	)

	deposits := router.Group("deposits")
	deposits.Use(auth)
	deposits.GET("",
		authorizationContainer.AuthorizationMiddleware.Middleware("deposit-list"),
		middlewares.PaginationMiddleware,
		serviceContainer.DepositController.GetList,
	)
	deposits.GET(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("deposit-show"),
		serviceContainer.DepositController.GetByUuid,
	)

	dashboard := router.Group("dashboard")
	dashboard.Use(auth)
	dashboard.GET("",
		authorizationContainer.AuthorizationMiddleware.Middleware("dashboard-show"),
		serviceContainer.DashboardController.Summary,
	)
	dashboard.GET("deposits-over-time",
		authorizationContainer.AuthorizationMiddleware.Middleware("dashboard-show"),
		serviceContainer.DashboardController.DepositsOverTime,
	)
	dashboard.GET("deposits-by-blockchain",
		authorizationContainer.AuthorizationMiddleware.Middleware("dashboard-show"),
		serviceContainer.DashboardController.DepositsByBlockchain,
	)
}
