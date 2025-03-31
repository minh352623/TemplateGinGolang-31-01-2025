package interest

import (
	"ecom/internal/model"
	"fmt"
	"strconv"
	"time"
)

type PercentSetting struct {
	PercentPrincipal float64 `json:"percentPrincipal"`
	LockTime         int64   `json:"lockTime"`
	PercentInterest  float64 `json:"percentInterest"`
	PercentForAdmin  float64 `json:"percentForAdmin"`
}

type Percent struct {
	Seconds  int              `json:"seconds"`
	Settings []PercentSetting `json:"settings"`
}

type LockTimeDefault struct {
	PercentDefault         float64 `json:"percentDefault"`
	PercentPrincipal       float64 `json:"percentPrincipal"`
	LockTimeDefault        int64   `json:"lockTimeDefault"`
	PercentForAdminDefault float64 `json:"percentForAdminDefault"`
}

type FeeSettingItem struct {
	TransactionType model.TransactionType `json:"transactionType"`
	FeeFixed        float64               `json:"feeFixed"`
	FeePercent      float64               `json:"feePercent"`
	FeeIn           bool                  `json:"feeIn"`
}
type CurrencySupportWithdraw struct {
	CurrencySupportInput  []model.Currency `json:"currencySupportInput"`
	CurrencySupportOutput []model.Currency `json:"currencySupportOutput"`
}
type Deposit struct {
	CurrencySupportDeposit []model.Currency `json:"currencySupportDeposit"`
	MinDeposit             float64          `json:"minDeposit"`
}
type Withdrawn struct {
	CurrencySupportInput  []model.Currency `json:"currencySupportInput"`
	CurrencySupportOutput []model.Currency `json:"currencySupportOutput"`
	MinWithdrawn          float64          `json:"minWithdrawn"`
	MaxWithdrawn          float64          `json:"maxWithdrawn"`
}
type InterestSetting struct {
	ID                         int32            `json:"id"`
	LockTimeDefault            LockTimeDefault  `json:"lockTimeDefault"`
	Percents                   []Percent        `json:"percents"`
	DepositLockTime            int              `json:"depositLockTime"`
	UrgentWithdrawalFeePercent float64          `json:"urgentWithdrawalFeePercent"`
	FreeWithdrawalsCount       int              `json:"freeWithdrawalsCount"`
	FreeWithdrawalLimit        float64          `json:"freeWithdrawalLimit"`
	CurrencyMergeEnabled       bool             `json:"currencyMergeEnabled"`
	FeeSetting                 []FeeSettingItem `json:"feeSetting"`
	IsAutoTakeProfit           bool             `json:"isAutoTakeProfit"`
	ProfitTakingCycle          int              `json:"profitTakingCycle"`
	Cronjob                    string           `json:"cronjob"`
	Deposit                    Deposit          `json:"deposit"`
	Withdrawn                  Withdrawn        `json:"withdrawn"`
	Platform                   string           `json:"platform"`
}

func CalculateInterest(balanceBeforeUpdate *model.Wallet, settings *InterestSetting, nowTime int64) (lastTimeUpdate int64, amountInterestUpdate float64) {
	//
	now := time.Now().Unix()
	if nowTime > 0 {
		now = nowTime
	}

	closeTime := now

	fmt.Println("closeTime", closeTime)
	percentDefault, lockTimeDefault, percentPrincipalDefault := settings.LockTimeDefault.PercentDefault, settings.LockTimeDefault.LockTimeDefault, settings.LockTimeDefault.PercentPrincipal
	percents := settings.Percents
	lastTimeUpdate, err := strconv.ParseInt(balanceBeforeUpdate.LastTimeUpdate, 10, 64)
	if err != nil {
		return 0, 0
	}
	timeDeposit, err := strconv.ParseInt(balanceBeforeUpdate.TimeDeposit, 10, 64)
	if err != nil {
		return 0, 0
	}
	balance, err := strconv.ParseFloat(balanceBeforeUpdate.Balance, 64)
	if err != nil {
		return 0, 0
	}
	amountInterest := 0.0

	// Calculate interest for each period
	if balanceBeforeUpdate.IsNew {
		for _, percent := range percents {
			periodEndTime := timeDeposit + int64(percent.Seconds)
			if closeTime != 0 && closeTime > lastTimeUpdate {
				timeCalculate := min(closeTime, periodEndTime)
				diffTime := timeCalculate - lastTimeUpdate
				if diffTime > 0 {
					for _, setting := range percent.Settings {
						numberBlock := float64(diffTime) / float64(setting.LockTime)
						amountClaimed := numberBlock *
							(balance * (setting.PercentPrincipal / 100)) *
							(setting.PercentInterest / 100)

						amountInterest += amountClaimed
					}
					lastTimeUpdate = timeCalculate
				}

			}
		}
	}

	// Apply the default interest after the periods
	fmt.Println("lastTimeUpdate", lastTimeUpdate, closeTime)
	if closeTime > lastTimeUpdate {
		diffTime := closeTime - lastTimeUpdate
		numberBlock := float64(diffTime) / float64(lockTimeDefault)
		if numberBlock > 0.0 {
			amountClaimed := numberBlock * balance * (percentPrincipalDefault / 100) * (percentDefault / 100)
			amountInterest += amountClaimed
			lastTimeUpdate = closeTime
		}
	}

	return closeTime, amountInterest
}

func CalculateOutput(now int64, timteDeposit int64, stepTime int64) int64 {
	fmt.Println("CalculateOutput", now, timteDeposit, stepTime)

	elapsedTime := now - timteDeposit
	remainingTime := elapsedTime % stepTime

	// Handle the conditional subtraction of stepTime
	var extraTime int64 = 0
	if remainingTime == 0 {
		extraTime = stepTime
	}

	output := now - remainingTime - extraTime
	return output
}
