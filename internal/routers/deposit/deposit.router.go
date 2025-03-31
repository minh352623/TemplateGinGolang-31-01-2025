package deposit

import (
	"ecom/internal/middlewares"
	"ecom/internal/wire"

	"github.com/gin-gonic/gin"
)

type DepositRouter struct{}

func (u *DepositRouter) InitDepositRouter(Router *gin.RouterGroup) {
	depositController, err := wire.InitializeDepositHandler()
	if err != nil {
		panic(err)
	}

	depositRouterPrivate := Router.Group("/deposit")
	depositRouterPrivate.Use(middlewares.AuthMiddleware())
	{
		// depositRouterPrivate.POST("", depositController.Deposit)
		depositRouterPrivate.POST("/test", depositController.Test)
	}
}
