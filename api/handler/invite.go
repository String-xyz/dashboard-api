package handler

import (
	"net/http"
	"strings"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Invite interface {
	Send(e echo.Context) error
	Accept(e echo.Context) error
	List(e echo.Context) error
	Resend(e echo.Context) error
	Update(e echo.Context) error
	Deactivate(e echo.Context) error
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
	callerId := c.Get("memberId").(string)
	platformId := c.Get("platformId").(string)

	body := model.RequestInviteSend{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "invite: send bind")
		return BadRequestError(c)
	}

	// validate email
	if err := c.Validate(body); err != nil {
		common.LogStringError(c, err, "invite: send validate")
		return InvalidPayloadError(c, err)
	}

	m, err := i.service.Send(c.Request().Context(), body, &callerId, platformId)
	if err != nil {
		common.LogStringError(c, err, "invite: send")

		if strings.Contains(err.Error(), "already in use") {
			return ConflictError(c, errors.Cause(err).Error())
		}

		return InternalError(c)
	}
	return c.JSON(http.StatusCreated, m)
}

func (i invite) Accept(c echo.Context) error {
	body := model.RequestInviteAcceptance{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "invite: accept bind")
		return BadRequestError(c)
	}
	id := c.Param("id")

	if !IsValidUUID(id) {
		return BadRequestError(c, "invalid id")
	}

	body.Id = &id
	m, jwt, err := i.service.Accept(c.Request().Context(), body)
	if err != nil {
		common.LogStringError(c, err, "invite: accept")

		if strings.Contains(err.Error(), "invite is not pending") {
			return ConflictError(c, errors.Cause(err).Error())
		}

		return InternalError(c)
	}

	err = SetAuthCookies(c, jwt)
	if err != nil {
		common.LogStringError(c, err, "invite: set auth cookies")
		return InternalError(c)
	}

	return c.JSON(http.StatusOK, m)
}

func (i invite) List(c echo.Context) error {
	platformId := c.Get("platformId").(string)
	status := c.QueryParam("status") // optional
	// if status == "" {
	// 	return BadRequestError(c)
	// }
	m, err := i.service.List(c.Request().Context(), status, platformId)
	if err != nil {
		common.LogStringError(c, err, "invite: list")
		return InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (i invite) Resend(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	id := c.Param("id")
	if !IsValidUUID(id) {
		return BadRequestError(c, "invalid id")
	}

	m, err := i.service.Resend(c.Request().Context(), id, callerId)
	if err != nil {
		common.LogStringError(c, err, "invite: resend")

		if strings.Contains(err.Error(), "not found") {
			return NotFoundError(c, errors.Cause(err).Error())
		}

		return InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (i invite) Update(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	id := c.Param("id")
	if !IsValidUUID(id) {
		return BadRequestError(c, "invalid id")
	}

	body := model.RequestInviteUpdate{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "invite: update bind")
		return BadRequestError(c)
	}

	if err := c.Validate(body); err != nil {
		common.LogStringError(c, err, "invite: update validate")
		return InvalidPayloadError(c, err)
	}

	m, err := i.service.Update(c.Request().Context(), body, id, callerId)
	if err != nil {
		common.LogStringError(c, err, "invite: update")

		if strings.Contains(err.Error(), "not found") {
			return NotFoundError(c, errors.Cause(err).Error())
		}

		if strings.Contains(err.Error(), "cannot elevate") {
			return ForbiddenError(c, errors.Cause(err).Error())
		}

		return InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (i invite) Deactivate(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	id := c.Param("id")
	if !IsValidUUID(id) {
		return BadRequestError(c, "invalid id")
	}

	m, err := i.service.Deactivate(c.Request().Context(), id, callerId)
	if err != nil {
		common.LogStringError(c, err, "invite: update")

		if strings.Contains(err.Error(), "not found") {
			return NotFoundError(c, errors.Cause(err).Error())
		}

		if strings.Contains(err.Error(), "lacks authority") {
			return ForbiddenError(c, errors.Cause(err).Error())
		}

		return InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (i invite) Get(c echo.Context) error {
	id := c.Param("id")
	if !IsValidUUID(id) {
		return BadRequestError(c, "invalid id")
	}

	m, err := i.service.Get(c.Request().Context(), id)
	if err != nil {
		common.LogStringError(c, err, "invite: get")

		if strings.Contains(err.Error(), "not found") {
			return NotFoundError(c, errors.Cause(err).Error())
		}

		return InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (i invite) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("no group attached to the invite handler")
	}
	i.Group = g

	g.POST("", i.Send, ms...)
	g.POST("/:id", i.Accept)
	g.GET("", i.List, ms...)
	g.POST("/:id/resend", i.Resend, ms...)
	g.PUT("/:id", i.Update, ms...)
	g.PUT("/:id/deactivate", i.Deactivate, ms...)
	g.GET("/:id", i.Get)

}
