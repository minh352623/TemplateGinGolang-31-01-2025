package impl

import (
	"context"
	"ecom/internal/database"
	"ecom/internal/service/deposit/types"
)

type depositHandlerImpl struct {
	r *database.Queries
}

func NewDepositHandlerImpl(r *database.Queries) *depositHandlerImpl {
	return &depositHandlerImpl{
		r: r,
	}
}

func (d *depositHandlerImpl) ProcessDeposit(ctx context.Context, deposit *types.Deposit) (*types.Deposit, error) {
	return &types.Deposit{
		ID: deposit.ID,
	}, nil
}
