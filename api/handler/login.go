package handler

import (
	"net/http"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/httperror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type Login interface {
	Login(e echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type login struct {
	service service.Login
	Group   *echo.Group
}

func NewLogin(service service.Login) Login {
	return &login{service: service}
}

func (l login) Login(c echo.Context) error {
	body := model.RequestLogin{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "login: login bind")
		return httperror.BadRequestError(c)
	}

	m, err := l.service.Login(c.Request().Context(), body)
	if err != nil {
		common.LogStringError(c, err, "login: login")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusAccepted, m)
}

func (l login) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("no group attached to the login handler")
	}
	l.Group = g
	// g.Use(ms...)
	g.POST("", l.Login)
}
