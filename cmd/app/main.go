package main

import (
	"os"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/platform-admin-api/api"
	"github.com/String-xyz/platform-admin-api/config"
	"github.com/String-xyz/platform-admin-api/pkg/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	// load env vars
	err := config.LoadEnv()
	if err != nil {
		panic(err)
	}
	lg := zerolog.New(os.Stdout)
	if !common.IsLocalEnv() {
		tracer.Start()
		defer tracer.Stop()
	}

	port := config.Var.PORT

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
