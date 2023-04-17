package handler

import (
	"net/http"
	"strings"

	"github.com/String-xyz/go-lib/common"

	httperror "github.com/String-xyz/go-lib/httperror"
	serror "github.com/String-xyz/go-lib/stringerror"
	validator "github.com/String-xyz/go-lib/validator"
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
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid callerId")
	}

	platformId, ok := c.Get("platformId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid platformId")
	}

	body := model.RequestInviteSend{}
	err := c.Bind(&body)
	if err != nil {
		return httperror.BadRequestError(c, "invalid payload")
	}

	body.Email = strings.ToLower(body.Email)

	if err := c.Validate(body); err != nil {
		common.LogStringError(c, err, "invite: send validate")
		return httperror.InvalidPayloadError(c, err)
	}

	m, err := i.service.Send(c.Request().Context(), body, &callerId, platformId)
	if err != nil {
		return DefaultErrorHandler(c, err, "invite: send")
	}

	return c.JSON(http.StatusCreated, m)
}

func (i invite) Accept(c echo.Context) error {
	id := c.Param("id")

	body := model.RequestInviteAcceptance{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "invite: accept bind")
		return httperror.BadRequestError(c)
	}

	if err := c.Validate(body); err != nil {
		return httperror.InvalidPayloadError(c, err)
	}

	if !validator.IsUUID(id) {
		return httperror.BadRequestError(c, "invalid id")
	}

	m, jwt, err := i.service.Accept(c.Request().Context(), id, body)
	if err != nil {
		common.LogStringError(c, err, "invite: accept")

		if serror.Is(err, serror.ALREADY_IN_USE) {
			return httperror.ConflictError(c, "Invite is not pending")
		}

		// the token is invalid or expired
		if serror.Is(err, serror.FORBIDDEN) {
			return httperror.BadRequestError(c, "Invalid password reset token")
		}

		return httperror.InternalError(c)
	}

	err = SetAuthCookies(c, jwt)
	if err != nil {
		common.LogStringError(c, err, "invite: set auth cookies")
		return httperror.InternalError(c)
	}

	return c.JSON(http.StatusOK, m)
}

func (i invite) List(c echo.Context) error {
	platformId, ok := c.Get("platformId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid platformId")
	}

	status := c.QueryParam("status") // optional
	// if status == "" {
	// 	return httperror.BadRequestError(c)
	// }
	m, err := i.service.List(c.Request().Context(), status, platformId)
	if err != nil {
		return DefaultErrorHandler(c, err, "invite: list")
	}
	return c.JSON(http.StatusOK, m)
}

func (i invite) Resend(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid callerId")
	}

	id := c.Param("id")
	if !validator.IsUUID(id) {
		return httperror.BadRequestError(c, "invalid id")
	}

	m, err := i.service.Resend(c.Request().Context(), id, callerId)
	if err != nil {
		common.LogStringError(c, err, "invite: resend")

		if serror.Is(err, serror.NOT_FOUND) {
			return httperror.NotFoundError(c, errors.Cause(err).Error())
		}

		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (i invite) Update(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid memberId")
	}

	id := c.Param("id")
	if !validator.IsUUID(id) {
		return httperror.BadRequestError(c, "invalid id")
	}

	body := model.RequestInviteUpdate{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "invite: update bind")
		return httperror.BadRequestError(c)
	}

	if err := c.Validate(body); err != nil {
		common.LogStringError(c, err, "invite: update validate")
		return httperror.InvalidPayloadError(c, err)
	}

	m, err := i.service.Update(c.Request().Context(), body, id, callerId)
	if err != nil {
		common.LogStringError(c, err, "invite: update")

		if serror.Is(err, serror.NOT_FOUND) {
			return httperror.NotFoundError(c, errors.Cause(err).Error())
		}

		if serror.Is(err, serror.FORBIDDEN) {
			return httperror.ForbiddenError(c, "cannot elevate member to owner")
		}

		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (i invite) Deactivate(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.InternalError(c, "missing or invalid memberId")
	}

	id := c.Param("id")
	if !validator.IsUUID(id) {
		return httperror.BadRequestError(c, "invalid id")
	}

	m, err := i.service.Deactivate(c.Request().Context(), id, callerId)
	if err != nil {
		common.LogStringError(c, err, "invite: update")

		if serror.Is(err, serror.NOT_FOUND) {
			return httperror.NotFoundError(c, errors.Cause(err).Error())
		}

		if serror.Is(err, serror.FORBIDDEN) {
			return httperror.ForbiddenError(c, errors.Cause(err).Error())
		}

		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (i invite) Get(c echo.Context) error {
	id := c.Param("id")
	if !validator.IsUUID(id) {
		return httperror.BadRequestError(c, "invalid id")
	}

	m, err := i.service.Get(c.Request().Context(), id)
	if err != nil {
		return DefaultErrorHandler(c, err, "invite: get")
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
	g.PATCH("/:id", i.Update, ms...)
	g.PATCH("/:id/deactivate", i.Deactivate, ms...)
	g.GET("/:id", i.Get) // No auth required. This is used for the invite link

}
