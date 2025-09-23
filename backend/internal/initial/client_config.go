package initial

import (
	"github.com/baimhons/stadiumhub/internal"
	"github.com/baimhons/stadiumhub/internal/utils"
	"gorm.io/gorm"
)

type clientConfig struct {
	DB    *gorm.DB
	Redis utils.RedisClient
}

func newClientConfig() *clientConfig {
	db := utils.ConnectMySQLDatabase(
		internal.ENV.Database.Host,
		internal.ENV.Database.Port,
		internal.ENV.Database.User,
		internal.ENV.Database.Password,
		internal.ENV.Database.Name,
	)

	redis := utils.ConnectRedis(
		internal.ENV.Redis.Host,
		internal.ENV.Redis.Port,
		internal.ENV.Redis.Password,
	)

	redisWrapper := utils.NewRedisClient(redis)

	return &clientConfig{
		DB:    db,
		Redis: redisWrapper,
	}
}
