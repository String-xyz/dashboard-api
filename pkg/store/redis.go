package store

import (
	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	"github.com/String-xyz/platform-admin-api/env"
)

func NewRedis() database.RedisStore {
	host, _ := env.Get("REDIS_HOST")
	port, _ := env.Get("REDIS_PORT")
	password, _ := env.Get("REDIS_PASSWORD")
	opts := database.RedisConfigOptions{
		Host:        host,
		Port:        port,
		Password:    password,
		ClusterMode: !common.IsLocalEnv(),
	}
	return database.NewRedisStore(opts)
}
