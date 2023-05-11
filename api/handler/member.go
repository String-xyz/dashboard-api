package handler

import (
	"net/http"
	"strings"

	"github.com/String-xyz/go-lib/common"
	httperror "github.com/String-xyz/go-lib/httperror"
	validator "github.com/String-xyz/go-lib/validator"

	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type Member interface {
	GetAll(e echo.Context) error
	Get(e echo.Context) error
	Update(e echo.Context) error
	Deactivate(e echo.Context) error
	PasswordReset(e echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type member struct {
	service service.Member
	auth    service.Auth
	Group   *echo.Group
}

func NewMember(services service.Services) Member {
	return &member{service: services.Member, auth: services.Auth}
}

// @Summary Get all members
// @Tags Member
// @Produce  json
// @Success 200 {array} repository.OrganizationMemberWithRole
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /members [get]
func (a member) GetAll(c echo.Context) error {
	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid organizationId")
	}

	m, err := a.service.GetAll(c.Request().Context(), organizationId)
	if err != nil {
		common.LogStringError(c, err, "member: get all")
		return httperror.Internal500(c)
	}
	return c.JSON(http.StatusOK, m)
}

// @Summary Get a member by ID
// @Tags Member
// @Produce  json
// @Param id path string true "Member ID"
// @Success 200 {object} repository.OrganizationMemberWithRole
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /members/{id} [get]
func (a member) Get(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid memberId")
	}

	organizationId, ok := c.Get("organizationId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid organizationId")
	}

	memberId := c.Param("id")

	if memberId == "" {
		return httperror.BadRequest400(c)
	}

	m, err := a.service.Get(c.Request().Context(), callerId, organizationId, memberId)
	if err != nil {
		return DefaultErrorHandler(c, err, "member: get")
	}

	return c.JSON(http.StatusOK, m)
}

// @Summary Update member
// @Tags Member
// @Accept  json
// @Produce  json
// @Param id path string true "Member ID"
// @Param body body model.RequestMemberUpdateOther true "Update Member Request"
// @Success 200 {object} repository.OrganizationMemberWithRole
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /members/{id} [patch]
func (a member) Update(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid memberId")
	}

	memberId := c.Param("id")

	if !validator.IsUUID(memberId) {
		return httperror.BadRequest400(c, "invalid id")
	}

	body := model.RequestMemberUpdateOther{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "member: update bind")
		return httperror.BadRequest400(c)
	}

	m, err := a.service.UpdateMember(c.Request().Context(), body, callerId, memberId)
	if err != nil {
		return DefaultErrorHandler(c, err, "member: update")
	}

	return c.JSON(http.StatusOK, m)
}

// @Summary Update self
// @Tags Member
// @Accept  json
// @Produce  json
// @Param body body model.RequestMemberUpdateSelf true "Update Self Request"
// @Success 200 {object} repository.OrganizationMemberWithRole
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /members [patch]
func (a member) UpdateSelf(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid memberId")
	}

	body := model.RequestMemberUpdateSelf{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "member: update self bind")
		return httperror.BadRequest400(c)
	}

	// password must be at least 8 characters
	if body.NewPassword != nil && len(*body.NewPassword) < 8 {
		return httperror.BadRequest400(c, "password must be at least 8 characters")
	}

	// TODO: Fix this validation. This code breaks the update self name endpoint.
	// if c.Validate(body) != nil {
	// 	return httperror.BadRequest400(c)
	// }

	m, err := a.service.UpdateSelf(c.Request().Context(), body, callerId)
	if err != nil {
		return DefaultErrorHandler(c, err, "member: update self")
	}
	return c.JSON(http.StatusOK, m)
}

// @Summary Transfer ownership
// @Tags Member
// @Accept  json
// @Produce  json
// @Param id path string true "Member ID"
// @Param body body model.RequestTransferOwnership true "Transfer Ownership Request"
// @Success 200 {object} repository.OrganizationMemberWithRole
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /members/{id}/transferOwner [patch]
func (a member) TransferOwnership(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid memberId")
	}

	memberId := c.Param("id")

	if !validator.IsUUID(memberId) {
		return httperror.BadRequest400(c, "invalid id")
	}

	body := model.RequestTransferOwnership{}

	if err := c.Bind(&body); err != nil {
		common.LogStringError(c, err, "member: transferOwnership bind")
		return httperror.BadRequest400(c)
	}

	m, err := a.service.TransferOwnership(c.Request().Context(), body, callerId, memberId)
	if err != nil {
		return DefaultErrorHandler(c, err, "member: transferOwnership")
	}

	return c.JSON(http.StatusOK, m)
}

// @Summary Deactivate member
// @Tags Member
// @Produce  json
// @Param id path string true "Member ID"
// @Success 200 {object} repository.OrganizationMemberWithRole
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /members/{id}/deactivate [patch]
func (a member) Deactivate(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid memberId")
	}

	memberId := c.Param("id")

	if !validator.IsUUID(memberId) {
		return httperror.BadRequest400(c, "invalid id")
	}

	if memberId == "" || memberId == callerId {
		return httperror.BadRequest400(c)
	}

	m, err := a.service.Deactivate(c.Request().Context(), callerId, memberId)
	if err != nil {
		return DefaultErrorHandler(c, err, "member: deactivate")
	}

	return c.JSON(http.StatusOK, m)
}

// @Summary Reactivate member
// @Tags Member
// @Produce  json
// @Param id path string true "Member ID"
// @Success 200 {object} repository.OrganizationMemberWithRole
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /members/{id}/reactivate [patch]
func (a member) Reactivate(c echo.Context) error {
	callerId, ok := c.Get("memberId").(string)
	if !ok {
		return httperror.Internal500(c, "missing or invalid memberId")
	}

	memberId := c.Param("id")
	if memberId == "" || memberId == callerId {
		return httperror.BadRequest400(c)
	}

	m, err := a.service.Reactivate(c.Request().Context(), callerId, memberId)
	if err != nil {
		return DefaultErrorHandler(c, err, "member: reactivate")
	}

	return c.JSON(http.StatusOK, m)
}

// @Summary Send password reset email
// @Tags Member
// @Accept  json
// @Produce  json
// @Param body body model.RequestPasswordResetEmail true "Password Reset Email Request"
// @Success 200 {object} map[string]string "Email sent message"
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /members/password-reset [get]
func (a member) SendPasswordResetEmail(c echo.Context) error {
	body := model.RequestPasswordResetEmail{}

	if err := c.Bind(&body); err != nil {
		common.LogStringError(c, err, "member: send password reset email bind")
		return httperror.BadRequest400(c)
	}

	if err := c.Validate(body); err != nil {
		return httperror.InvalidPayload400(c, err)
	}

	body.Email = strings.ToLower(body.Email)

	err := a.service.SendPasswordResetEmail(c.Request().Context(), body)
	if err != nil {
		return DefaultErrorHandler(c, err, "member: send password reset email")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "email sent"})
}

// @Summary Reset password
// @Tags Member
// @Accept  json
// @Produce  json
// @Param body body model.RequestPasswordReset true "Password Reset Request"
// @Success 200 {object} map[string]string "Password reset message"
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /members/password-reset [post]
func (a member) PasswordReset(c echo.Context) error {
	body := model.RequestPasswordReset{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "member: password reset bind")
		return httperror.BadRequest400(c)
	}

	if err := c.Validate(body); err != nil {
		return httperror.InvalidPayload400(c, err)
	}

	err = a.service.PasswordReset(c.Request().Context(), body)
	if err != nil {
		return DefaultErrorHandler(c, err, "member: password reset")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "password reset"})
}

func (a member) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("no group attached to the member handler")
	}
	a.Group = g

	g.GET("", a.GetAll, ms...)
	g.GET("/:id", a.Get, ms...)
	g.PATCH("/:id", a.Update, ms...)
	g.PATCH("", a.UpdateSelf, ms...)
	g.PATCH("/:id/transferOwner", a.TransferOwnership, ms...)
	g.PATCH("/:id/deactivate", a.Deactivate, ms...)
	g.PATCH("/:id/reactivate", a.Reactivate, ms...)
	g.GET("/password-reset", a.SendPasswordResetEmail)
	g.POST("/password-reset", a.PasswordReset)
}
