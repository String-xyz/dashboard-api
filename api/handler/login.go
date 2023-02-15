package handler

import (
	"net/http"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/httperror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/labstack/echo/v4"
)

type Login interface {
	Login(e echo.Context) error
	RefreshToken(c echo.Context) error
	Logout(c echo.Context) error
	RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc)
}

type login struct {
	service service.Login
	auth    service.Auth
	group   *echo.Group
}

func NewLogin(service service.Services) Login {
	return &login{service: service.Login, auth: service.Auth}
}

func (l login) Login(c echo.Context) error {
	body := model.RequestLogin{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "login: login bind")
		return httperror.BadRequestError(c)
	}

	member, jwt, err := l.service.Login(c.Request().Context(), body)
	if err != nil {
		common.LogStringError(c, err, "login: login")
		return httperror.InternalError(c)
	}
	// If member is denied, do not log in
	denied, err := l.auth.IsDenied(member.ID)
	if err != nil {
		common.LogStringError(c, err, "login: denylist failed")
		return httperror.InternalError(c)
	}
	if denied {
		common.LogStringError(c, err, "login: access denied")
		return httperror.Unauthorized(c)
	}

	err = SetAuthCookies(c, jwt)
	if err != nil {
		common.LogStringError(c, err, "login: set auth cookies")
		return httperror.InternalError(c)
	}

	return c.JSON(http.StatusOK, member)
}

func (l login) RefreshToken(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		common.LogStringError(c, err, "RefreshToken: unable to get refresh_token cookie")
		return httperror.Unauthorized(c)
	}

	resp, err := l.auth.RefreshToken(cookie.Value)
	if err != nil {
		common.LogStringError(c, err, "login: refresh token")
		return httperror.BadRequestError(c, "Invalid or expired token")
	}

	// If member is denied, do not refresh token
	denied, err := l.auth.IsDenied(resp.Member.ID)
	if err != nil {
		return httperror.InternalError(c)
	}
	if denied {
		return httperror.BadRequestError(c)
	}

	// set auth in cookies
	err = SetAuthCookies(c, resp.JWT)
	if err != nil {
		common.LogStringError(c, err, "RefreshToken: unable to set auth cookies")
		return httperror.InternalError(c)
	}

	return c.JSON(http.StatusOK, resp)
}

// logout
func (l login) Logout(c echo.Context) error {
	// get refresh token from cookie
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		common.LogStringError(c, err, "Logout: unable to get refresh_token cookie")
		return httperror.Unauthorized(c)
	}

	// invalidate refresh token. Returns error if token is not found
	err = l.auth.InvalidateRefreshToken(cookie.Value)
	if err != nil {
		common.LogStringError(c, err, "Token not found")
	}
	// There is no need to invalidate the access token since it is a short lived token

	// delete auth cookies
	err = DeleteAuthCookies(c)
	if err != nil {
		common.LogStringError(c, err, "Logout: unable to delete auth cookies")
		return httperror.InternalError(c)
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (l login) RegisterRoutes(g *echo.Group, ms ...echo.MiddlewareFunc) {
	if g == nil {
		panic("no group attached to the login handler")
	}
	l.group = g
	g.POST("", l.Login)
	g.POST("/refresh", l.RefreshToken)
	g.POST("/logout", l.Logout)
}
