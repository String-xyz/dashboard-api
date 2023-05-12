package main

import (
	"log"
	"os"

	"github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/platform-admin-api/api"
	"github.com/String-xyz/platform-admin-api/config"
	"github.com/String-xyz/platform-admin-api/pkg/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func main() {
	// load env vars
	err := config.LoadEnv()
	if err != nil {
		panic(err)
	}
	lg := zerolog.New(os.Stdout)
	if !common.IsLocalEnv() {
		setupTracer()
		defer profiler.Stop()
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

func setupTracer() {
	rules := []tracer.SamplingRule{tracer.RateRule(1)}
	tracer.Start(
		tracer.WithSamplingRules(rules),
		tracer.WithService("platform-api"),
		tracer.WithEnv(config.Var.ENV),
	)

	err := profiler.Start(
		profiler.WithService("platform-api"),
		profiler.WithEnv(config.Var.ENV),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
		))

	if err != nil {
		log.Fatal(err)
	}
}
