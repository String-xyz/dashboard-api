package main

import (
	"log"
	"os"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/platform-admin-api/api"
	"github.com/String-xyz/platform-admin-api/pkg/store"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func main() {
	// load .env file
	godotenv.Load(".env") // removed the err since in cloud this wont be loaded
	lg := zerolog.New(os.Stdout)
	if !common.IsLocalEnv() {
		setupTracer()
		defer profiler.Stop()
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

func setupTracer() {
	rules := []tracer.SamplingRule{tracer.RateRule(1)}
	tracer.Start(
		tracer.WithSamplingRules(rules),
		tracer.WithService("platform-api"),
		tracer.WithEnv(os.Getenv("ENV")),
	)

	err := profiler.Start(
		profiler.WithService("platform-api"),
		profiler.WithEnv(os.Getenv("ENV")),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
		))

	if err != nil {
		log.Fatal(err)
	}
}
