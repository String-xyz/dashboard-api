package handler

import (
	"net/http"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/httperror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type Member interface {
	GetAll(e echo.Context) error
	Get(e echo.Context) error
	Update(e echo.Context) error
	PasswordReset(e echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type member struct {
	service service.Member
	Group   *echo.Group
}

func NewMember(service service.Member) Member {
	return &member{service: service}
}

func (a member) GetAll(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	platformId := c.Get("platformId").(string)
	m, err := a.service.GetAll(c.Request().Context(), callerId, platformId)
	if err != nil {
		common.LogStringError(c, err, "member: get all")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusAccepted, m)
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
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusAccepted, m)
}

func (a member) Update(c echo.Context) error {
	callerId := c.Get("memberId").(string)
	memberId := c.Param("id")
	if memberId == "" {
		return httperror.BadRequestError(c)
	}
	body := model.RequestMemberUpdate{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "member: update bind")
		return httperror.BadRequestError(c)
	}

	m, err := a.service.Update(c.Request().Context(), body, callerId, memberId)
	if err != nil {
		common.LogStringError(c, err, "member: update")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusAccepted, m)
}

func (a member) SendPasswordResetEmail(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		return httperror.BadRequestError(c)
	}
	err := a.service.SendPasswordResetEmail(c.Request().Context(), email)
	if err != nil {
		common.LogStringError(c, err, "member: send password reset email")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusAccepted, nil)
}

func (a member) PasswordReset(c echo.Context) error {
	body := model.RequestPasswordReset{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "member: password reset bind")
		return httperror.BadRequestError(c)
	}

	err = a.service.PasswordReset(c.Request().Context(), body)
	if err != nil {
		common.LogStringError(c, err, "member: password reset")
		return httperror.InternalError(c)
	}
	return c.JSON(http.StatusAccepted, nil)
}

func (a member) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("no group attached to the member handler")
	}
	a.Group = g

	g.Use(ms...)
	g.GET("", a.GetAll, ms...)
	g.GET("/:id", a.Get, ms...)
	g.PUT("/:id", a.Update, ms...)
	g.GET("/password-reset", a.SendPasswordResetEmail)
	g.POST("/password-reset", a.PasswordReset)
}
