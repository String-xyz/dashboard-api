package handler

import (
	"net/http"

	"github.com/String-xyz/go-lib/common"
	httperror "github.com/String-xyz/go-lib/httperror"
	serror "github.com/String-xyz/go-lib/stringerror"
	validator "github.com/String-xyz/go-lib/validator"

	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type Contract interface {
	Create(e echo.Context) error
	GetAll(e echo.Context) error
	Get(e echo.Context) error
	Deactivate(e echo.Context) error
	Reactivate(e echo.Context) error
	Update(e echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type contract struct {
	service service.Contract
	Group   *echo.Group
}

func NewContract(service service.Contract) Contract {
	return &contract{service: service}
}

// Create a new contract
// @Summary Create a new contract
// @Tags Contract
// @Accept  json
// @Produce  json
// @Param RequestContractCreate body model.RequestContractCreate true "Contract Create Request"
// @Success 201 {object} model.Contract
// @Failure 400 {object} httperror.HTTPError
// @Failure 409 {object} httperror.HTTPError
// @Failure 500 {object} httperror.HTTPError
// @Router /contracts [post]
func (a contract) Create(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid memberId")
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid organizationId")
	}

	body := model.RequestContractCreate{}
	if err := c.Bind(&body); err != nil {
		return httperror.BadRequestError(c)
	}

	ok = validator.IsUUID(body.PlatformId)
	if !ok {
		return httperror.BadRequestError(c, "missing or invalid platformId")
	}

	SanitizeChecksums(&body.Address)

	m, err := a.service.Create(c.Request().Context(), body, callerId, organizationId)
	if err != nil {
		common.LogStringError(c, err, "contract: create")

		if serror.Is(err, serror.ALREADY_IN_USE) {
			return httperror.ConflictError(c, err.Error())
		}
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusCreated, m)
}

// Get all contracts
// @Summary Get all contracts
// @Tags Contract
// @Accept  json
// @Produce  json
// @Param platformId query string false "Platform ID"
// @Success 200 {array} model.Contract
// @Failure 500 {object} httperror.HTTPError
// @Router /contracts [get]
func (a contract) GetAll(c echo.Context) error {
	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid organizationId")
	}

	platformId := c.QueryParam("platformId")

	m, err := a.service.GetAll(c.Request().Context(), platformId, organizationId)
	if err != nil {
		return DefaultErrorHandler(c, err, "contract: get all")
	}
	return c.JSON(http.StatusOK, m)
}

// Get a specific contract
// @Summary Get a specific contract
// @Tags Contract
// @Accept  json
// @Produce  json
// @Param id path string true "Contract ID"
// @Success 200 {object} model.Contract
// @Failure 400 {object} httperror.HTTPError
// @Failure 500 {object} httperror.HTTPError
// @Router /contracts/{id} [get]
func (a contract) Get(c echo.Context) error {
	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid organizationId")
	}

	contractId := c.Param("id")
	if contractId == "" {
		return httperror.BadRequestError(c)
	}

	m, err := a.service.Get(c.Request().Context(), contractId, organizationId)
	if err != nil {
		return DefaultErrorHandler(c, err, "contract: get")
	}
	return c.JSON(http.StatusOK, m)
}

// Deactivate a contract
// @Summary Deactivate a contract
// @Tags Contract
// @Accept  json
// @Produce  json
// @Param id path string true "Contract ID"
// @Success 200 {object} model.Contract
// @Failure 400 {object} httperror.HTTPError
// @Failure 500 {object} httperror.HTTPError
// @Router /contracts/{id}/deactivate [patch]
func (a contract) Deactivate(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid callerId")
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid organizationId")
	}

	contractId := c.Param("id")
	if !validator.IsUUID(contractId, organizationId, callerId) {
		return httperror.BadRequestError(c, "invalid id")
	}

	m, err := a.service.Deactivate(c.Request().Context(), contractId, callerId, organizationId)
	if err != nil {
		return DefaultErrorHandler(c, err, "contract: deactivate")
	}
	return c.JSON(http.StatusOK, m)
}

// Reactivate a contract
// @Summary Reactivate a contract
// @Tags Contract
// @Accept  json
// @Produce  json
// @Param id path string true "Contract ID"
// @Success 200 {object} model.Contract
// @Failure 400 {object} httperror.HTTPError
// @Failure 500 {object} httperror.HTTPError
// @Router /contracts/{id}/reactivate [patch]
func (a contract) Reactivate(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid callerId")
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid organizationId")
	}

	contractId := c.Param("id")
	if !validator.IsUUID(contractId, organizationId, callerId) {
		return httperror.BadRequestError(c, "invalid id")
	}

	m, err := a.service.Reactivate(c.Request().Context(), contractId, callerId, organizationId)
	if err != nil {
		return DefaultErrorHandler(c, err, "contract: reactivate")
	}
	return c.JSON(http.StatusOK, m)
}

// Update a contract
// @Summary Update a contract
// @Tags Contract
// @Accept  json
// @Produce  json
// @Param id path string true "Contract ID"
// @Param RequestContractUpdate body model.RequestContractUpdate true "Contract Update Request"
// @Success 200 {object} model.Contract
// @Failure 400 {object} httperror.HTTPError
// @Failure 500 {object} httperror.HTTPError
// @Router /contracts/{id} [patch]
func (a contract) Update(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid callerId")
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid organizationId")
	}

	contractId := c.Param("id")
	if !validator.IsUUID(contractId) {
		return httperror.BadRequestError(c, "invalid id")
	}

	body := model.RequestContractUpdate{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "contract: update bind")
		return httperror.BadRequestError(c)
	}
	if body.Address != nil {
		SanitizeChecksums(body.Address)
	}

	m, err := a.service.Update(c.Request().Context(), body, contractId, callerId, organizationId)
	if err != nil {
		return DefaultErrorHandler(c, err, "contract: update")
	}

	return c.JSON(http.StatusOK, m)
}

func (a contract) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("no group attached to the apikey handler")
	}
	a.Group = g
	g.Use(ms...)
	g.POST("", a.Create, ms...)
	g.GET("", a.GetAll, ms...)
	g.GET("/:id", a.Get, ms...)
	g.PATCH("/:id/deactivate", a.Deactivate, ms...)
	g.PATCH("/:id/reactivate", a.Reactivate, ms...)
	g.PATCH("/:id", a.Update, ms...)
}
