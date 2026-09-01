package middlewares

import (
	"athena/src/api/http/response"
	"athena/src/pkg/logger"
	"athena/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthorizationMiddleware struct {
	AuthorizationService services.IAuthorizationService
}

// Middleware wraps the AuthorizationMiddleware method to make it compatible with Gin.
func (service *AuthorizationMiddleware) Middleware(permission string) gin.HandlerFunc {
	return func(context *gin.Context) {
		if !service.userIsAuthorized(context, permission) {
			response.Api(context).SetStatusCode(http.StatusForbidden).SetMessage("request-unauthorized").SetLog().Send()
			context.Abort()
			return
		}
		context.Next()
	}
}

func (service *AuthorizationMiddleware) userIsAuthorized(context *gin.Context, permission string) bool {

	inquiry, err := service.AuthorizationService.IsAuthorized(context,
		context.GetUint("authenticated-user-id"),
		context.GetString("authenticated-user-type"),
		permission,
	)

	if !inquiry || err != nil {
		logger.LogInfo(context, "user is not authorized to do this operation", service, zap.String("permission", permission))
		return false
	}

	return true
}
