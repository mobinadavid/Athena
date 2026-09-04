package routes

import (
	"athena/src/providers"

	"github.com/gin-gonic/gin"
)

func WalletAddressRouter(router *gin.RouterGroup) {
	serviceContainer := providers.GetContainer()
	router.POST("/webhook/:hash", serviceContainer.WalletAddressController.Webhook)
}
