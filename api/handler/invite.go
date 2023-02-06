package handler

import (
	"net/http"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/httperror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type Invite interface {
	Send(e echo.Context) error
	Accept(e echo.Context) error
	List(e echo.Context) error
	Resend(e echo.Context) error
	Update(e echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type invite struct {
	service service.Invite
	Group   *echo.Group
}

func NewInvite(service service.Invite) Invite {
	return &invite{service: service}
}

func (i invite) Send(c echo.Context) error {
	body := model.RequestInviteSend{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "invite: send bind")
		return httperror.BadRequestError(c)
	}
	// TODO derive platform from JWT
	platform := model.Platform{}

	// TODO derive member_id from JWT
	m, err := i.service.Send(c.Request().Context(), body, platform)
	if err != nil {
		common.LogStringError(c, err, "invite: send")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusCreated, m)
}

func (i invite) Accept(c echo.Context) error {
	body := model.RequestInviteAcceptance{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "invite: accept bind")
		return httperror.BadRequestError(c)
	}
	id := c.Param("id")
	if id == "" {
		return httperror.BadRequestError(c)
	}
	body.Id = &id
	m, err := i.service.Accept(c.Request().Context(), body)
	if err != nil {
		common.LogStringError(c, err, "invite: accept")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusAccepted, m)
}

func (i invite) List(c echo.Context) error {
	status := c.QueryParam("status")
	if status == "" {
		return httperror.BadRequestError(c)
	}
	m, err := i.service.List(c.Request().Context(), status)
	if err != nil {
		common.LogStringError(c, err, "invite: list")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusAccepted, m)
}

func (i invite) Resend(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return httperror.BadRequestError(c)
	}
	m, err := i.service.Resend(c.Request().Context(), id)
	if err != nil {
		common.LogStringError(c, err, "invite: resend")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusAccepted, m)
}

func (i invite) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return httperror.BadRequestError(c)
	}
	body := model.RequestInviteUpdate{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "invite: update bind")
		return httperror.BadRequestError(c)
	}

	m, err := i.service.Update(c.Request().Context(), body, id)
	if err != nil {
		common.LogStringError(c, err, "invite: update")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusAccepted, m)
}

func (i invite) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("no group attached to the invite handler")
	}
	i.Group = g

	g.Use(ms...)
	g.POST("", i.Send, ms...)
	g.POST("/:id", i.Accept)
	g.GET("", i.List, ms...)
	g.POST("/:id/resend", i.Resend, ms...)
	g.PUT("/:id", i.Update, ms...)
}
