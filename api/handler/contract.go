package handler

import (
	"net/http"

	"github.com/String-xyz/go-lib/common"
	httperror "github.com/String-xyz/go-lib/httperror"
	validator "github.com/String-xyz/go-lib/validator"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Contract interface {
	Create(e echo.Context) error
	GetAll(e echo.Context) error
	Get(e echo.Context) error
	Deactivate(e echo.Context) error
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
	callerId := c.Get("memberId").(string)
	platformId := c.Get("platformId").(string)

	body := model.RequestContractUpdate{}
	if err := c.Bind(&body); err != nil {
		return httperror.BadRequestError(c)
	}

	m, err := a.service.Create(c.Request().Context(), body, callerId, platformId)
	if err != nil {
		if errors.Cause(err).Error() == "contract already exists" {
			return httperror.ConflictError(c, err.Error())
		}
		common.LogStringError(c, err, "contract: create")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusCreated, m)
}

func (a contract) GetAll(c echo.Context) error {
	platformId := c.Get("platformId").(string)

	m, err := a.service.GetAll(c.Request().Context(), platformId)
	if err != nil {
		common.LogStringError(c, err, "contract: get all")
		return DefaultErrorHandler(c, err)
	}
	return c.JSON(http.StatusOK, m)
}

func (a contract) Get(c echo.Context) error {
	platformId := c.Get("platformId").(string)
	contractId := c.Param("id")
	if contractId == "" {
		return httperror.BadRequestError(c)
	}

	m, err := a.service.Get(c.Request().Context(), platformId, contractId)
	if err != nil {
		common.LogStringError(c, err, "contract: get")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (a contract) Deactivate(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	platformId := c.Get("platformId").(string)
	contractId := c.Param("id")
	if !validator.IsUUID(contractId, platformId, callerId) {
		return httperror.BadRequestError(c, "invalid id")
	}

	m, err := a.service.Deactivate(c.Request().Context(), callerId, platformId, contractId)
	if err != nil {
		common.LogStringError(c, err, "contract: deactivate")
		return DefaultErrorHandler(c, err)
	}
	return c.JSON(http.StatusOK, m)
}

func (a contract) Update(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	platformId := c.Get("platformId").(string)
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

	m, err := a.service.Update(c.Request().Context(), body, callerId, platformId, contractId)
	if err != nil {
		common.LogStringError(c, err, "contract: update")
		return DefaultErrorHandler(c, err)
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
	g.PUT("/:id/deactivate", a.Deactivate, ms...)
	g.PUT("/:id", a.Update, ms...)
}
