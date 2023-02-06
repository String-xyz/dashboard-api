package handler

import (
	"net/http"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/httperror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type Apikey interface {
	Create(e echo.Context) error
	GetAll(e echo.Context) error
	Get(e echo.Context) error
	Deactivate(e echo.Context) error
	Update(e echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type apikey struct {
	service service.Apikey
	Group   *echo.Group
}

func NewApikey(service service.Apikey) Apikey {
	return &apikey{service: service}
}

func (a apikey) Create(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	platformId := c.Get("platformId").(string)

	m, err := a.service.Create(c.Request().Context(), callerId, platformId)
	if err != nil {
		common.LogStringError(c, err, "apikey: create")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusCreated, m)
}

func (a apikey) GetAll(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	platformId := c.Get("platformId").(string)

	m, err := a.service.GetAll(c.Request().Context(), callerId, platformId)
	if err != nil {
		common.LogStringError(c, err, "apikey: get all")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (a apikey) Get(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	platformId := c.Get("platformId").(string)
	keyId := c.Param("id")
	if keyId == "" {
		return httperror.BadRequestError(c)
	}

	m, err := a.service.Get(c.Request().Context(), callerId, platformId, keyId)
	if err != nil {
		common.LogStringError(c, err, "apikey: get all")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (a apikey) Deactivate(c echo.Context) error {
	// TODO: Get platform ID from the JWT
	id := ""

	keyId := c.Param("id")
	if keyId == "" {
		return httperror.BadRequestError(c)
	}

	m, err := a.service.Deactivate(c.Request().Context(), id, keyId)
	if err != nil {
		common.LogStringError(c, err, "apikey: deactivate")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (a apikey) Update(c echo.Context) error {
	// TODO: Get platform ID from the JWT
	id := ""
	body := model.RequestApikeyUpdate{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "apikey: update bind")
		return httperror.BadRequestError(c)
	}

	m, err := a.service.Update(c.Request().Context(), body, id)
	if err != nil {
		common.LogStringError(c, err, "apikey: update")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (a apikey) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("no group attached to the apikey handler")
	}
	a.Group = g
	g.Use(ms...)
	g.POST("", a.Create, ms...)
	g.GET("", a.GetAll, ms...)
	g.GET("/:id", a.Get, ms...)
	g.PUT("/:id/deactivate", a.Deactivate, ms...)
	g.PUT("/:id", a.Update, ms...)
}
