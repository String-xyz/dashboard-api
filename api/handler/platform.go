package handler

import (
	"net/http"

	"github.com/String-xyz/dashboard-api/pkg/model"
	"github.com/String-xyz/dashboard-api/pkg/service"
	"github.com/String-xyz/go-lib/v2/common"
	httperror "github.com/String-xyz/go-lib/v2/httperror"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
	validator "github.com/String-xyz/go-lib/v2/validator"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
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
// @Param body body model.RequestPlatformCreate true "Platform Create Request"
// @Success 201 {object} model.Platform
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /platforms [post]
func (p platform) Create(c echo.Context) error {
	body := model.RequestPlatformCreate{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "platform: create bind")
		return httperror.BadRequest400(c)
	}

	if err := c.Validate(body); err != nil {
		common.LogStringError(c, err, "platform: create validate body")
		return httperror.InvalidPayload400(c, err)
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid organizationId")
	}

	m, err := p.service.Create(c.Request().Context(), body, organizationId)
	if err != nil {
		common.LogStringError(c, err, "platform: create")
		return httperror.Internal500(c)
	}
	return c.JSON(http.StatusCreated, m)
}

// @Summary Get platform
// @Tags Platform
// @Produce  json
// @Param id path string true "Platform ID"
// @Success 200 {object} model.Platform
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /platforms/{id} [get]
func (p platform) Get(c echo.Context) error {
	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid organizationId")
	}

	platformId := c.Param("id")
	if platformId == "" {
		return httperror.BadRequest400(c)
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
// @Success 200 {array} model.Platform
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /platforms [get]
func (p platform) GetAll(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid memberId")
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid organizationId")
	}

	m, err := p.service.GetAll(c.Request().Context(), callerId, organizationId)
	if err != nil && errors.Cause(err) != serror.NOT_FOUND {
		return DefaultErrorHandler(c, err, "platform: get all")
	}

	return c.JSON(http.StatusOK, m)
}

// @Summary Update platform
// @Tags Platform
// @Accept  json
// @Produce  json
// @Param id path string true "Platform ID"
// @Param body body model.RequestPlatformUpdate true "Platform Update Request"
// @Success 200 {object} model.Platform
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /platforms/{id} [patch]
func (p platform) Update(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid memberId")
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid organizationId")
	}

	platformId := c.Param("id")
	if platformId == "" {
		return httperror.BadRequest400(c)
	}

	body := model.RequestPlatformUpdate{}

	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "platform: update bind")
		return httperror.BadRequest400(c)
	}

	if body == (model.RequestPlatformUpdate{}) {
		return httperror.BadRequest400(c, "No fields to update")
	}

	// validate body
	if err := c.Validate(body); err != nil {
		common.LogStringError(c, err, "platform: update validate body")
		return httperror.InvalidPayload400(c, err)
	}

	m, err := p.service.Update(c.Request().Context(), body, platformId, callerId, organizationId)
	if err != nil {
		return DefaultErrorHandler(c, err, "platform: update")
	}
	return c.JSON(http.StatusOK, m)
}

// @Summary Deactivate platform
// @Tags Platform
// @Accept  json
// @Produce  json
// @Param id path string true "Platform ID"
// @Success 200 {object} model.Platform
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /platforms/{id}/deactivate [patch]
func (p platform) Deactivate(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid callerId")
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid organizationId")
	}

	id := c.Param("id")
	if !validator.IsUUID(id) {
		return httperror.BadRequest400(c, "invalid id")
	}

	m, err := p.service.Deactivate(c.Request().Context(), id, callerId, organizationId)
	if err != nil {
		return DefaultErrorHandler(c, err, "contract: deactivate")
	}

	return c.JSON(http.StatusOK, m)

}

// @Summary Reactivate platform
// @Tags Platform
// @Accept  json
// @Produce  json
// @Param id path string true "Platform ID"
// @Success 200 {object} model.Platform
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /platforms/{id}/reactivate [patch]
func (p platform) Reactivate(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid callerId")
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid organizationId")
	}

	id := c.Param("id")
	if !validator.IsUUID(id) {
		return httperror.BadRequest400(c, "invalid id")
	}

	m, err := p.service.Reactivate(c.Request().Context(), id, callerId, organizationId)
	if err != nil {
		return DefaultErrorHandler(c, err, "contract: reactivate")
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
	g.PATCH("/:id/deactivate", p.Deactivate, ms...)
	g.PATCH("/:id/reactivate", p.Reactivate, ms...)
}
