package providers

import (
	"athena/src/api/http/controllers"
	"athena/src/services"
	"github.com/google/wire"
)

var IpgContainer = wire.NewSet(
	ProvideIpgController,
	ProvideIpgService,
	wire.Bind(new(services.IIpgService), new(*services.IpgService)),
)

func ProvideIpgController(service services.IIpgService) *controllers.IpgController {
	return &controllers.IpgController{
		IIpgService: service,
	}
}

func ProvideIpgService(igpService services.IIGPService) *services.IpgService {
	return &services.IpgService{
		IIGPService: igpService,
	}
}
