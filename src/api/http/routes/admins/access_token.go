package admins

import (
	"athena/src/api/http/middlewares"
	"athena/src/models"
	"athena/src/providers"

	"github.com/gin-gonic/gin"
)

func RegisterAccessTokensRouter(router *gin.RouterGroup) {
	authContainer := providers.GetAuthenticationContainer()
	accessTokens := router.Group("access-tokens")
	{
		accessTokens.GET("",
			authContainer.AuthenticationMiddleware.Middleware("admin"),
			middlewares.QueryParametersBuilderMiddleware(models.AccessTokenModel{}),
			authContainer.AdminAccessTokenController.GetList,
		)

		accessTokens.GET(":uuid",
			authContainer.AuthenticationMiddleware.Middleware("admin"),
			authContainer.AdminAccessTokenController.GetByUuid,
		)

		accessTokens.POST("refresh",
			authContainer.AdminAccessTokenController.RefreshAccessToken,
		)

		revoke := accessTokens.Group("revoke").Use(authContainer.AuthenticationMiddleware.Middleware("admin"))
		{
			revoke.DELETE("",
				authContainer.AdminAccessTokenController.RevokeTokens,
			)

			revoke.DELETE(":uuid",
				authContainer.AdminAccessTokenController.RevokeTokenByUuid,
			)

			revoke.DELETE("current-token",
				authContainer.AdminAccessTokenController.RevokeCurrentToken,
			)
		}
	}

	activeAccessTokens := router.Group("active-access-tokens")
	{
		activeAccessTokens.GET("", authContainer.AuthenticationMiddleware.Middleware("admin"),
			middlewares.QueryParametersBuilderMiddleware(models.AccessTokenModel{}),
			authContainer.AdminAccessTokenController.GetActiveTokens,
		)
	}
}
