package inittiallize

import (
	"ecom/docs"
	"ecom/global"
	"ecom/internal/routers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	var r *gin.Engine
	docs.SwaggerInfo.Title = "Swagger Fortune Vault"
	docs.SwaggerInfo.Description = "This is a sample server Fortune Vault."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8001"
	docs.SwaggerInfo.BasePath = "/v1/api"
	docs.SwaggerInfo.Schemes = []string{"http"}
	if global.Config.Server.Mode == "dev" {
		gin.SetMode((gin.DebugMode))
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode((gin.ReleaseMode))
		r = gin.New()
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// middleware
	// r.Use(middlewares.AuthMiddleware())
	depositRouter := routers.RouterGroupApp.Deposit
	testRouter := routers.RouterGroupApp.Test
	MainGroup := r.Group("v1/api")
	{
		MainGroup.GET("checkStatus", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "ok"})
		})
	}
	{
		depositRouter.InitDepositRouter(MainGroup)
		testRouter.InitTestRouter(MainGroup)
	}

	return r
}
