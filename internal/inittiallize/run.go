package inittiallize

import (
	"ecom/global"
	"ecom/internal/wire"
	"fmt"

	"go.uber.org/zap"
)

// swagger embed files

func Run() {
	LoadConfig()
	fmt.Println("config database", global.Config.Postgres.Host, global.Config.Postgres.Port, global.Config.Postgres.User, global.Config.Postgres.DBName)
	initLogger()
	initSecurity()
	initPostgresC()
	initPostgresSetting()
	global.Logger.Info("hello world", zap.String("name", "John"))

	initRedis()
	InitRabbitMQ()
	consumeMessage, err := wire.InitializeConsumeHandler()
	if err != nil {
		global.Logger.Error("Failed to initialize consume message", zap.Error(err))
		return
	}
	consumeMessage.RegisterConsumers()
	// Start worker in a separate goroutine
	// go func() {
	// 	w := worker.NewWorker()
	// 	w.Start()
	// }()

	r := InitRouter()
	InitCronJob()

	port := global.Config.Server.Port
	fmt.Println("port", port)
	r.Run(":" + port)

}
