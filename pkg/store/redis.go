package store

import (
	"github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	"github.com/String-xyz/platform-admin-api/config"
)

func NewRedis() database.RedisStore {
	opts := database.RedisConfigOptions{
		Host:        config.Var.REDIS_HOST,
		Port:        config.Var.REDIS_PORT,
		Password:    config.Var.REDIS_PASSWORD,
		ClusterMode: !common.IsLocalEnv(),
	}
	return database.NewRedisStore(opts)
}
