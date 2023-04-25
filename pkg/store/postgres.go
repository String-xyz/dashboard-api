package store

import (
	"fmt"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/platform-admin-api/env"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	sqlxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"
)

var pgDB *sqlx.DB
var DBDriver = "postgres"

func strConnection() string {
	username, _ := env.Get("DB_USERNAME")
	password, _ := env.Get("DB_PASSWORD")
	name, _ := env.Get("DB_NAME")
	host, _ := env.Get("DB_HOST")
	port, _ := env.Get("DB_PORT")

	var (
		DBUser     = username
		DBPassword = password
		DBName     = name
		DBHost     = host
		DBPort     = port
	)

	var SSLMode string

	if common.IsLocalEnv() {
		SSLMode = "disable"
	} else {
		SSLMode = "require"
	}

	str := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		DBHost,
		DBPort,
		DBUser,
		DBName,
		DBPassword,
		SSLMode,
	)

	return str
}

func MustNewPG() *sqlx.DB {
	if pgDB != nil {
		return pgDB
	}
	sqltrace.Register(DBDriver, &pq.Driver{}, sqltrace.WithServiceName("platform-admin-api"))
	connection, err := sqlxtrace.Open(DBDriver, strConnection())
	if err != nil {
		panic(err)
	}

	if err := connection.Ping(); err != nil {
		panic(err)
	}

	pgDB = connection
	return pgDB
}
