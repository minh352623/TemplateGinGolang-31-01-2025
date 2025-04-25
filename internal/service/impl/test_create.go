package impl

import (
	"context"
	"ecom/internal/database"
)

type testCreateImpl struct {
	r *database.Queries
}

func NewTestCreateImpl(r *database.Queries) *testCreateImpl {
	return &testCreateImpl{
		r: r,
	}
}

func (t *testCreateImpl) CreateTest(req *database.CreateTestParams) (database.Test, error) {
	order, err := t.r.CreateTest(context.Background(), *req)
	if err != nil {
		return database.Test{}, err
	}
	return order, nil
}
