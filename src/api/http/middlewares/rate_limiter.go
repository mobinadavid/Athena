package middlewares

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

// RateLimiter struct to hold the limiter and mutex
type RateLimiter struct {
	limiter *rate.Limiter
	mutex   sync.Mutex
}

// NewRateLimiter creates a new RateLimiter instance
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(r, b),
	}
}

// Middleware function to limit requests
func (rateLimiter *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rateLimiter.mutex.Lock()
		if !rateLimiter.limiter.Allow() {
			rateLimiter.mutex.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}
		rateLimiter.mutex.Unlock()
		c.Next()
	}
}
