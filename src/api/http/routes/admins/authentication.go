package admins

import (
	"athena/src/providers"
	"athena/src/services"

	"github.com/gin-gonic/gin"
)

func AuthenticationRouter(router *gin.RouterGroup) {
	authContainer := providers.GetAuthenticationContainer()
	loginController := authContainer.AdminLoginController
	recoveryController := authContainer.AdminRecoveryPasswordController
	twoFaController := authContainer.AdminTwoFaAuthController

	adminLoginViaPassRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.AdminLoginCriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("admin-login-critical"))

	adminLoginTwoFaRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.AdminLoginCriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("admin-login-2fa"))

	recoverPasswordRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.CriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("admin-recover-password"))

	authentication := router.Group("authentication")
	login := authentication.Group("login")
	{
		login.POST("", adminLoginViaPassRateLimiter.Middleware, loginController.LoginViaPassword)
		login.POST("verify-2fa", adminLoginTwoFaRateLimiter.Middleware, loginController.VerifyTwoFa)
	}

	recoverPassword := authentication.Group("recover-password")
	{
		recoverPassword.POST("", recoverPasswordRateLimiter.Middleware, recoveryController.RequestOtp)
		recoverPassword.POST("verify-otp", recoverPasswordRateLimiter.Middleware, recoveryController.RecoveryPasswordViaOtp)
	}

	profile := router.Group("profile")
	profile.Use(authContainer.AuthenticationMiddleware.Middleware("admin"))
	{
		twoFa := profile.Group("2fa")
		twoFa.GET("status", twoFaController.Status)
		twoFa.POST("enable", twoFaController.Enable)
		twoFa.POST("enable/verify", twoFaController.VerifyCode)
		twoFa.POST("disable", twoFaController.Disable)
	}
}
