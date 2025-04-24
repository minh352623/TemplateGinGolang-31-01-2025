package service

import (
	"ecom/internal/database"
	"ecom/internal/repo"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ITestService interface {
	GetTestById(id uuid.UUID) (database.Test, error)
	CreateTest(req *database.CreateTestParams) (database.Test, error)
	UpdateTest(req *database.UpdateTestParams) (database.Test, error)
}

type testService struct {
	repo repo.ITestRepository
}

func NewTestService(repo repo.ITestRepository) ITestService {
	return &testService{
		repo: repo,
	}
}

func (s *testService) GetTestById(id uuid.UUID) (database.Test, error) {
	data, err := s.repo.GetTestById(id)
	if err != nil {
		fmt.Println("error", err)
		return database.Test{}, err
	}
	return data, nil
}

func (s *testService) CreateTest(req *database.CreateTestParams) (database.Test, error) {
	data, err := s.repo.CreateTest(req)
	if err != nil {
		fmt.Println("error", err)
		return database.Test{}, err
	}
	return data, nil
}

func (s *testService) UpdateTest(req *database.UpdateTestParams) (database.Test, error) {
	fmt.Println("UpdateTest", req)
	data, err := s.repo.UpdateTest(req)
	if err != nil {
		fmt.Println("error", err)
		return database.Test{}, err
	}
	time.Sleep(10 * time.Second)
	return data, nil
}
