package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/String-xyz/dashboard-api/pkg/model"
	"github.com/String-xyz/dashboard-api/pkg/service"
	"github.com/String-xyz/go-lib/v2/common"
	httperror "github.com/String-xyz/go-lib/v2/httperror"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
	validator "github.com/String-xyz/go-lib/v2/validator"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Apikey interface {
	Create(e echo.Context) error
	GetAll(e echo.Context) error
	Get(e echo.Context) error
	Delete(e echo.Context) error
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

// @Summary Create API key
// @Description Create an API key
// @Tags Apikey
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param type query string false "Type of API key (public or secret)"
// @Param platformId query string false "Platform ID"
// @Success 201 {object} model.Apikey
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /apikeys [post]
func (a apikey) Create(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid memberId")
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid organizationId")
	}

	keyType := c.QueryParam("type")
	if keyType == "" {
		keyType = "public"
	}

	keyType = strings.ToLower(keyType)

	if keyType != "public" && keyType != "secret" {
		return httperror.BadRequest400(c, "invalid key type")
	}

	platformId := c.QueryParam("platformId")
	if keyType == "public" && platformId == "" {
		return httperror.BadRequest400(c)
	}

	m, err := a.service.Create(c.Request().Context(), keyType, callerId, platformId, organizationId)
	if err != nil {
		return DefaultErrorHandler(c, err, "apikey: create")
	}

	return c.JSON(http.StatusCreated, m)
}

// @Summary Get all API keys
// @Description Get all API keys
// @Tags Apikey
// @Accept json
// @Produce json
// @Param platformId query string false "Platform ID"
// @Param limit query int false "Limit the number of returned results"
// @Param offset query int false "Offset the starting point of returned results"
// @Success 200 {array} model.Apikey
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /apikeys [get]
func (a apikey) GetAll(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid memberId")
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid organizationId")
	}

	platformId := c.QueryParam("platformId")

	var err error
	var limit int
	var offset int
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return httperror.BadRequest400(c, "invalid limit")
		}
	}
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return httperror.BadRequest400(c, "invalid offset")
		}
	}

	m, err := a.service.GetAll(c.Request().Context(), callerId, platformId, organizationId, limit, offset)
	if err != nil && errors.Cause(err) != serror.NOT_FOUND {
		return DefaultErrorHandler(c, err, "apikey: get all")
	}

	return c.JSON(http.StatusOK, m)
}

// @Summary Get API key
// @Description Get API key by ID
// @Tags Apikey
// @Accept json
// @Produce json
// @Param id path string true "API Key ID"
// @Success 200 {object} model.Apikey
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /apikeys/{id} [get]
func (a apikey) Get(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid memberId")
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid organizationId")
	}

	keyId := c.Param("id")
	if keyId == "" {
		return httperror.BadRequest400(c)
	}

	m, err := a.service.Get(c.Request().Context(), keyId, callerId, organizationId)
	if err != nil {
		return DefaultErrorHandler(c, err, "apikey: get")
	}

	return c.JSON(http.StatusOK, m)
}

// @Summary Delete API key
// @Description Delete API key by ID
// @Tags Apikey
// @Accept json
// @Produce json
// @Param id path string true "API Key ID"
// @Success 204 "No Content"
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /apikeys/{id} [delete]
func (a apikey) Delete(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid callerId")
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid organizationId")
	}

	keyId := c.Param("id")

	if !validator.IsUUID(keyId) {
		return httperror.BadRequest400(c, "invalid id")
	}

	err := a.service.Delete(c.Request().Context(), keyId, callerId, organizationId)
	if err != nil {
		return DefaultErrorHandler(c, err, "apikey: deactivate")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "success"})
}

// @Summary Update API key
// @Description Update API key by ID
// @Tags Apikey
// @Accept json
// @Produce json
// @Param id path string true "API Key ID"
// @Param body body model.RequestApikeyUpdate true "API key Update Request"
// @Success 200 {object} model.Apikey
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /apikeys/{id} [put]
func (a apikey) Update(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid organizationId")
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid organizationId")
	}

	keyId := c.Param("id")

	if !validator.IsUUID(keyId) {
		return httperror.BadRequest400(c, "invalid id")
	}

	body := model.RequestApikeyUpdate{}

	if err := c.Bind(&body); err != nil {
		common.LogStringError(c, err, "apikey: update bind")
		return httperror.BadRequest400(c)
	}

	if err := c.Validate(&body); err != nil {
		common.LogStringError(c, err, "apikey: update validate")
		return httperror.BadRequest400(c)
	}

	m, err := a.service.Update(c.Request().Context(), keyId, body, callerId, organizationId)
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
	g.PATCH("/:id", a.Update, ms...)
	g.DELETE("/:id", a.Delete, ms...)
}
