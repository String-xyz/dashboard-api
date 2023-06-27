package handler

import (
	"net/http"
	"strings"

	"github.com/String-xyz/dashboard-api/pkg/model"
	"github.com/String-xyz/dashboard-api/pkg/service"
	"github.com/String-xyz/go-lib/v2/common"
	httperror "github.com/String-xyz/go-lib/v2/httperror"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
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

// @Summary Login
// @Tags Login
// @Accept  json
// @Produce  json
// @Param body body model.RequestLogin true "Login Request"
// @Success 200 {object} repository.OrganizationMemberWithRole
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /login [post]
func (l login) Login(c echo.Context) error {
	body := model.RequestLogin{}
	err := c.Bind(&body)
	if err != nil {
		common.LogStringError(c, err, "login: login bind")
		return httperror.BadRequest400(c)
	}

	body.Email = strings.ToLower(body.Email)

	// validate body
	if err := c.Validate(body); err != nil {
		common.LogStringError(c, err, "login: login validate body")
		return httperror.InvalidPayload400(c, err)
	}

	member, jwt, err := l.service.Login(c.Request().Context(), body)
	if err != nil {
		common.LogStringError(c, err, "login: login")

		if serror.Is(err, serror.DEACTIVATED, serror.INVALID_PASSWORD, serror.NOT_FOUND) {
			return httperror.Unauthorized401(c, "Invalid email or password")
		}

		if serror.Is(err, serror.NOT_FOUND) {
			return httperror.Unauthorized401(c, "Invalid email or password")
		}
	}

	err = SetAuthCookies(c, jwt)
	if err != nil {
		common.LogStringError(c, err, "login: set auth cookies")
		return httperror.Internal500(c)
	}

	return c.JSON(http.StatusOK, member)
}

// @Summary Refresh token
// @Tags Login
// @Produce  json
// @Success 200 {object} service.MemberCreateResponse
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /login/refresh [post]
func (l login) RefreshToken(c echo.Context) error {
	cookie, err := c.Cookie("StringAdminRefreshToken")
	if err != nil {
		common.LogStringError(c, err, "RefreshToken: unable to get StringAdminRefreshToken cookie")
		return httperror.Unauthorized401(c, "Invalid or expired token")
	}

	response, err := l.auth.RefreshToken(cookie.Value)
	if err != nil {
		common.LogStringError(c, err, "login: refresh token")
		return httperror.Unauthorized401(c, "Invalid or expired token")
	}

	// If member is denied, do not refresh token
	denied, err := l.auth.IsDenied(response.Member.Id)
	if err != nil {
		common.LogStringError(c, err, "login: fail denylist check")
		return httperror.Internal500(c)
	}
	if denied {
		common.LogStringError(c, err, "login: user is denied")
		return httperror.Unauthorized401(c)
	}

	// set auth in cookies
	err = SetAuthCookies(c, response.JWT)
	if err != nil {
		common.LogStringError(c, err, "RefreshToken: unable to set auth cookies")
		return httperror.Internal500(c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary Logout
// @Tags Login
// @Success 204
// @Failure 500 {object} error
// @Router /login/logout [post]
func (l login) Logout(c echo.Context) error {
	// get refresh token from cookie
	cookie, err := c.Cookie("StringAdminRefreshToken")
	if err != nil {
		// already logged out, idempotent
		common.LogStringError(c, err, "Logout: unable to get StringAdminRefreshToken cookie")
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
		return httperror.Internal500(c)
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
