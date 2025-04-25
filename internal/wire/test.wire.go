//go:build wireinject

package wire

import (
	"ecom/internal/controller"
	"ecom/internal/repo"
	"ecom/internal/service"

	"github.com/google/wire"
)

func InitializeTestControllerHandler() (*controller.TestController, error) {
	wire.Build(
		controller.NewTestController,
		service.NewTestService,
		repo.NewTestRepository,
	)
	return new(controller.TestController), nil
}
