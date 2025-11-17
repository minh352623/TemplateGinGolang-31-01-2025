package deposit

import (
	"context"
	"ecom/internal/service/deposit/types"
)

type (
	IDepositHandler interface {
		ProcessDeposit(ctx context.Context, deposit *types.Deposit) (*types.Deposit, error)
	}
	IDepositInfo interface {
		GetDepositById(ctx context.Context, id string) (*types.Deposit, error)
	}
)

var (
	localDepositHandler IDepositHandler
	localDepositInfo    IDepositInfo
)

func InitDepositHandler(handler IDepositHandler) {
	localDepositHandler = handler
}

func InitDepositInfo(input IDepositInfo) {
	localDepositInfo = input
}

func DepositHandler() IDepositHandler {
	if localDepositHandler == nil {
		panic("localDepositHandler is nil")
	}
	return localDepositHandler
}

func DepositInfo() IDepositInfo {
	if localDepositInfo == nil {
		panic("localDepositInfo is nil")
	}
	return localDepositInfo
}
