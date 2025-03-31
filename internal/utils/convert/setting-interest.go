package convert

import (
	"ecom/internal/model"
	"ecom/internal/utils/interest"
	"fmt"
)

func ConvertSettingInterest(settingInterest model.WalletIntegration) *interest.InterestSetting {
	fmt.Println("settingInterest", settingInterest)
	return &interest.InterestSetting{}
}
