package handler

import (
	"net/http"
	"strings"

	"github.com/String-xyz/go-lib/common"
	httperror "github.com/String-xyz/go-lib/httperror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Platform interface {
	Create(e echo.Context) error
	Get(e echo.Context) error
	Update(e echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type platform struct {
	service service.Platform
	Group   *echo.Group
}

func NewPlatform(service service.Platform) Platform {
	return &platform{service: service}
}

func (p platform) Create(c echo.Context) error {
	body := model.RequestPlatformCreate{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "platform: create bind")
		return httperror.BadRequestError(c)
	}

	if err := c.Validate(body); err != nil {
		return httperror.InvalidPayloadError(c, err)
	}

	m, err := p.service.Create(c.Request().Context(), body)
	if err != nil {
		common.LogStringError(c, err, "platform: create")

		errMessage := errors.Cause(err).Error()
		if strings.Contains(errMessage, "already in use") {
			return httperror.ConflictError(c, errMessage)
		}

		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusCreated, m)
}

func (p platform) Get(c echo.Context) error {
	platformId := c.Get("platformId").(string)
	m, err := p.service.Get(c.Request().Context(), platformId)
	if err != nil {
		common.LogStringError(c, err, "platform: get")

		if errors.Cause(err).Error() == "sql: no rows in result set" || strings.Contains(errors.Cause(err).Error(), "not found") {
			return httperror.NotFoundError(c)
		}

		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusAccepted, m)
}

func (p platform) Update(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	platformId := c.Get("platformId").(string)
	body := model.RequestPlatformUpdate{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "platform: update bind")
		return httperror.BadRequestError(c)
	}

	// validate body
	if err := c.Validate(body); err != nil {
		return httperror.InvalidPayloadError(c, err)
	}

	m, err := p.service.Update(c.Request().Context(), body, platformId, callerId)
	if err != nil {
		common.LogStringError(c, err, "platform: update")

		if errors.Cause(err).Error() == "sql: no rows in result set" || strings.Contains(errors.Cause(err).Error(), "not found") {
			return httperror.NotFoundError(c)
		}

		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (p platform) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("no group attached to the platform handler")
	}
	p.Group = g
	g.POST("", p.Create)
	g.GET("", p.Get, ms...)
	g.PUT("", p.Update, ms...)
}
