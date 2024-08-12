package routes

import (
	"athena/src/providers"
	"github.com/gin-gonic/gin"
)

func IpgRouter(router *gin.RouterGroup) {
	serviceContainer := providers.GetContainer()

	ipg := router.Group("ipg")
	ipg.POST("request-payment", serviceContainer.IpgController.RequestPayment)
	ipg.POST("ipg-callback/:uuid", serviceContainer.IpgController.VerifyPayment)
	ipg.GET("igp/:uuid", serviceContainer.IpgController.GetIGPByUuid)

}
