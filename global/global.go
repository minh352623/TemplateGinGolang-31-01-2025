package global

import (
	"ecom/pkg/logger"
	"ecom/pkg/rabbitmq"
	"ecom/pkg/security"
	"ecom/pkg/setting"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Config           setting.Config
	Logger           *logger.LoggerZap
	Pdb              *gorm.DB
	PdbSetting       *gorm.DB
	Rdb              *redis.Client
	ServerSetting    *setting.ServerSetting
	SecuritySetting  *setting.SecuritySetting
	SecurityService  *security.SecurityService
	RabbitMQ         *rabbitmq.RabbitMQ
	RabbitMQProducer *rabbitmq.Producer
)
