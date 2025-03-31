package interest

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"ecom/internal/model"
	"ecom/internal/utils/interest"

	"github.com/stretchr/testify/assert"
)

func TestCalculateInterestDeposit(t *testing.T) {
	// Test Case 1: Normal Case
	t.Run("Normal case", func(t *testing.T) {
		// Mock data for balanceBeforeUpdate (this simulates the wallet data)
		balanceBeforeUpdate := model.Wallet{
			LastTimeUpdate: "1741686996",        // Mock last time update as UNIX timestamp
			TimeDeposit:    "1741591420",        // Mock time deposit as UNIX timestamp
			Balance:        "18.94508622222222", // Mock balance as a string
			IsNew:          false,               // Simulate a new wallet
			AmountInterest: "0",
			ProviderKey:    "123456abc",
		}

		// Mock settings for interest calculation
		settings := interest.InterestSetting{
			LockTimeDefault: interest.LockTimeDefault{
				PercentDefault:         10,
				PercentPrincipal:       100,
				LockTimeDefault:        300,
				PercentForAdminDefault: 0,
			},
			Percents: []interest.Percent{
				{
					Seconds: 300,
					Settings: []interest.PercentSetting{
						{
							PercentPrincipal: 50,
							LockTime:         300,
							PercentInterest:  10,
							PercentForAdmin:  0,
						},
					},
				},
			},
		}
		now := time.Now().Unix()

		lastTimeUpdate, amountInterestUpdate := interest.CalculateInterest(&balanceBeforeUpdate, &settings, now)
		fmt.Println("lastTimeUpdate", lastTimeUpdate, "amountInterestUpdate", amountInterestUpdate)

	})
}
func TestCalculateInterest(t *testing.T) {
	// Test Case 1: Normal Case
	t.Run("Normal case", func(t *testing.T) {
		// Mock data for balanceBeforeUpdate (this simulates the wallet data)
		balanceBeforeUpdate := model.Wallet{
			LastTimeUpdate: "1633014000", // Mock last time update as UNIX timestamp
			TimeDeposit:    "1633014000", // Mock time deposit as UNIX timestamp
			Balance:        "1000",       // Mock balance as a string
			IsNew:          true,         // Simulate a new wallet
		}

		// Mock settings for interest calculation
		settings := interest.InterestSetting{
			LockTimeDefault: interest.LockTimeDefault{
				PercentDefault:  5.0, // 5% default interest
				LockTimeDefault: 30,  // 30 seconds lock time for default interest

			},
			Percents: []interest.Percent{
				{
					Seconds: 30 * 24 * 60 * 60, // 30 days in seconds
					Settings: []interest.PercentSetting{
						{
							PercentPrincipal: 10.0, // 10% principal
							LockTime:         30,   // 30 seconds lock time
							PercentInterest:  5.0,  // 5% interest
						},
					},
				},
			},
		}

		// Test 1: Default case with the mocked data and current time
		lastTimeUpdateInt, _ := strconv.Atoi(balanceBeforeUpdate.LastTimeUpdate)
		nowTime := int64(lastTimeUpdateInt + 60) // Get current Unix timestamp
		lastTimeUpdate, amountInterestUpdate := interest.CalculateInterest(&balanceBeforeUpdate, &settings, nowTime)
		fmt.Println("lastTimeUpdate", lastTimeUpdate, "amountInterestUpdate", amountInterestUpdate)

		// Expected values (based on your business logic, replace with actual expected results)
		expectedLastTimeUpdate := nowTime
		expectedAmountInterest := 10.0 // Example expected interest

		// Assert that the actual and expected results are the same
		assert.Equal(t, expectedLastTimeUpdate, lastTimeUpdate, "LastTimeUpdate should match the current time")
		assert.Equal(t, expectedAmountInterest, amountInterestUpdate, "AmountInterest should match the expected interest")
	})

	// Test Case 2: Case when IsNew is false (should not calculate any interest)
	t.Run("Wallet is not new", func(t *testing.T) {
		balanceBeforeUpdate := model.Wallet{
			LastTimeUpdate: "1633014000",
			TimeDeposit:    "1633014000",
			Balance:        "1000",
			IsNew:          false, // Wallet is not new, so interest should not be calculated
		}

		settings := interest.InterestSetting{
			LockTimeDefault: interest.LockTimeDefault{
				PercentDefault:  5.0,
				LockTimeDefault: 30,
			},
			Percents: []interest.Percent{
				{
					Seconds: 30 * 24 * 60 * 60,
					Settings: []interest.PercentSetting{
						{
							PercentPrincipal: 10.0,
							LockTime:         30,
							PercentInterest:  5.0,
						},
					},
				},
			},
		}

		nowTime := time.Now().Unix()
		_, amountInterestUpdate := interest.CalculateInterest(&balanceBeforeUpdate, &settings, nowTime)

		// Assert that no interest is calculated when IsNew is false
		assert.Equal(t, float64(0), amountInterestUpdate, "AmountInterest should be 0 when wallet is not new")
	})

	// Test Case 3: Case when lastTimeUpdate equals nowTime (no time difference)
	t.Run("No time difference", func(t *testing.T) {
		balanceBeforeUpdate := model.Wallet{
			LastTimeUpdate: "1633014000",
			TimeDeposit:    "1633014000",
			Balance:        "1000",
			IsNew:          true,
		}

		settings := interest.InterestSetting{
			LockTimeDefault: interest.LockTimeDefault{
				PercentDefault:  5.0,
				LockTimeDefault: 30,
			},
			Percents: []interest.Percent{
				{
					Seconds: 30 * 24 * 60 * 60,
					Settings: []interest.PercentSetting{
						{
							PercentPrincipal: 10.0,
							LockTime:         30,
							PercentInterest:  5.0,
						},
					},
				},
			},
		}

		// Last time update equals current time
		lastTimeUpdateInt, _ := strconv.Atoi(balanceBeforeUpdate.LastTimeUpdate)
		nowTime := int64(lastTimeUpdateInt)
		_, amountInterestUpdate := interest.CalculateInterest(&balanceBeforeUpdate, &settings, nowTime)

		// Assert that no interest is calculated because no time has passed
		assert.Equal(t, float64(0), amountInterestUpdate, "AmountInterest should be 0 if no time has passed")
	})

	// Test Case 4: Case when balanceBeforeUpdate is invalid (empty balance)
	t.Run("Invalid balance", func(t *testing.T) {
		balanceBeforeUpdate := model.Wallet{
			LastTimeUpdate: "1633014000",
			TimeDeposit:    "1633014000",
			Balance:        "", // Invalid balance
			IsNew:          true,
		}

		settings := interest.InterestSetting{
			LockTimeDefault: interest.LockTimeDefault{
				PercentDefault:  5.0,
				LockTimeDefault: 30,
			},
			Percents: []interest.Percent{
				{
					Seconds: 30 * 24 * 60 * 60,
					Settings: []interest.PercentSetting{
						{
							PercentPrincipal: 10.0,
							LockTime:         30,
							PercentInterest:  5.0,
						},
					},
				},
			},
		}

		nowTime := time.Now().Unix()
		_, amountInterestUpdate := interest.CalculateInterest(&balanceBeforeUpdate, &settings, nowTime)

		// Assert that the function handles the error gracefully
		assert.Equal(t, float64(0), amountInterestUpdate, "AmountInterest should be 0 when balance is invalid")
	})
}
