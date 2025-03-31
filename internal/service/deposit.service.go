package service

import (
	"ecom/internal/model"
	"ecom/internal/repo"
)

type IDepositService interface {
	Test(userID string, email string, messageID string, routingKey string, hashKey string) (model.Cycle, error)
}

type depositService struct {
	cycleRepository repo.ICycleRepository
}

func NewDepositService(cycleRepository repo.ICycleRepository) IDepositService {
	return &depositService{
		cycleRepository: cycleRepository,
	}
}

func (ds *depositService) Test(userID string, email string, messageID string, routingKey string, hashKey string) (model.Cycle, error) {
	cycle, err := ds.cycleRepository.GetCycleById(1)
	if err != nil {
		return model.Cycle{}, err
	}

	return cycle, nil
}
