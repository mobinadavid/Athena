package middlewares

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func PaginationMiddleware(context *gin.Context) {
	page, _ := strconv.Atoi(context.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(context.DefaultQuery("limit", "10"))
	// Adjust limit to avoid excessively large requests
	switch {
	case limit > 100:
		limit = 100
	case limit <= 0:
		limit = 10
	}
	// Set values in context
	context.Set("page", page)
	context.Set("limit", limit)

	context.Next()
}
