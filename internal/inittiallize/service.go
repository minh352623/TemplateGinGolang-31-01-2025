package inittiallize

import (
	"ecom/global"
	"ecom/internal/database"
	"ecom/internal/service/deposit"
	"ecom/internal/service/deposit/impl"
)

func InitServiceInterface() {
	query := database.New(global.Pdbc)
	deposit.InitDepositHandler(impl.NewDepositHandlerImpl(query))
	deposit.InitDepositInfo(impl.NewDepositInfoImpl(query))
}
