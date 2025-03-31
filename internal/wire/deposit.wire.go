//go:build wireinject

package wire

import (
	"ecom/internal/controller"
	"ecom/internal/repo"
	"ecom/internal/service"

	"github.com/google/wire"
)

func InitializeDepositHandler() (*controller.DepositController, error) {
	wire.Build(
		service.NewDepositService,
		controller.NewDepositController,
		repo.NewCycleRepository,
	)
	return new(controller.DepositController), nil
}
