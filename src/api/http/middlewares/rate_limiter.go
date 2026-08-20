package middlewares

import (
	"athena/src/api/errs"
	"athena/src/api/http/response"
	"athena/src/services"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
)

// RateLimiterMiddleware is the middleware for gin.
type RateLimiterMiddleware struct {
	RateLimitService services.IRateLimitService
}

func (middleware *RateLimiterMiddleware) SetLimiter(limiter *limiter.Limiter) *RateLimiterMiddleware {
	middleware.RateLimitService.SetLimiter(limiter)
	return middleware
}

func (middleware *RateLimiterMiddleware) SetKey(keyGetter func(*gin.Context) string) *RateLimiterMiddleware {
	middleware.RateLimitService.SetKey(keyGetter)
	return middleware
}

func (middleware *RateLimiterMiddleware) SetExcludedKey(excludedKey func(string) bool) *RateLimiterMiddleware {
	middleware.RateLimitService.SetExcludedKey(excludedKey)
	return middleware
}

func (middleware *RateLimiterMiddleware) Middleware(c *gin.Context) {
	key := middleware.RateLimitService.KeyFunc(c)
	if middleware.RateLimitService.GetExcludedKey() != nil && middleware.RateLimitService.ExcludedKeyFunc(key) {
		c.Next()
		return
	}
	context, err := middleware.RateLimitService.GetLimiter().Get(c, key)
	if err != nil {
		response.Api(c).SetStatusCode(http.StatusForbidden).
			SetMessage(errs.SomeThingWentWrong.Error()).SetLog().Send()
		c.Abort()
		return
	}

	c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
	c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
	c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

	resetIn := math.Ceil(float64(context.Reset-time.Now().Unix()) / 60)
	if resetIn < 0 {
		resetIn = 0
	}

	// handle rate limit exceed or request contains captcha
	if context.Reached {
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead || c.Request.Method == http.MethodOptions || c.Request.Method == http.MethodDelete {
			if middleware.RateLimitService.GetOnExceedHandler(c, strconv.FormatFloat(resetIn, 'f', -1, 64)) {
				c.Next()
			}
			c.Abort()
			return
		} else if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch {
			if middleware.RateLimitService.PostOnExceedHandler(c, strconv.FormatFloat(resetIn, 'f', -1, 64)) {
				c.Next()
			}
			c.Abort()
			return
		}
	}

	c.Next()
}
