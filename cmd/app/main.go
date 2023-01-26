package main

import (
	"os"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/platform-admin-api/api"
	"github.com/String-xyz/platform-admin-api/pkg/store"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	// load .env file
	godotenv.Load(".env") // removed the err since in cloud this wont be loaded
	lg := zerolog.New(os.Stdout)
	if !common.IsLocalEnv() {
		tracer.Start()
		defer tracer.Stop()
	}

	port := os.Getenv("PORT")
	if port == "" {
		panic("no port!")
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	db := store.MustNewPG()

	// setup api
	api.Start(api.APIConfig{
		DB:     db,
		Redis:  store.NewRedis(),
		Port:   port,
		Logger: &lg,
	})
}
