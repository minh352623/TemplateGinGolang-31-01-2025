package middlewares

import "github.com/gin-gonic/gin"

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}