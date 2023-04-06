package handler

import (
	"net/http"

	"github.com/String-xyz/go-lib/common"
	httperror "github.com/String-xyz/go-lib/httperror"
	validator "github.com/String-xyz/go-lib/validator"
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
	keyType := c.QueryParam("type")

	m, err := a.service.Create(c.Request().Context(), callerId, platformId, keyType)
	if err != nil {
		return DefaultErrorHandler(c, err, "apikey: create")
	}

	return c.JSON(http.StatusCreated, m)
}

func (a apikey) GetAll(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	platformId := c.Get("platformId").(string)

	m, err := a.service.GetAll(c.Request().Context(), callerId, platformId)
	if err != nil {
		return DefaultErrorHandler(c, err, "apikey: get all")
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
		return DefaultErrorHandler(c, err, "apikey: get")
	}

	return c.JSON(http.StatusOK, m)
}

func (a apikey) Deactivate(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	platformId := c.Get("platformId").(string)
	keyId := c.Param("id")

	if !validator.IsUUID(keyId, platformId, callerId) {
		return httperror.BadRequestError(c, "invalid id")
	}

	m, err := a.service.Deactivate(c.Request().Context(), callerId, platformId, keyId)
	if err != nil {
		return DefaultErrorHandler(c, err, "apikey: deactivate")
	}

	return c.JSON(http.StatusOK, m)
}

func (a apikey) Update(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	platformId := c.Get("platformId").(string)
	keyId := c.Param("id")

	if platformId == "" || callerId == "" {
		return httperror.BadRequestError(c)
	}

	if !validator.IsUUID(keyId) {
		return httperror.BadRequestError(c, "invalid id")
	}

	body := model.RequestApikeyUpdate{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "apikey: update bind")
		return httperror.BadRequestError(c)
	}

	m, err := a.service.Update(c.Request().Context(), body, callerId, platformId, keyId)
	if err != nil {
		return DefaultErrorHandler(c, err, "apikey: update")
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
