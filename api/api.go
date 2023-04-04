package api

import (
	"net/http"

	"github.com/String-xyz/go-lib/database"
	"github.com/String-xyz/go-lib/middleware"
	"github.com/String-xyz/go-lib/validator"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type APIConfig struct {
	DB     database.Queryable
	Redis  database.RedisStore
	Logger *zerolog.Logger
	Port   string
}

func heartbeat(c echo.Context) error {
	return c.JSON(http.StatusOK, "alive")
}

func baseMiddleware(logger *zerolog.Logger, e *echo.Echo) {
	e.Use(middleware.Tracer())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger(logger))
	e.Use(middleware.LogRequest())
}

func Start(config APIConfig) {
	e := echo.New()
	e.Validator = validator.New()
	e.GET("/heartbeat", heartbeat)

	baseMiddleware(config.Logger, e)

	// Routes
	repos := newRepos(config)
	redis := newRedis()
	service := newServices(config, repos, redis)

	platformRoute(service, e)
	memberRoute(service, e)
	loginRoute(service, e)
	inviteRoute(service, e)
	apikeyRoute(service, e)
	contractRoute(service, e)
	networkRoute(service, e)

	e.Logger.Fatal(e.Start(":" + config.Port))
}
