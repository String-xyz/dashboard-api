package store

import (
	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	"github.com/String-xyz/platform-admin-api/env"
)

func NewRedis() database.RedisStore {
	opts := database.RedisConfigOptions{
		Host:        env.Var.REDIS_HOST,
		Port:        env.Var.REDIS_PORT,
		Password:    env.Var.REDIS_PASSWORD,
		ClusterMode: !common.IsLocalEnv(),
	}
	return database.NewRedisStore(opts)
}
