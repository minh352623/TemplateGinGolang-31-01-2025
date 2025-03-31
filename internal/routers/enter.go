package routers

import (
	"ecom/internal/routers/deposit"
)

type RouterGroup struct {
	Deposit deposit.DepositRouterGroup
}

var RouterGroupApp = new(RouterGroup)
