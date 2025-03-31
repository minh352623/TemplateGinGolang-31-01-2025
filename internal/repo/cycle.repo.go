package repo

import (
	"ecom/global"
	"ecom/internal/model"
)

type ICycleRepository interface {
	GetCycleById(id int32) (model.Cycle, error)
}

type cycleRepository struct {
}

func NewCycleRepository() ICycleRepository {
	return &cycleRepository{}
}

func (r *cycleRepository) GetCycleById(id int32) (model.Cycle, error) {
	// ingnore column updated_at, created_at, user_created, user_updated
	cycle := model.Cycle{}
	err := global.PdbSetting.Select("key, value").Where("id = ?", id).First(&cycle).Error
	if err != nil {
		return model.Cycle{}, err
	}
	return cycle, nil
}
