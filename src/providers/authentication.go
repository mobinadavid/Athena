package providers

import (
	adminAuthcontrollers "athena/src/api/http/controllers/admins/authentication"
	"athena/src/api/http/controllers/users/authentication"
	"athena/src/api/http/middlewares"
	"athena/src/services"
	authServices "athena/src/services/authentication"
)

func ProvideUserRegisterController(registerService *authServices.RegisterService) *authentication.RegisterController {
	return &authentication.RegisterController{
		RegisterService: registerService,
	}
}

func ProvideUserLoginController(loginService *authServices.LoginService) *authentication.LoginController {
	return &authentication.LoginController{
		LoginService: loginService,
	}
}

func ProvideUserRecoverPasswordController(recoverPasswordService *authServices.RecoveryPasswordService) *authentication.RecoverPasswordController {
	return &authentication.RecoverPasswordController{
		RecoverPasswordService: recoverPasswordService,
	}
}

func ProvideRegisterService(userService *services.UserService, otpService *services.OTPService) *authServices.RegisterService {
	return &authServices.RegisterService{
		UserService: userService,
		OTPService:  otpService,
	}
}

func ProvideLoginService(userService *services.UserService, jwtService *authServices.JwtService, accessTokenService *authServices.AccessTokenService, otpService *services.OTPService, adminService *services.AdminService, twoFaService *authServices.TwoFaService) *authServices.LoginService {
	return &authServices.LoginService{
		UserService:        userService,
		JwtService:         jwtService,
		AccessTokenService: accessTokenService,
		OTPService:         otpService,
		AdminService:       adminService,
		TwoFaService:       twoFaService,
	}
}

func ProvideJwtService() *authServices.JwtService {
	return &authServices.JwtService{}
}

func ProvideAuthenticationMiddleware(accessTokenService *authServices.AccessTokenService) *middlewares.AuthenticationMiddleware {
	return &middlewares.AuthenticationMiddleware{
		AccessTokenService: accessTokenService,
	}
}

func ProvideRecoveryPasswordService(OtpService *services.OTPService, UserService *services.UserService, AdminService *services.AdminService) *authServices.RecoveryPasswordService {
	return &authServices.RecoveryPasswordService{
		OTPService:   OtpService,
		UserService:  UserService,
		AdminService: AdminService,
	}
}

func ProvideAdminLoginController(loginService *authServices.LoginService, accessToken *authServices.AccessTokenService) *adminAuthcontrollers.LoginController {
	return &adminAuthcontrollers.LoginController{
		LoginService: loginService,
		AccessToken:  accessToken,
	}
}

func ProvideAdminPasswordRecoveryController(RecoveryPassword *authServices.RecoveryPasswordService) *adminAuthcontrollers.PasswordRecoveryController {
	return &adminAuthcontrollers.PasswordRecoveryController{
		RecoveryPasswordService: RecoveryPassword,
	}
}

func ProvideUserTwoFactorAutController(twoFactorAuth *authServices.TwoFaService) *authentication.TwoFAController {
	return &authentication.TwoFAController{
		TwoFaService: twoFactorAuth,
	}
}

func ProvideTwoFactorService(userService *services.UserService, adminService *services.AdminService, otpService *services.OTPService) *authServices.TwoFaService {
	return &authServices.TwoFaService{
		UserService:  userService,
		AdminService: adminService,
		OTPService:   otpService,
	}
}

func ProvideAdminTwoFactorAutController(twoFactorAuth *authServices.TwoFaService) *adminAuthcontrollers.TwoFAController {
	return &adminAuthcontrollers.TwoFAController{
		TwoFaService: twoFactorAuth,
	}
}
