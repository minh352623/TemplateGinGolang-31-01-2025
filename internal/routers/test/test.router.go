package test

import (
	"ecom/internal/wire"

	"github.com/gin-gonic/gin"
)

type TestRouter struct{}

func (u *TestRouter) InitTestRouter(Router *gin.RouterGroup) {
	testController, err := wire.InitializeTestControllerHandler()
	if err != nil {
		panic(err)
	}

	testRouterPrivate := Router.Group("/test")
	// testRouterPrivate.Use(middlewares.AuthMiddleware())
	{
		testRouterPrivate.GET("/:id", testController.GetTestById)
		testRouterPrivate.POST("/update", testController.UpdateTest)
	}
}
