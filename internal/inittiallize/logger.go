package inittiallize

import (
	"ecom/global"
	"ecom/pkg/logger"
)

func initLogger() {
	global.Logger = logger.NewLogger(global.Config.Logger)
}
