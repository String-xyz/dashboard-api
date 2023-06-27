package config

type vars struct {
	ENV                   string `required:"true"`
	PORT                  string `required:"true"`
	BASE_DASHBOARD_URL    string `required:"true"`
	DB_NAME               string `required:"true"`
	DB_USERNAME           string `required:"true"`
	DB_PASSWORD           string `required:"true"`
	DB_HOST               string `required:"true"`
	DB_PORT               string `required:"true"`
	REDIS_PASSWORD        string `required:"true"`
	REDIS_HOST            string `required:"true"`
	REDIS_PORT            string `required:"true"`
	MEMBER_ROLE_OWNER_ID  string `required:"true"`
	MEMBER_ROLE_ADMIN_ID  string `required:"true"`
	MEMBER_ROLE_MEMBER_ID string `required:"true"`
	JWT_SECRET_KEY        string `required:"true"`
	SENDGRID_API_KEY      string `required:"true"`
	STRING_ENCRYPTION_KEY string `required:"true"`
	AUTH_EMAIL_ADDRESS    string `required:"true"`
}

var Var vars
