package routes

import (
	"athena/src/api/http/middlewares"
	"athena/src/providers"

	"github.com/gin-gonic/gin"
)

func RegisterPaymentTrackingRouter(router *gin.RouterGroup) {
	container := providers.GetContainer()
	auth := providers.GetAuthenticationContainer()

	payments := router.Group("payments")
	payments.Use(auth.AuthenticationMiddleware.Middleware("user"))
	{
		payments.POST("", container.PaymentController.Create)
		payments.GET("", middlewares.PaginationMiddleware, container.PaymentController.GetList)
		payments.GET(":uuid", container.PaymentController.GetByUuid)
		payments.GET(":uuid/transactions", middlewares.PaginationMiddleware, container.PaymentController.GetTransactions)
	}

	deposits := router.Group("deposits")
	deposits.Use(auth.AuthenticationMiddleware.Middleware("user"))
	{
		deposits.GET("", middlewares.PaginationMiddleware, container.DepositController.GetList)
		deposits.GET(":uuid", container.DepositController.GetByUuid)
	}

	notifications := router.Group("notifications")
	notifications.Use(auth.AuthenticationMiddleware.Middleware("user"))
	{
		notifications.GET("", middlewares.PaginationMiddleware, container.NotificationController.GetList)
		notifications.POST("read-all", container.NotificationController.MarkAllRead)
		notifications.POST(":uuid/read", container.NotificationController.MarkRead)
	}

	wallets := router.Group("wallets")
	wallets.Use(auth.AuthenticationMiddleware.Middleware("user"))
	{
		wallets.GET("", container.WalletAddressController.GetMyAllocated)
		wallets.GET(":uuid/transactions", middlewares.PaginationMiddleware, container.WalletAddressController.GetOwnedTransactions)
	}

	dashboard := router.Group("dashboard")
	dashboard.Use(auth.AuthenticationMiddleware.Middleware("user"))
	{
		dashboard.GET("", container.DashboardController.Summary)
		dashboard.GET("deposits-over-time", container.DashboardController.DepositsOverTime)
		dashboard.GET("deposits-by-blockchain", container.DashboardController.DepositsByBlockchain)
	}
}
