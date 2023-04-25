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

func (a contract) Create(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid memberId")
	}

	platformId, ok := c.Get("platformId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid platformId")
	}

	body := model.RequestContractCreate{}
	if err := c.Bind(&body); err != nil {
		return httperror.BadRequestError(c)
	}

	SanitizeChecksums(&body.Address)
	err := common.SanitizeIdInput(&body)
	if err != nil {
		return httperror.BadRequestError(c)
	}

	m, err := a.service.Create(c.Request().Context(), body, callerId, platformId)
	if err != nil {
		common.LogStringError(c, err, "contract: create")

		if serror.Is(err, serror.ALREADY_IN_USE) {
			return httperror.ConflictError(c, err.Error())
		}
		return httperror.InternalError(c)
	}
	err = common.SanitizeIdOutput(&m)
	if err != nil {
		return httperror.InternalError(c, "failed to sanitize id output")
	}
	return c.JSON(http.StatusCreated, m)
}

func (a contract) GetAll(c echo.Context) error {
	platformId, ok := c.Get("platformId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid platformId")
	}

	m, err := a.service.GetAll(c.Request().Context(), platformId)
	if err != nil {
		return DefaultErrorHandler(c, err, "contract: get all")
	}

	for i := range m {
		err = common.SanitizeIdOutput(&m[i])
		if err != nil {
			return httperror.InternalError(c, "failed to sanitize id output")
		}
	}
	return c.JSON(http.StatusOK, m)
}

func (a contract) Get(c echo.Context) error {
	platformId, ok := c.Get("platformId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid platformId")
	}

	contractId := c.Param("id")
	if contractId == "" {
		return httperror.BadRequestError(c)
	}

	m, err := a.service.Get(c.Request().Context(), platformId, contractId)
	if err != nil {
		return DefaultErrorHandler(c, err, "contract: get")
	}
	err = common.SanitizeIdOutput(&m)
	if err != nil {
		return httperror.InternalError(c, "failed to sanitize id output")
	}
	return c.JSON(http.StatusOK, m)
}

func (a contract) Deactivate(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid callerId")
	}

	platformId, ok := c.Get("platformId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid platformId")
	}

	contractId := c.Param("id")
	if !validator.IsUUID(contractId, platformId, callerId) {
		return httperror.BadRequestError(c, "invalid id")
	}

	m, err := a.service.Deactivate(c.Request().Context(), callerId, platformId, contractId)
	if err != nil {
		return DefaultErrorHandler(c, err, "contract: deactivate")
	}

	err = common.SanitizeIdOutput(&m)
	if err != nil {
		return httperror.InternalError(c, "failed to sanitize id output")
	}
	return c.JSON(http.StatusOK, m)
}

func (a contract) Reactivate(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid callerId")
	}

	platformId, ok := c.Get("platformId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid platformId")
	}

	contractId := c.Param("id")
	if !validator.IsUUID(contractId, platformId, callerId) {
		return httperror.BadRequestError(c, "invalid id")
	}

	m, err := a.service.Reactivate(c.Request().Context(), callerId, platformId, contractId)
	if err != nil {
		return DefaultErrorHandler(c, err, "contract: reactivate")
	}

	err = common.SanitizeIdOutput(&m)
	if err != nil {
		return httperror.InternalError(c, "failed to sanitize id output")
	}
	return c.JSON(http.StatusOK, m)
}

func (a contract) Update(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid callerId")
	}

	platformId, ok := c.Get("platformId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid platformId")
	}

	contractId := c.Param("id")
	if platformId == "" || callerId == "" {
		return httperror.BadRequestError(c)
	}

	if !validator.IsUUID(contractId) {
		return httperror.BadRequestError(c, "invalid id")
	}

	body := model.RequestContractUpdate{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "contract: update bind")
		return httperror.BadRequestError(c)
	}

	SanitizeChecksums(body.Address)

	err = common.SanitizeIdInput(&body)
	if err != nil {
		return httperror.BadRequestError(c)
	}

	m, err := a.service.Update(c.Request().Context(), body, callerId, platformId, contractId)
	if err != nil {
		return DefaultErrorHandler(c, err, "contract: update")
	}

	err = common.SanitizeIdOutput(&m)
	if err != nil {
		return httperror.InternalError(c, "failed to sanitize id output")
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
