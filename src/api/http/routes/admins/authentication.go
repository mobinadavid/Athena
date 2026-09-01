package admins

import (
	"athena/src/providers"
	"athena/src/services"

	"github.com/gin-gonic/gin"
)

func AuthenticationRouter(router *gin.RouterGroup) {
	loginController := providers.GetAuthenticationContainer().AdminLoginController

	// get rate limiter for login
	adminLoginViaPassRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.AdminLoginCriticalLimiter()).SetKey(services.GenericCriticalKeyGetter("admin-login-critical"))

	authentication := router.Group("authentication")

	login := authentication.Group("login")
	{
		login.POST("", adminLoginViaPassRateLimiter.Middleware, loginController.LoginViaPassword)
	}
}
