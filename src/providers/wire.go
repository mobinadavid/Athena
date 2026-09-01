//go:build wireinject
// +build wireinject

package providers

import (
	"athena/src/api/http/controllers"
	adminController "athena/src/api/http/controllers/admins"
	adminAuthenticationController "athena/src/api/http/controllers/admins/authentication"
	userAuthenticationController "athena/src/api/http/controllers/users/authentication"
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
		UserRegisterController          *userAuthenticationController.RegisterController
		UserLoginController             *userAuthenticationController.LoginController
		UserAccessTokenController       *userAuthenticationController.AccessTokenController
		UserRecoveryPasswordController  *userAuthenticationController.RecoverPasswordController
		UserTwoFaAuthController         *userAuthenticationController.TwoFAController
		AdminLoginController            *adminAuthenticationController.LoginController
		AdminTwoFaAuthController        *adminAuthenticationController.TwoFAController
		AdminAccessTokenController      *adminController.AccessTokenController
		AdminRecoveryPasswordController *adminAuthenticationController.PasswordRecoveryController
		AuthenticationMiddleware        *middlewares.AuthenticationMiddleware
	}

	AdminContainer struct {
		AdminController *adminController.AdminController
	}

	AuthorizationContainer struct {
		AuthorizationMiddleware      *middlewares.AuthorizationMiddleware
		AdminAuthorizationController *adminController.AuthorizationController
	}

	UserContainer struct {
		AdminUserController *adminController.UserController
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
		ProvideAdminRepository,
		ProvideRoleRepository,
		// Services
		ProvideRegisterService,
		ProvideTwoFactorService,
		ProvideAdminService,
		ProvideLoginService,
		ProvideJwtService,
		ProvideUserService,
		ProvideAccessTokenService,
		ProvideRecoveryPasswordService,
		ProvideOTPService,
		// Controllers
		ProvideUserTwoFactorAutController,
		ProvideAdminTwoFactorAutController,
		ProvideUserRegisterController,
		ProvideUserLoginController,
		ProvideUserRecoverPasswordController,
		ProvideUserAccessTokenController,
		ProvideAdminLoginController,
		ProvideAdminAccessTokenController,
		ProvideAdminPasswordRecoveryController,
		// Middlewares
		ProvideAuthenticationMiddleware,
		wire.Struct(new(AuthenticationContainer), "*"),
	)
	return nil
}

func GetAdminContainer() *AdminContainer {
	wire.Build(
		database.GetInstance,
		ProvideAdminRepository,
		ProvideRoleRepository,
		ProvideAdminService,
		ProvideOTPService,
		ProvideAdminController,
		wire.Struct(new(AdminContainer), "*"),
	)
	return nil
}

func GetAuthorizationContainer() *AuthorizationContainer {
	wire.Build(
		// Repositories
		database.GetInstance,
		ProvideRoleRepository,
		ProvidePermissionRepository,
		ProvidePermissionGroupRepository,
		ProvideUserRepository,
		ProvideAdminRepository,
		// Services
		ProvideAuthorizationService,
		// Controllers
		ProvideAdminAuthorizationController,
		// Middlewares
		ProvideAuthorizationMiddleware,
		wire.Struct(new(AuthorizationContainer), "*"),
	)
	return nil
}

func GetUserContainer() *UserContainer {
	wire.Build(
		// Repositories
		database.GetInstance,
		ProvideUserRepository,
		// Services
		ProvideUserService,
		// Controllers
		ProvideAdminUserController,
		wire.Struct(new(UserContainer), "*"),
	)
	return nil
}
