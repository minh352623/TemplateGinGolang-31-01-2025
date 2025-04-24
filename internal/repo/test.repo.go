package repo

import (
	"context"
	"database/sql"
	"ecom/global"
	"ecom/internal/database"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ITestRepository interface {
	GetTestById(id uuid.UUID) (database.Test, error)
	CreateTest(req *database.CreateTestParams) (database.Test, error)
	UpdateTest(req *database.UpdateTestParams) (database.Test, error)
}

type testRepository struct {
	sqlc *database.Queries
}

func NewTestRepository() ITestRepository {
	return &testRepository{
		sqlc: database.New(global.Pdbc),
	}
}

func (r *testRepository) GetTestById(id uuid.UUID) (database.Test, error) {
	order, err := r.sqlc.GetTestById(context.Background(), id)
	if err != nil {
		return database.Test{}, err
	}
	return order, nil
}

func (r *testRepository) CreateTest(req *database.CreateTestParams) (database.Test, error) {
	order, err := r.sqlc.CreateTest(context.Background(), *req)
	if err != nil {
		return database.Test{}, err
	}
	return order, nil
}

func (r *testRepository) UpdateTest(req *database.UpdateTestParams) (database.Test, error) {
	// check if id is exist
	id := uuid.MustParse(req.ID.String())
	fmt.Println("id", id)
	test, err := r.sqlc.GetTestById(context.Background(), id)
	if err != nil {
		global.Logger.Error("GetTestById", zap.Error(err))
		return database.Test{}, err
	}
	if test.ID == uuid.Nil {
		global.Logger.Error("GetTestById", zap.Error(errors.New("id not found")))
		return database.Test{}, errors.New("id not found")
	}
	balance, err := strconv.ParseFloat(test.Balance.String, 64)
	if err != nil {
		global.Logger.Error("ParseFloat", zap.Error(err))
		return database.Test{}, err
	}
	balanceFloat, err := strconv.ParseFloat(req.Balance.String, 64)
	if err != nil {
		global.Logger.Error("ParseFloat", zap.Error(err))
		return database.Test{}, err
	}
	balance += balanceFloat
	req.Balance = sql.NullString{
		String: strconv.FormatFloat(balance, 'f', -1, 64),
		Valid:  true,
	}
	order, err := r.sqlc.UpdateTest(context.Background(), *req)
	if err != nil {
		global.Logger.Error("UpdateTest", zap.Error(err))
		return database.Test{}, err
	}
	return order, nil
}
