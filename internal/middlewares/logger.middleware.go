package middlewares

import "github.com/gin-gonic/gin"

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}