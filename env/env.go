package env

import (
	"errors"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type vars struct {
	ENV                   string
	PORT                  string
	BASE_DASHBOARD_URL    string
	DB_NAME               string
	DB_USERNAME           string
	DB_PASSWORD           string
	DB_HOST               string
	DB_PORT               string
	REDIS_PASSWORD        string
	REDIS_HOST            string
	REDIS_PORT            string
	MEMBER_ROLE_OWNER_ID  string
	MEMBER_ROLE_ADMIN_ID  string
	MEMBER_ROLE_MEMBER_ID string
	SERVICE_NAME          string
	JWT_SECRET_KEY        string
	SENDGRID_API_KEY      string
	STRING_ENCRYPTION_KEY string
	AUTH_EMAIL_ADDRESS    string
}

var Var vars

func LoadEnv() error {
	godotenv.Load(".env")
	stype := reflect.ValueOf(Var).Elem()
	for i := 0; i < stype.NumField(); i++ {
		field := stype.Field(i)
		key := stype.Type().Field(i).Name
		value := os.Getenv(key)
		if value == "" {
			return errors.New("Missing environment variable: " + key)
		}
		field.SetString(value)
	}
	return nil
}
