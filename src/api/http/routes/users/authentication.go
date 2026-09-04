package routes

import (
	"athena/src/providers"
	"athena/src/services"

	"github.com/gin-gonic/gin"
)

func AuthenticationRouter(router *gin.RouterGroup) {
	loginController := providers.GetAuthenticationContainer().UserLoginController
	registerController := providers.GetAuthenticationContainer().UserRegisterController
	recoverPasswordController := providers.GetAuthenticationContainer().UserRecoveryPasswordController
	twoFaController := providers.GetAuthenticationContainer().UserTwoFaAuthController
	authMiddleware := providers.GetAuthenticationContainer().AuthenticationMiddleware

	registerVerifyOTPRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.RegisterCriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("register-verify-otp"))

	registerResendOTPRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.CriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("register-resend-otp"))

	loginVerifyOTPRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.LoginCriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("login-verify-otp"))

	loginViaOTPRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.LoginCriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("login-via-otp"))

	loginViaPasswordRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.LoginCriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("user-login"))

	loginVerifyTwoFaRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.LoginCriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("user-login-2fa"))

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

	authentication := router.Group("authentication")

	register := authentication.Group("register")
	{
		register.POST("send-otp", registerController.Register)
		register.POST("verify-otp", registerVerifyOTPRateLimiter.Middleware, registerController.VerifyRegister)
		register.POST("resend-otp", registerResendOTPRateLimiter.Middleware, registerController.ResendOTP)
	}

	login := authentication.Group("login")
	{
		login.POST("", loginViaPasswordRateLimiter.Middleware, loginController.LoginViaPassword)
		login.POST("via-otp/send-otp", loginViaOTPRateLimiter.Middleware, loginController.LoginViaOtpSendOtp)
		login.POST("via-otp/resend-otp", loginVerifyOTPRateLimiter.Middleware, loginController.ResendOTP)
		login.POST("via-otp/verify-otp", loginViaOTPRateLimiter.Middleware, loginController.VerifyOTP)
		login.POST("verify-2fa", loginVerifyTwoFaRateLimiter.Middleware, loginController.VerifyTwoFa)
	}

	recoverPassword := authentication.Group("recover-password")
	{
		recoverPassword.POST("", recoverPasswordRateLimiter.Middleware, recoverPasswordController.RecoverPassword)
		recoverPassword.POST("verify-otp", recoverPasswordVerifyOTPRateLimiter.Middleware, recoverPasswordController.VerifyOtpRecoverPassword)
		recoverPassword.POST("set-password", recoverPasswordSetPasswordRateLimiter.Middleware, recoverPasswordController.SetPassword)
		recoverPassword.POST("resend-otp", recoverPasswordResendOTPRateLimiter.Middleware, recoverPasswordController.ResendOTP)
	}

	profile := router.Group("profile")
	profile.Use(authMiddleware.Middleware("user"))
	{
		twoFa := profile.Group("2fa")
		twoFa.GET("status", twoFaController.Status)
		twoFa.POST("enable", twoFaController.Enable)
		twoFa.POST("enable/verify", twoFaController.VerifyCode)
		twoFa.POST("disable", twoFaController.Disable)
	}
}
