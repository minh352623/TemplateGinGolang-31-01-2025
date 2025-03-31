package middlewares

import "github.com/gin-gonic/gin"

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}