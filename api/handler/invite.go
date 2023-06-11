package handler

import (
	"net/http"
	"strings"

	"github.com/String-xyz/go-lib/v2/common"

	"github.com/String-xyz/dashboard-api/pkg/model"
	"github.com/String-xyz/dashboard-api/pkg/service"
	httperror "github.com/String-xyz/go-lib/v2/httperror"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
	validator "github.com/String-xyz/go-lib/v2/validator"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Invite interface {
	Send(e echo.Context) error
	Accept(e echo.Context) error
	List(e echo.Context) error
	Resend(e echo.Context) error
	Update(e echo.Context) error
	Revoke(e echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type invite struct {
	service service.Invite
	Group   *echo.Group
}

func NewInvite(service service.Invite) Invite {
	return &invite{service: service}
}

// @Summary Send an invite
// @Tags Invite
// @Accept  json
// @Produce  json
// @Param RequestInviteSend body model.RequestInviteSend true "Invite Send Request"
// @Success 201 {object} repository.MemberInviteInfo
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /invites [post]
func (i invite) Send(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid callerId")
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid organizationId")
	}

	body := model.RequestInviteSend{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "invite: send bind")
		return httperror.BadRequest400(c, "invalid payload")
	}

	body.Email = strings.ToLower(body.Email)

	if err := c.Validate(body); err != nil {
		common.LogStringError(c, err, "invite: send validate")
		return httperror.InvalidPayload400(c, err)
	}

	m, err := i.service.Send(c.Request().Context(), body, &callerId, organizationId)
	if err != nil {
		return DefaultErrorHandler(c, err, "invite: send")
	}

	return c.JSON(http.StatusCreated, m)
}

// @Summary Accept an invite
// @Tags Invite
// @Accept  json
// @Produce  json
// @Param id path string true "Invite ID"
// @Param RequestInviteAcceptance body model.RequestInviteAcceptance true "Invite Acceptance Request"
// @Success 200 {object} model.OrganizationMember
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /invites/{id} [post]
func (i invite) Accept(c echo.Context) error {
	id := c.Param("id")

	body := model.RequestInviteAcceptance{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "invite: accept bind")
		return httperror.BadRequest400(c)
	}

	if err := c.Validate(body); err != nil {
		return httperror.InvalidPayload400(c, err)
	}

	if !validator.IsUUID(id) {
		return httperror.BadRequest400(c, "invalid id")
	}

	m, jwt, err := i.service.Accept(c.Request().Context(), id, body)
	if err != nil {
		common.LogStringError(c, err, "invite: accept")

		if serror.Is(err, serror.ALREADY_IN_USE) {
			return httperror.Conflict409(c, "Invite is not pending")
		}

		// the token is invalid or expired
		if serror.Is(err, serror.FORBIDDEN) {
			return httperror.BadRequest400(c, "Invalid password reset token")
		}

		return httperror.Internal500(c)
	}

	err = SetAuthCookies(c, jwt)
	if err != nil {
		common.LogStringError(c, err, "invite: set auth cookies")
		return httperror.Internal500(c)
	}

	return c.JSON(http.StatusOK, m)
}

// @Summary List invites
// @Tags Invite
// @Produce  json
// @Param status query string false "Invite status"
// @Success 200 {array} repository.MemberInviteInfo
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /invites [get]
func (i invite) List(c echo.Context) error {
	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid organizationId")
	}

	status := c.QueryParam("status")
	m, err := i.service.List(c.Request().Context(), status, organizationId)
	if err != nil {
		return DefaultErrorHandler(c, err, "invite: list")
	}
	return c.JSON(http.StatusOK, m)
}

// @Summary Resend an invite
// @Tags Invite
// @Produce  json
// @Param id path string true "Invite ID"
// @Success 200 {object} repository.MemberInviteInfo
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /invites/{id}/resend [post]
func (i invite) Resend(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid callerId")
	}

	id := c.Param("id")
	if !validator.IsUUID(id) {
		return httperror.BadRequest400(c, "invalid id")
	}

	m, err := i.service.Resend(c.Request().Context(), id, callerId)
	if err != nil {
		common.LogStringError(c, err, "invite: resend")

		if serror.Is(err, serror.NOT_FOUND) {
			return httperror.NotFound404(c, errors.Cause(err).Error())
		}

		return httperror.Internal500(c)
	}
	return c.JSON(http.StatusOK, m)
}

// @Summary Update an invite
// @Tags Invite
// @Accept  json
// @Produce  json
// @Param id path string true "Invite ID"
// @Param RequestInviteUpdate body model.RequestInviteUpdate true "Invite Update Request"
// @Success 200 {object} repository.MemberInviteInfo
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 403 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /invites/{id} [patch]
func (i invite) Update(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid memberId")
	}

	id := c.Param("id")
	if !validator.IsUUID(id) {
		return httperror.BadRequest400(c, "invalid id")
	}

	body := model.RequestInviteUpdate{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "invite: update bind")
		return httperror.BadRequest400(c)
	}

	if err := c.Validate(body); err != nil {
		common.LogStringError(c, err, "invite: update validate")
		return httperror.InvalidPayload400(c, err)
	}

	m, err := i.service.Update(c.Request().Context(), body, id, callerId)
	if err != nil {
		common.LogStringError(c, err, "invite: update")

		if serror.Is(err, serror.NOT_FOUND) {
			return httperror.NotFound404(c, errors.Cause(err).Error())
		}

		if serror.Is(err, serror.FORBIDDEN) {
			return httperror.Forbidden403(c, "cannot elevate member to owner")
		}

		return httperror.Internal500(c)
	}
	return c.JSON(http.StatusOK, m)
}

// @Summary Revoke an invite
// @Tags Invite
// @Produce  json
// @Param id path string true "Invite ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 403 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /invites/{id} [delete]
func (i invite) Revoke(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid memberId")
	}

	id := c.Param("id")
	if !validator.IsUUID(id) {
		return httperror.BadRequest400(c, "invalid id")
	}

	err := i.service.Revoke(c.Request().Context(), id, callerId)
	if err != nil {
		common.LogStringError(c, err, "invite: update")

		if serror.Is(err, serror.NOT_FOUND) {
			return httperror.NotFound404(c, errors.Cause(err).Error())
		}

		if serror.Is(err, serror.FORBIDDEN) {
			return httperror.Forbidden403(c, errors.Cause(err).Error())
		}

		return httperror.Internal500(c)
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "success"})
}

// @Summary Get an invite
// @Tags Invite
// @Produce  json
// @Param id path string true "Invite ID"
// @Success 200 {object} repository.MemberInviteInfo
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /invites/{id} [get]
func (i invite) Get(c echo.Context) error {
	id := c.Param("id")
	if !validator.IsUUID(id) {
		return httperror.BadRequest400(c, "invalid id")
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
	g.DELETE("/:id", i.Revoke, ms...)
	g.GET("/:id", i.Get) // No auth required. This is used for the invite link

}
