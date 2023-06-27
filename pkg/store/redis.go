package store

import (
	"github.com/String-xyz/dashboard-api/config"
	"github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
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
