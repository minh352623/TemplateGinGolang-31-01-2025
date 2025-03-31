package inittiallize

import (
	"ecom/global"
	"fmt"

	"go.uber.org/zap"
)

// swagger embed files

func Run() {
	LoadConfig()
	fmt.Println("config database", global.Config.Postgres.Host, global.Config.Postgres.Port, global.Config.Postgres.User, global.Config.Postgres.DBName)
	initLogger()
	initSecurity()
	initPostgres()
	initPostgresSetting()
	global.Logger.Info("hello world", zap.String("name", "John"))

	initRedis()
	// initRabbitMQ()

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
