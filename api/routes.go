package api

import (
	"github.com/String-xyz/platform-admin-api/api/handler"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/labstack/echo/v4"
)

func platformRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewPlatform(services.Platform)
	handler.RegisterRoutes(e.Group("/platforms"))
}
