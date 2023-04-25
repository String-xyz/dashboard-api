package env

import (
	"errors"
	"os"

	"github.com/String-xyz/go-lib/common"
	"github.com/joho/godotenv"
)

var keys = []string{
	"ENV",
	"PORT",
	"BASE_DASHBOARD_URL",
	"DB_NAME",
	"DB_USERNAME",
	"DB_PASSWORD",
	"DB_HOST",
	"DB_PORT",
	"REDIS_PASSWORD",
	"REDIS_HOST",
	"REDIS_PORT",
	"MEMBER_ROLE_OWNER_ID",
	"MEMBER_ROLE_ADMIN_ID",
	"MEMBER_ROLE_MEMBER_ID",
	"SERVICE_NAME",
	"JWT_SECRET_KEY",
	"SENDGRID_API_KEY",
	"STRING_ENCRYPTION_KEY",
	"AUTH_EMAIL_ADDRESS",
}

var envMap map[string]string

func LoadEnv() error {
	godotenv.Load(".env")
	envMap = make(map[string]string)
	for _, key := range keys {
		val := os.Getenv(key)
		if val == "" {
			return common.StringError(errors.New("Missing env var: " + key))
		}
		envMap[key] = val
	}
	return nil
}

func Get(key string) (string, error) {
	res, ok := envMap[key]
	if !ok {
		return "", common.StringError(errors.New("No such env var: " + key))
	}
	return res, nil
}
