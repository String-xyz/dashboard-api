package handler

import (
	"net/http"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/httperror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type Platform interface {
	Create(e echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type platform struct {
	service service.Platform
}

func NewPlatform(service service.Platform) Platform {
	return &platform{service: service}
}

func (p platform) Create(c echo.Context) error {
	body := model.Platform{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "platform: create bind")
		return httperror.BadRequestError(c)
	}

	m, err := p.service.Create()
	if err != nil {
		common.LogStringError(c, err, "platform: create")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusCreated, m)
}

func (p platform) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("no group attached to the platform handler")
	}
	g.Use(ms...)
	g.POST("", p.Create)
}
