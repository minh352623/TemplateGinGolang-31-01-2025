package deposit

import (
	controller "ecom/internal/controller/deposit"
	"ecom/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type DepositRouter struct{}

func (u *DepositRouter) InitDepositRouter(Router *gin.RouterGroup) {
	depositRouterPrivate := Router.Group("/deposit")
	depositRouterPrivate.Use(middlewares.AuthMiddleware())
	{
		depositRouterPrivate.POST("/:id", controller.Deposit.ProcessDeposit)
	}
}
