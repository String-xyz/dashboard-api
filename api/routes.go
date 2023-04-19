package api

import (
	"github.com/String-xyz/platform-admin-api/api/handler"
	"github.com/String-xyz/platform-admin-api/api/middleware"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/labstack/echo/v4"
)

func organizationRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewOrganization(services.Organization)
	handler.RegisterRoutes(e.Group("/organizations"), middleware.JWT(services.Auth))
}

func platformRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewPlatform(services.Platform)
	handler.RegisterRoutes(e.Group("/platforms"), middleware.JWT(services.Auth))
}

func memberRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewMember(services)
	handler.RegisterRoutes(e.Group("/members"), middleware.JWT(services.Auth))
}

func loginRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewLogin(services)
	handler.RegisterRoutes(e.Group("/login"))
}

func inviteRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewInvite(services.Invite)
	handler.RegisterRoutes(e.Group("/invites"), middleware.JWT(services.Auth))
}

func apikeyRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewApikey(services.Apikey)
	handler.RegisterRoutes(e.Group("/apikeys"), middleware.JWT(services.Auth))
}

func contractRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewContract(services.Contract)
	handler.RegisterRoutes(e.Group("/contracts"), middleware.JWT(services.Auth))
}

func networkRoute(services service.Services, e *echo.Echo) {
	handler := handler.NewNetwork(services.Network)
	handler.RegisterRoutes(e.Group("/networks"))
}
