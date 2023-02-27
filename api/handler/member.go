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

func (a member) GetAll(c echo.Context) error {
	platformId := c.Get("platformId").(string)
	m, err := a.service.GetAll(c.Request().Context(), platformId)
	if err != nil {
		common.LogStringError(c, err, "member: get all")
		return InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (a member) Get(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	platformId := c.Get("platformId").(string)
	memberId := c.Param("id")
	if memberId == "" {
		return BadRequestError(c)
	}
	m, err := a.service.Get(c.Request().Context(), callerId, platformId, memberId)
	if err != nil {
		common.LogStringError(c, err, "member: get")
		return InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (a member) Update(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	memberId := c.Param("id")

	if !IsValidUUID(memberId) {
		return BadRequestError(c, "invalid id")
	}

	body := model.RequestMemberUpdateOther{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "member: update bind")
		return BadRequestError(c)
	}

	m, err := a.service.UpdateMember(c.Request().Context(), body, callerId, memberId)
	if err != nil {
		common.LogStringError(c, err, "member: update")

		if errors.Cause(err).Error() == "invoking member lacks authority" {
			return ForbiddenError(c, "invoking member lacks authority")
		}

		if errors.Cause(err).Error() == "not found" {
			return NotFoundError(c, "member not found")
		}

		return InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (a member) UpdateSelf(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	body := model.RequestMemberUpdateSelf{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "member: update self bind")
		return BadRequestError(c)
	}

	// validate body
	if err := c.Validate(body); err != nil {
		return InvalidPayloadError(c, err)
	}

	m, err := a.service.UpdateSelf(c.Request().Context(), body, callerId)
	if err != nil {
		common.LogStringError(c, err, "member: update self")

		if strings.Contains(errors.Cause(err).Error(), "invalid password") {
			return BadRequestError(c, "invalid password")
		}

		return InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (a member) Deactivate(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	memberId := c.Param("id")

	if !IsValidUUID(memberId) {
		return BadRequestError(c, "invalid id")
	}

	if memberId == "" || memberId == callerId {
		return BadRequestError(c)
	}

	m, err := a.service.Deactivate(c.Request().Context(), callerId, memberId)
	if err != nil {
		common.LogStringError(c, err, "member: deactivate")

		if errors.Cause(err).Error() == "invoking member lacks authority" {
			return ForbiddenError(c, "invoking member lacks authority")
		}

		if errors.Cause(err).Error() == "not found" {
			return NotFoundError(c, "member not found")
		}

		return InternalError(c)
	}

	return c.JSON(http.StatusOK, m)
}

func (a member) Reactivate(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	memberId := c.Param("id")
	if memberId == "" || memberId == callerId {
		return httperror.BadRequestError(c)
	}

	m, err := a.service.Reactivate(c.Request().Context(), callerId, memberId)
	if err != nil {
		common.LogStringError(c, err, "member: deactivate")
		return httperror.InternalError(c)
	}

	return c.JSON(http.StatusOK, m)
}

func (a member) SendPasswordResetEmail(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		return BadRequestError(c, "email is required")
	}
	err := a.service.SendPasswordResetEmail(c.Request().Context(), email)
	if err != nil {
		common.LogStringError(c, err, "member: send password reset email")

		if errors.Cause(err).Error() == "not found" {
			return NotFoundError(c, "member not found")
		}

		return InternalError(c)
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "email sent"})
}

func (a member) PasswordReset(c echo.Context) error {
	body := model.RequestPasswordReset{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "member: password reset bind")
		return BadRequestError(c)
	}

	if err := c.Validate(body); err != nil {
		return InvalidPayloadError(c, err)
	}

	err = a.service.PasswordReset(c.Request().Context(), body)
	if err != nil {
		common.LogStringError(c, err, "member: password reset")

		if errors.Cause(err).Error() == "invalid password reset token" {
			return BadRequestError(c, "invalid password reset token")
		}

		if strings.Contains(errors.Cause(err).Error(), "invalid password") {
			return BadRequestError(c, "invalid password")
		}

		return InternalError(c)
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
	g.PUT("/:id", a.Update, ms...)
	g.PUT("", a.UpdateSelf, ms...)
	g.PUT("/:id/deactivate", a.Deactivate, ms...)
	g.PUT("/:id/reactivate", a.Reactivate, ms...)
	g.GET("/password-reset", a.SendPasswordResetEmail)
	g.POST("/password-reset", a.PasswordReset)
}
