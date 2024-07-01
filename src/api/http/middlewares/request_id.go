package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID(context *gin.Context) {
	if context.Request != nil {
		context.Set("request-uuid", uuid.NewString())
	}
	context.Next()
}
