package routes

import (
	"athena/src/api/http/middlewares"
	"athena/src/models"
	"athena/src/providers"

	"github.com/gin-gonic/gin"
)

func RegisterAccessTokenRouter(router *gin.RouterGroup) {
	authContainer := providers.GetAuthenticationContainer()
	accessTokens := router.Group("access-tokens")
	{
		accessTokens.GET("",
			authContainer.AuthenticationMiddleware.Middleware("user"),
			middlewares.QueryParametersBuilderMiddleware(models.AccessTokenModel{}),
			authContainer.UserAccessTokenController.GetList,
		)

		accessTokens.GET(":uuid",
			authContainer.AuthenticationMiddleware.Middleware("user"),
			authContainer.UserAccessTokenController.GetByUuid,
		)

		accessTokens.POST("refresh",
			authContainer.UserAccessTokenController.RefreshAccessToken,
		)

		revoke := accessTokens.Group("revoke").Use(authContainer.AuthenticationMiddleware.Middleware("user"))
		{
			revoke.DELETE("",
				authContainer.UserAccessTokenController.RevokeTokens,
			)

			revoke.DELETE(":uuid",
				authContainer.UserAccessTokenController.RevokeTokenByUUID,
			)

			revoke.DELETE("current-token",
				authContainer.UserAccessTokenController.RevokeCurrentToken,
			)
		}
	}

	activeAccessTokens := router.Group("active-access-tokens")
	{
		activeAccessTokens.GET("", authContainer.AuthenticationMiddleware.Middleware("user"),
			middlewares.QueryParametersBuilderMiddleware(models.AccessTokenModel{}),
			authContainer.UserAccessTokenController.GetActiveTokens,
		)
	}
}
