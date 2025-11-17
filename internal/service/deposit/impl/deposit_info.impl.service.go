package impl

import (
	"context"
	"ecom/internal/database"
	"ecom/internal/service/deposit/types"
)

type depositInfoImpl struct {
	r *database.Queries
}

func NewDepositInfoImpl(r *database.Queries) *depositInfoImpl {
	return &depositInfoImpl{
		r: r,
	}
}

func (d *depositInfoImpl) GetDepositById(ctx context.Context, id string) (*types.Deposit, error) {
	return &types.Deposit{
		ID: id,
	}, nil
}
