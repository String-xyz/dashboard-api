package handler

import (
	"net/http"
	"strings"

	"github.com/String-xyz/go-lib/common"
	httperror "github.com/String-xyz/go-lib/httperror"
	serror "github.com/String-xyz/go-lib/stringerror"
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

	body.Email = strings.ToLower(body.Email)

	// validate body
	if err := c.Validate(body); err != nil {
		return httperror.InvalidPayloadError(c, err)
	}

	member, jwt, err := l.service.Login(c.Request().Context(), body)
	if err != nil {
		common.LogStringError(c, err, "login: login")

		if serror.Is(err, serror.DEACTIVATED, serror.INVALID_PASSWORD, serror.NOT_FOUND) {
			return httperror.Unauthorized(c, "Invalid email or password")
		}

		if serror.Is(err, serror.NOT_FOUND) {
			return httperror.Unauthorized(c, "Invalid email or password")
		}

		return DefaultErrorHandler(c, err, "login: login")
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
		return httperror.Unauthorized(c, "Invalid or expired token")
	}

	resp, err := l.auth.RefreshToken(cookie.Value)
	if err != nil {
		common.LogStringError(c, err, "login: refresh token")
		return httperror.Unauthorized(c, "Invalid or expired token")
	}

	// If member is denied, do not refresh token
	denied, err := l.auth.IsDenied(resp.Member.ID)
	if err != nil {
		common.LogStringError(c, err, "login: fail denylist check")
		return httperror.InternalError(c)
	}
	if denied {
		common.LogStringError(c, err, "login: user is denied")
		return httperror.Unauthorized(c)
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
		// already logged out, idempotent
		return c.JSON(http.StatusNoContent, nil)
	}

	// invalidate refresh token. Returns error if token is not found
	err = l.auth.InvalidateRefreshToken(cookie.Value)
	if err != nil {
		common.LogStringError(c, err, "Token not found")
	}

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
