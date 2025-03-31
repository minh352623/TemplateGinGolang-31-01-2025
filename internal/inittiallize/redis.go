package inittiallize

import (
	"context"
	"ecom/global"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

func initRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.Config.Redis.Host, global.Config.Redis.Port),
		Password: global.Config.Redis.Password,
		DB:       global.Config.Redis.DB,
		PoolSize: global.Config.Redis.PoolSize,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		global.Logger.Error("Redis connect error", zap.Error(err))
	}

	fmt.Println("Redis connect success")

	global.Rdb = rdb
}
