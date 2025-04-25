//go:build wireinject

package wire

import (
	"ecom/internal/messaging"
	"ecom/internal/repo"
	"ecom/internal/service"

	"github.com/google/wire"
)

func InitializeConsumeHandler() (*messaging.ConsumeMessage, error) {
	wire.Build(
		service.NewTestService,
		repo.NewTestRepository,
		messaging.NewConsumeMessage,
	)
	return new(messaging.ConsumeMessage), nil
}
