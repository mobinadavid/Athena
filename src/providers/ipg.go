package providers

import (
	"athena/src/api/http/controllers"
	"athena/src/services"
	"github.com/google/wire"
)

var IpgContainer = wire.NewSet(
	ProvideIpgController,
	ProvideIpgService,
)

func ProvideIpgController(service services.IIpgService) *controllers.IpgController {
	return &controllers.IpgController{
		IIpgService: service,
	}
}

func ProvideIpgService() services.IIpgService {
	return &services.IpgService{}
}
