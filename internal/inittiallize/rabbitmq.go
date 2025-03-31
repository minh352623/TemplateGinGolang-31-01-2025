package inittiallize

import (
	"ecom/global"
	consts "ecom/pkg/const"
	"ecom/pkg/rabbitmq"

	"go.uber.org/zap"
)

func initRabbitMQ() {
	config := rabbitmq.Config{
		Host:     global.Config.RabbitMQ.Host,
		Port:     global.Config.RabbitMQ.Port,
		User:     global.Config.RabbitMQ.User,
		Password: global.Config.RabbitMQ.Password,
		VHost:    global.Config.RabbitMQ.VHost,
	}

	global.Logger.Info("Initializing RabbitMQ", zap.Any("config", config))

	rmq, err := rabbitmq.NewRabbitMQ(config, global.Logger.GetZapLogger())
	if err != nil {
		global.Logger.Fatal("Failed to initialize RabbitMQ", zap.Error(err))
	}

	// Declare default exchange
	err = rmq.DeclareExchange("ecom.events", "topic", true, false)
	if err != nil {
		global.Logger.Fatal("Failed to declare exchange", zap.Error(err))
	}

	// Declare consistent hash exchange
	err = rmq.DeclareConsistentHashExchange(consts.HashedExchangeName)
	if err != nil {
		global.Logger.Fatal("Failed to declare consistent hash exchange", zap.Error(err))
	}

	// Create producer
	producer := rabbitmq.NewProducer(rmq)

	// Store in global for access throughout the application
	global.RabbitMQ = rmq
	global.RabbitMQProducer = producer

	// Create queues with binding weights
	// queue, err := rmq.DeclareQueue("sync_queue", true, false, false, false, nil)

	// // Bind with weight
	// err = rmq.BindQueue(queue.Name, "1", consts.HashedExchangeName)
	// if err != nil {
	// 	global.Logger.Fatal("Failed to bind queue", zap.Error(err))
	// }

	global.Logger.Info("RabbitMQ initialized successfully")
}
