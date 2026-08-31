//go:build wireinject
// +build wireinject

package providers

import (
	"athena/src/api/http/controllers"
	"athena/src/api/http/controllers/users/authentication"
	"athena/src/api/http/middlewares"
	"athena/src/database"
	"athena/src/services"

	"github.com/google/wire"
)

type (
	Container struct {
		BlockchainController         *controllers.BlockchainController
		WalletAddressController      *controllers.WalletAddressController
		BlockchainExplorerController *controllers.BlockchainExplorerController
		IpgController                *controllers.IpgController
		DepositService               *services.DepositService
		IgpService                   *services.IGPService
	}

	AuthenticationContainer struct {
		UserRegisterController         *authentication.RegisterController
		UserLoginController            *authentication.LoginController
		UserAccessTokenController      *authentication.AccessTokenController
		UserRecoveryPasswordController *authentication.RecoverPasswordController
		AuthenticationMiddleware       *middlewares.AuthenticationMiddleware
	}
)

func GetContainer() *Container {
	wire.Build(
		database.GetInstance,
		BlockchainContainer,
		ExplorerContainer,
		WalletAddressContainer,
		DepositContainer,
		IpgContainer,
		IgpContainer,
		wire.Struct(new(Container), "*"),
	)

	return nil
}

func GetAuthenticationContainer() *AuthenticationContainer {
	wire.Build(
		// Repositories
		database.GetInstance,
		ProvideUserRepository,
		ProvideAccessTokenRepository,
		// Services
		ProvideRegisterService,
		ProvideLoginService,
		ProvideJwtService,
		ProvideUserService,
		ProvideAccessTokenService,
		ProvideRecoveryPasswordService,
		ProvideOTPService,
		// Controllers
		ProvideUserRegisterController,
		ProvideUserLoginController,
		ProvideUserRecoverPasswordController,
		ProvideUserAccessTokenController,
		// Middlewares
		ProvideAuthenticationMiddleware,
		wire.Struct(new(AuthenticationContainer), "*"),
	)
	return nil
}
