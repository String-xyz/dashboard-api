package api

import (
	"net/http"

	"github.com/String-xyz/go-lib/database"
	"github.com/String-xyz/go-lib/middleware"
	"github.com/String-xyz/go-lib/validator"
	_ "github.com/String-xyz/platform-admin-api/docs"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	echoSwagger "github.com/swaggo/echo-swagger"
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
	e.Use(middleware.Tracer("platform-api"))
	e.Use(middleware.CORS())
	e.Use(middleware.RequestId())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger(logger))
	e.Use(middleware.LogRequest())
}

// @title Platform Management API
// @version 1.0
// @description Platform Management API for managing Organizations, Platforms, API keys, etc.

// @contact.name Platform Management API Support
// @contact.url http://string.xyz
// @contact.email support@stringxyz.com

// @host api.sandbox.string-api.xyz
// @BasePath /
func Start(config APIConfig) {
	e := echo.New()
	e.Validator = validator.New()
	e.GET("/heartbeat", heartbeat)

	// setup swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	baseMiddleware(config.Logger, e)

	// Routes
	repos := newRepos(config)
	redis := newRedis()
	service := newServices(config, repos, redis)

	organizationRoute(service, e)
	platformRoute(service, e)
	memberRoute(service, e)
	loginRoute(service, e)
	inviteRoute(service, e)
	apikeyRoute(service, e)
	contractRoute(service, e)
	networkRoute(service, e)

	e.Logger.Fatal(e.Start(":" + config.Port))
}
