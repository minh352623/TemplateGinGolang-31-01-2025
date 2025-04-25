package inittiallize

import (
	"ecom/global"
	"ecom/pkg/rabbitmq"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func InitRabbitMQ() {
	connectUrl := fmt.Sprintf("amqp://%s:%s@%s:%d/", global.Config.RabbitMQ.User, global.Config.RabbitMQ.Password, global.Config.RabbitMQ.Host, global.Config.RabbitMQ.Port)
	fmt.Println("connectUrl", connectUrl)
	conn, err := amqp.Dial(connectUrl)
	if err != nil {
		global.Logger.Error("Failed to connect to RabbitMQ", zap.Error(err))
		panic(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		global.Logger.Error("Failed to open a channel", zap.Error(err))
		panic(err)
	}

	global.RabbitMQManager = &rabbitmq.QueueManager{
		Conn:    conn,
		Channel: channel,
		Queues:  make(map[string]amqp.Queue),
	}

	// Declare exchange
	global.RabbitMQManager.DeclareExchange(global.Config.Exchange.Test, "x-consistent-hash")

	// Declare queues
	global.RabbitMQManager.DeclareQueue(global.Config.Queue.Test)

	// Bind queues to exchange
	err = global.RabbitMQManager.BindQueue(global.Config.Queue.Test, global.Config.Exchange.Test)
	if err != nil {
		global.Logger.Error("Failed to bind queue to exchange", zap.Error(err))
		panic(err)
	}

	fmt.Println("Connected to RabbitMQ successfully")

}
