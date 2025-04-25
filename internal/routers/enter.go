package routers

import (
	"ecom/internal/routers/deposit"
	"ecom/internal/routers/test"
)

type RouterGroup struct {
	Deposit deposit.DepositRouterGroup
	Test    test.TestRouterGroup
}

var RouterGroupApp = new(RouterGroup)
