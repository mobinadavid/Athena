package middlewares

import (
	"athena/src/api/errs"
	"athena/src/api/http/response"
	"athena/src/pkg/logger"
	"athena/src/services/authentication"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthenticationMiddleware struct {
	AccessTokenService authentication.IAccessTokenService
}

// Middleware wraps the AuthenticationMiddleware method to make it compatible with Gin.
func (service *AuthenticationMiddleware) Middleware(ownerType string) gin.HandlerFunc {
	return func(context *gin.Context) {
		// Check authentication and handle token expiry
		if !service.isAuthenticated(context, ownerType) {
			context.Abort()
			return
		}
		// Continue processing if authenticated
		context.Next()
	}
}

// Handle is a middleware function for Gin to authenticate requests.
// It checks the Authorization header for a valid JWT, validates it, and sets user context
// information if the authentication is successful.
func (service *AuthenticationMiddleware) isAuthenticated(context *gin.Context, ownerType string) bool {
	// Define the expected prefix for the Authorization header.
	const BearerSchema = "Bearer "

	// Retrieve the Authorization header from the request.
	header := context.GetHeader("Authorization")
	// Series of if conditions to validate the presence and format of the 'Authorization' header.
	// Each condition sets isAuthenticated to false if a specific check fails.

	// Check if the Authorization header is missing, does not start with 'Bearer ', or is just 'Bearer ' without a token.
	if header == "" || !strings.HasPrefix(header, BearerSchema) || len(header) == len(BearerSchema) {
		logger.LogInfo(context, "failed to get access token from header", service)
		response.Api(context).SetMessage(errs.ErrAuthenticationFailed.Error()).SetStatusCode(http.StatusUnauthorized).SetLog().Send()
		return false
	}

	// Extract the JWT token from the Authorization header, removing the 'Bearer ' prefix.
	tokenString := header[len(BearerSchema):]

	// Validate the JWT both JWT and database.
	token, err := service.AccessTokenService.Validate(context, tokenString, authentication.AccessToken, ownerType)
	if err != nil {
		logger.LogErrorWithFieldsV2(context, "failed to validate access token", service, err)
		response.Api(context).SetMessage(errs.ErrAuthenticationFailed.Error()).SetStatusCode(http.StatusUnauthorized).SetLog().Send()
		return false
	}
	// Token is valid, set user context
	context.Set("authenticated-user-id", token.OwnerID)
	context.Set("authenticated-user-type", token.OwnerType)
	context.Set("access-token-uuid", token.Uuid.String())

	// log the owner session
	logger.LogOwnerSession(token)

	// Update last used timestamp
	defer func() {
		_, _ = service.AccessTokenService.UpdateLastUsedAt(context, token)
	}()

	return true
}
