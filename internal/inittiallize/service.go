package inittiallize

import (
	"ecom/global"
	"ecom/internal/database"
	"ecom/internal/service"
	"ecom/internal/service/impl"
)

func InitServiceInterface() {
	query := database.New(global.Pdbc)
	service.InitTestCreate(impl.NewTestCreateImpl(query))
}
