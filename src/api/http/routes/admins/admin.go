package admins

import (
	"athena/src/api/http/middlewares"
	"athena/src/models"
	"athena/src/providers"
	"athena/src/services"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRouter(router *gin.RouterGroup) {
	authenticationContainer := providers.GetAuthenticationContainer()
	authorizationContainer := providers.GetAuthorizationContainer()
	adminContainer := providers.GetAdminContainer()

	changePasswordRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.RegisterCriticalLimiter()).SetKey(services.CriticalChangePasswordKeySetter)
	changePasswordSendOTPRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.RegisterCriticalLimiter()).SetKey(services.CriticalChangePasswordSendOtpKeySetter)
	changePasswordResendOTPRateLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.RegisterCriticalLimiter()).SetKey(services.CriticalChangePasswordResendOtpKeySetter)

	admin := router.Group("admins")
	admin.Use(authenticationContainer.AuthenticationMiddleware.Middleware("admin"))
	admin.GET("",
		authorizationContainer.AuthorizationMiddleware.Middleware("admin-list"),
		middlewares.QueryParametersBuilderMiddleware(models.AdminModel{}),
		adminContainer.AdminController.GetList,
	)

	admin.GET(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("admin-show"),
		adminContainer.AdminController.GetByUuid,
	)

	router.POST("admins/get-profile",
		authenticationContainer.AuthenticationMiddleware.Middleware("admin"),
		authorizationContainer.AuthorizationMiddleware.Middleware("admin-show"),
		adminContainer.AdminController.GetProfile,
	)

	admin.POST("",
		authorizationContainer.AuthorizationMiddleware.Middleware("admin-create"),
		adminContainer.AdminController.Create,
	)

	admin.PATCH(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("admin-update"),
		adminContainer.AdminController.Update,
	)

	admin.DELETE(":uuid",
		authorizationContainer.AuthorizationMiddleware.Middleware("admin-delete"),
		adminContainer.AdminController.Delete,
	)

	changePasswordRoutes := admin.Group("change-password")
	{
		changePasswordRoutes.POST("send-otp",
			changePasswordSendOTPRateLimiter.Middleware, adminContainer.AdminController.ChangePasswordSendOtp)
		changePasswordRoutes.POST("",
			changePasswordRateLimiter.Middleware, adminContainer.AdminController.ChangePasswordVerifyOtp)
		changePasswordRoutes.POST("resend-otp",
			changePasswordResendOTPRateLimiter.Middleware, adminContainer.AdminController.ChangePasswordResendOTP)
	}
}
