package handler

import (
	"net/http"

	"github.com/String-xyz/go-lib/common"
	httperror "github.com/String-xyz/go-lib/httperror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type Platform interface {
	Create(e echo.Context) error
	Get(e echo.Context) error
	GetAll(c echo.Context) error
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

// @Summary Create platform
// @Tags Platform
// @Accept  json
// @Produce  json
// @Security JWTAuth
// @Security BearerAuth
// @Param body body model.RequestPlatformCreate true "Platform Create Request"
// @Success 201 {object} model.Platform
// @Failure 400 {object} httperror.HTTPError
// @Failure 500 {object} httperror.HTTPError
// @Router /platforms [post]
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

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid organizationId")
	}

	m, err := p.service.Create(c.Request().Context(), body, organizationId)
	if err != nil {
		common.LogStringError(c, err, "platform: create")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusCreated, m)
}

// @Summary Get platform
// @Tags Platform
// @Produce  json
// @Security JWTAuth
// @Security BearerAuth
// @Param id path string true "Platform ID"
// @Success 200 {object} model.Platform
// @Failure 400 {object} httperror.HTTPError
// @Failure 500 {object} httperror.HTTPError
// @Router /platforms/{id} [get]
func (p platform) Get(c echo.Context) error {
	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid organizationId")
	}

	platformId := c.Param("id")
	if platformId == "" {
		return httperror.BadRequestError(c)
	}

	m, err := p.service.Get(c.Request().Context(), platformId, organizationId)
	if err != nil {
		return DefaultErrorHandler(c, err, "platform: get")
	}

	return c.JSON(http.StatusAccepted, m)
}

// @Summary Get all platforms
// @Tags Platform
// @Produce  json
// @Security JWTAuth
// @Security BearerAuth
// @Success 200 {array} model.Platform
// @Failure 500 {object} httperror.HTTPError
// @Router /platforms [get]
func (p platform) GetAll(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid memberId")
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid organizationId")
	}

	m, err := p.service.GetAll(c.Request().Context(), callerId, organizationId)
	if err != nil {
		return DefaultErrorHandler(c, err, "apikey: get all")
	}

	return c.JSON(http.StatusOK, m)
}

// @Summary Update platform
// @Tags Platform
// @Accept  json
// @Produce  json
// @Security JWTAuth
// @Security BearerAuth
// @Param id path string true "Platform ID"
// @Param body body model.RequestPlatformUpdate true "Platform Update Request"
// @Success 200 {object} model.Platform
// @Failure 400 {object} httperror.HTTPError
// @Failure 500 {object} httperror.HTTPError
// @Router /platforms/{id} [patch]
func (p platform) Update(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid memberId")
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid organizationId")
	}

	platformId := c.Param("id")
	if platformId == "" {
		return httperror.BadRequestError(c)
	}

	body := model.RequestPlatformUpdate{}

	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "platform: update bind")
		return httperror.BadRequestError(c)
	}

	if body == (model.RequestPlatformUpdate{}) {
		return httperror.BadRequestError(c, "No fields to update")
	}

	// validate body
	if err := c.Validate(body); err != nil {
		return httperror.InvalidPayloadError(c, err)
	}

	m, err := p.service.Update(c.Request().Context(), body, platformId, callerId, organizationId)
	if err != nil {
		return DefaultErrorHandler(c, err, "platform: update")
	}
	return c.JSON(http.StatusOK, m)
}

func (p platform) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("no group attached to the platform handler")
	}
	p.Group = g
	g.POST("", p.Create, ms...)
	g.GET("", p.GetAll, ms...)
	g.GET("/:id", p.Get, ms...)
	g.PATCH("/:id", p.Update, ms...)
}
