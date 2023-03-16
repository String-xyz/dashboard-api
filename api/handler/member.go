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
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusOK, m)
}

func (a member) Get(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	platformId := c.Get("platformId").(string)
	memberId := c.Param("id")
	if memberId == "" {
		return httperror.BadRequestError(c)
	}
	m, err := a.service.Get(c.Request().Context(), callerId, platformId, memberId)
	if err != nil {
		common.LogStringError(c, err, "member: get")
		return DefaultErrorHandler(c, err)
	}
	return c.JSON(http.StatusOK, m)
}

func (a member) Update(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	memberId := c.Param("id")

	if !validator.IsUUID(memberId) {
		return httperror.BadRequestError(c, "invalid id")
	}

	body := model.RequestMemberUpdateOther{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "member: update bind")
		return httperror.BadRequestError(c)
	}

	m, err := a.service.UpdateMember(c.Request().Context(), body, callerId, memberId)
	if err != nil {
		common.LogStringError(c, err, "member: update")

		return DefaultErrorHandler(c, err)
	}

	return c.JSON(http.StatusOK, m)
}

func (a member) UpdateSelf(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	body := model.RequestMemberUpdateSelf{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "member: update self bind")
		return httperror.BadRequestError(c)
	}

	// password must be at least 8 characters
	if body.NewPassword != nil && len(*body.NewPassword) < 8 {
		return httperror.BadRequestError(c, "password must be at least 8 characters")
	}

	// TODO: Fix this validation. This code breaks the update self name endpoint.
	// if c.Validate(body) != nil {
	// 	return httperror.BadRequestError(c)
	// }

	m, err := a.service.UpdateSelf(c.Request().Context(), body, callerId)
	if err != nil {
		common.LogStringError(c, err, "member: update self")

		if serror.IsError(err, serror.NOT_FOUND) {
			/* Since we get the callerId from the token, this should never happen. If it does, it means the token is invalid */
			return httperror.Unauthorized(c)
		}

		return DefaultErrorHandler(c, err)
	}
	return c.JSON(http.StatusOK, m)
}

func (a member) TransferOwnership(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	memberId := c.Param("id")

	if !validator.IsUUID(memberId) {
		return httperror.BadRequestError(c, "invalid id")
	}

	body := model.RequestTransferOwnership{}

	if err := c.Bind(&body); err != nil {
		common.LogStringError(c, err, "member: transferOwnership bind")
		return httperror.BadRequestError(c)
	}

	m, err := a.service.TransferOwnership(c.Request().Context(), body, callerId, memberId)
	if err != nil {
		common.LogStringError(c, err, "member: transferOwnership")

		return DefaultErrorHandler(c, err)
	}

	return c.JSON(http.StatusOK, m)
}

func (a member) Deactivate(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	memberId := c.Param("id")

	if !validator.IsUUID(memberId) {
		return httperror.BadRequestError(c, "invalid id")
	}

	if memberId == "" || memberId == callerId {
		return httperror.BadRequestError(c)
	}

	m, err := a.service.Deactivate(c.Request().Context(), callerId, memberId)
	if err != nil {
		common.LogStringError(c, err, "member: deactivate")

		return DefaultErrorHandler(c, err)
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

		return DefaultErrorHandler(c, err)
	}

	return c.JSON(http.StatusOK, m)
}

func (a member) SendPasswordResetEmail(c echo.Context) error {
	body := model.RequestPasswordResetEmail{}

	if err := c.Bind(&body); err != nil {
		common.LogStringError(c, err, "member: send password reset email bind")
		return httperror.BadRequestError(c)
	}

	if err := c.Validate(body); err != nil {
		return httperror.InvalidPayloadError(c, err)
	}

	body.Email = strings.ToLower(body.Email)

	err := a.service.SendPasswordResetEmail(c.Request().Context(), body)
	if err != nil {
		common.LogStringError(c, err, "member: send password reset email")

		return DefaultErrorHandler(c, err)
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "email sent"})
}

func (a member) PasswordReset(c echo.Context) error {
	body := model.RequestPasswordReset{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "member: password reset bind")
		return httperror.BadRequestError(c)
	}

	if err := c.Validate(body); err != nil {
		return httperror.InvalidPayloadError(c, err)
	}

	err = a.service.PasswordReset(c.Request().Context(), body)
	if err != nil {
		common.LogStringError(c, err, "member: password reset")

		return DefaultErrorHandler(c, err)
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
	g.PUT("/:id/transferOwner", a.TransferOwnership, ms...)
	g.PUT("/:id/deactivate", a.Deactivate, ms...)
	g.PUT("/:id/reactivate", a.Reactivate, ms...)
	g.GET("/password-reset", a.SendPasswordResetEmail)
	g.POST("/password-reset", a.PasswordReset)
}

// custom error messages
