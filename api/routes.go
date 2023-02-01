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

func memberRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewMember(services.Member)
	handler.RegisterRoutes(e.Group("/members"))
}

func loginRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewLogin(services.Login)
	handler.RegisterRoutes(e.Group("/login"))
}

func inviteRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewInvite(services.Invite)
	handler.RegisterRoutes(e.Group("/invites"))
}

func apikeyRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewApikey(services.Apikey)
	handler.RegisterRoutes(e.Group("/apikeys"))
}
