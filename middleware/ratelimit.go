package middleware

import (
	"net/http"
	"time"

	ginpkg "github.com/gin-gonic/gin"

	"stockflow/component/ratelimit"
)

func RateLimit(limiter ratelimit.Limiter, limit int64, window time.Duration) ginpkg.HandlerFunc {
	return func(c *ginpkg.Context) {
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		clientID := c.ClientIP()
		if clientID == "" {
			clientID = "unknown"
		}

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		if !limiter.IsAllowed(c.Request.Context(), clientID, path, limit, window) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, ginpkg.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}
