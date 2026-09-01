package routes

import (
	"athena/src/providers"
	"athena/src/services"

	"github.com/gin-gonic/gin"
)

func AuthenticationRouter(router *gin.RouterGroup) {
	// get container
	loginController := providers.GetAuthenticationContainer().UserLoginController
	registerController := providers.GetAuthenticationContainer().UserRegisterController
	recoverPasswordController := providers.GetAuthenticationContainer().UserRecoveryPasswordController

	// get rate limiter for register
	registerVerifyOTPRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.RegisterCriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("register-verify-otp"))

	registerResendOTPRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.CriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("register-resend-otp"))

	// get rate limiter for login
	loginVerifyOTPRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.LoginCriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("login-verify-otp"))

	loginViaOTPRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.LoginCriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("login-via-otp"))

	// get rate limiter for recover password
	recoverPasswordRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.CriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("register"))

	recoverPasswordVerifyOTPRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.CriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("register-verify-otp"))

	recoverPasswordSetPasswordRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.CriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("register-set-password"))

	recoverPasswordResendOTPRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.CriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("register-resend-otp"))

	// define route
	authentication := router.Group("authentication")

	// register
	register := authentication.Group("register")
	{
		register.POST("send-otp", registerController.Register)
		register.POST("verify-otp", registerVerifyOTPRateLimiter.Middleware, registerController.VerifyRegister)
		register.POST("resend-otp", registerResendOTPRateLimiter.Middleware, registerController.ResendOTP)
	}

	// login
	login := authentication.Group("login")
	{
		login.POST("", loginVerifyOTPRateLimiter.Middleware, loginController.LoginViaPassword)
		login.POST("via-otp/send-otp", loginViaOTPRateLimiter.Middleware, loginController.LoginViaOtpSendOtp)
		login.POST("via-otp/resend-otp", loginVerifyOTPRateLimiter.Middleware, loginController.ResendOTP)
		login.POST("via-otp/verify-otp", loginViaOTPRateLimiter.Middleware, loginController.VerifyOTP)

	}

	// recover password
	recoverPassword := authentication.Group("recover-password")
	{
		recoverPassword.POST("", recoverPasswordRateLimiter.Middleware, recoverPasswordController.RecoverPassword)
		recoverPassword.POST("verify-otp", recoverPasswordVerifyOTPRateLimiter.Middleware, recoverPasswordController.VerifyOtpRecoverPassword)
		recoverPassword.POST("set-password", recoverPasswordSetPasswordRateLimiter.Middleware, recoverPasswordController.SetPassword)
		recoverPassword.POST("resend-otp", recoverPasswordResendOTPRateLimiter.Middleware, recoverPasswordController.ResendOTP)
	}
}
