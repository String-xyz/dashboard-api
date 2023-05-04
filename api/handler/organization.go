package handler

import (
	"net/http"
	"strings"

	"github.com/String-xyz/go-lib/common"
	httperror "github.com/String-xyz/go-lib/httperror"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type Organization interface {
	Create(e echo.Context) error
	Get(e echo.Context) error
	Update(e echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type organization struct {
	service service.Organization
	Group   *echo.Group
}

func NewOrganization(service service.Organization) Organization {
	return &organization{service: service}
}

// @Summary Create organization
// @Tags Organization
// @Accept  json
// @Produce  json
// @Param body body model.RequestOrganizationCreate true "Organization Create Request"
// @Success 201 {object} model.Organization
// @Failure 400 {object} httperror.HTTPError
// @Failure 409 {object} httperror.HTTPError
// @Failure 500 {object} httperror.HTTPError
// @Router /organizations [post]
func (o organization) Create(c echo.Context) error {
	body := model.RequestOrganizationCreate{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "organization: create bind")
		return httperror.BadRequestError(c)
	}

	body.Email = strings.ToLower(body.Email)

	if err := c.Validate(body); err != nil {
		return httperror.InvalidPayloadError(c, err)
	}

	m, err := o.service.Create(c.Request().Context(), body)
	if err != nil {
		common.LogStringError(c, err, "organization: create")

		if serror.Is(err, serror.ALREADY_IN_USE) {
			return httperror.ConflictError(c, "email already in use")
		}

		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusCreated, m)
}

// @Summary Get organization
// @Tags Organization
// @Produce  json
// @Success 200 {object} model.Organization
// @Failure 500 {object} httperror.HTTPError
// @Router /organizations [get]
func (o organization) Get(c echo.Context) error {
	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid organizationId")
	}

	m, err := o.service.Get(c.Request().Context(), organizationId)
	if err != nil {
		return DefaultErrorHandler(c, err, "organization: get")
	}

	return c.JSON(http.StatusAccepted, m)
}

// @Summary Update organization
// @Tags Organization
// @Accept  json
// @Produce  json
// @Security JWTAuth
// @Security BearerAuth
// @Param body body model.RequestOrganizationUpdate true "Organization Update Request"
// @Success 200 {object} model.Organization
// @Failure 400 {object} httperror.HTTPError
// @Failure 500 {object} httperror.HTTPError
// @Router /organizations [patch]
func (o organization) Update(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid memberId")
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid organizationId")
	}

	body := model.RequestOrganizationUpdate{}

	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "organization: update bind")
		return httperror.BadRequestError(c)
	}

	if body == (model.RequestOrganizationUpdate{}) {
		return httperror.BadRequestError(c, "No fields to update")
	}

	// validate body
	if err := c.Validate(body); err != nil {
		return httperror.InvalidPayloadError(c, err)
	}

	m, err := o.service.Update(c.Request().Context(), body, organizationId, callerId)
	if err != nil {
		return DefaultErrorHandler(c, err, "organization: update")
	}
	return c.JSON(http.StatusOK, m)
}

func (o organization) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("no group attached to the organization handler")
	}
	o.Group = g
	g.POST("", o.Create)
	g.GET("", o.Get, ms...)
	g.PATCH("", o.Update, ms...)
}
