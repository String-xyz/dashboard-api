package handler

import (
	"net/http"
	"strconv"

	"github.com/String-xyz/dashboard-api/pkg/service"
	"github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/httperror"
	"github.com/labstack/echo/v4"
)

type Network interface {
	GetAll(e echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type network struct {
	service service.Network
	Group   *echo.Group
}

func NewNetwork(service service.Network) Network {
	return &network{service: service}
}

// @Summary Get all networks
// @Tags Network
// @Produce  json
// @Success 200 {array} model.NetworkData
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /networks [get]
func (a network) GetAll(c echo.Context) error {
	var err error
	var limit int
	var offset int
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			common.LogStringError(c, err, "network get all: invalid limit")
			return httperror.BadRequest400(c, "invalid limit")
		}
	}
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			common.LogStringError(c, err, "network get all: invalid offset")
			return httperror.BadRequest400(c, "invalid offset")
		}
	}

	m, err := a.service.GetAll(c.Request().Context(), limit, offset)
	if err != nil {
		common.LogStringError(c, err, "network: get all")
		return DefaultErrorHandler(c, err, "network: get all")
	}
	return c.JSON(http.StatusOK, m)
}

func (a network) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("no group attached to the apikey handler")
	}
	a.Group = g
	g.Use(ms...)
	g.GET("", a.GetAll, ms...)
}
