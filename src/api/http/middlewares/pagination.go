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

	sortBy := context.DefaultQuery("sort_by", "created_at")
	sortOrder := context.DefaultQuery("sort_order", "asc")

	// Set values in context
	context.Set("page", page)
	context.Set("limit", limit)
	context.Set("sort_by", sortBy)
	context.Set("sort_order", sortOrder)

	context.Next()
}
