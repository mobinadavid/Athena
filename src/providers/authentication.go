package providers

import (
	userAuthcontrollers "athena/src/api/http/controllers/authentication"
	"athena/src/api/http/middlewares"
	"athena/src/services"
	authServices "athena/src/services/authentication"
)

func ProvideUserRegisterController(registerService *authServices.RegisterService) *userAuthcontrollers.RegisterController {
	return &userAuthcontrollers.RegisterController{
		RegisterService: registerService,
	}
}

func ProvideUserLoginController(loginService *authServices.LoginService) *userAuthcontrollers.LoginController {
	return &userAuthcontrollers.LoginController{
		LoginService: loginService,
	}
}

func ProvideUserRecoverPasswordController(recoverPasswordService *authServices.RecoveryPasswordService) *userAuthcontrollers.RecoverPasswordController {
	return &userAuthcontrollers.RecoverPasswordController{
		RecoverPasswordService: recoverPasswordService,
	}
}

func ProvideRegisterService(userService *services.UserService, otpService *services.OTPService) *authServices.RegisterService {
	return &authServices.RegisterService{
		UserService: userService,
		OTPService:  otpService,
	}
}

func ProvideLoginService(userService *services.UserService, jwtService *authServices.JwtService, accessTokenService *authServices.AccessTokenService, twoFaService *authServices.TwoFaService, otpService *services.OTPService) *authServices.LoginService {
	return &authServices.LoginService{
		UserService:        userService,
		JwtService:         jwtService,
		AccessTokenService: accessTokenService,
		OTPService:         otpService,
		TwoFaService:       twoFaService,
	}
}

func ProvideTwoFaService() *authServices.TwoFaService {
	return &authServices.TwoFaService{}
}

func ProvideJwtService() *authServices.JwtService {
	return &authServices.JwtService{}
}

func ProvideAuthenticationMiddleware(accessTokenService *authServices.AccessTokenService) *middlewares.AuthenticationMiddleware {
	return &middlewares.AuthenticationMiddleware{
		AccessTokenService: accessTokenService,
	}
}

func ProvideRecoveryPasswordService(OtpService *services.OTPService, UserService *services.UserService) *authServices.RecoveryPasswordService {
	return &authServices.RecoveryPasswordService{
		OTPService:  OtpService,
		UserService: UserService,
	}
}

//func ProvideAdminLoginController(loginService *authServices.LoginService, accessToken *authServices.AccessTokenService) *adminAuthcontrollers.LoginController {
//	return &adminAuthcontrollers.LoginController{
//		LoginService: loginService,
//		AccessToken:  accessToken,
//	}
//}
//
//func ProvideAdminPasswordRecoveryController(RecoveryPassword *authServices.RecoveryPasswordService) *adminAuthcontrollers.PasswordRecoveryController {
//	return &adminAuthcontrollers.PasswordRecoveryController{
//		RecoveryPasswordService: RecoveryPassword,
//	}
//}
