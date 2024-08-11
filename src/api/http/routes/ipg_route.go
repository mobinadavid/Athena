package routes

import (
	"athena/src/providers"
	"github.com/gin-gonic/gin"
)

func IpgRouter(router *gin.RouterGroup) {
	serviceContainer := providers.GetContainer()

	ipg := router.Group("ipg")
	ipg.POST("", serviceContainer.IpgController.RequestPayment)

}
