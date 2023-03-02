package handler

import (
	"net/http"
	"os"
	"time"

	httperror "github.com/String-xyz/go-lib/httperror"
	serrors "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/labstack/echo/v4"
)

// TODO: Move all of this into the go-lib

func SetJWTCookie(c echo.Context, jwt service.JWT) error {
	cookie := new(http.Cookie)
	cookie.Name = "StringJWT"
	cookie.Value = jwt.Token
	// cookie.HttpOnly = true // due the short expiration time it is not needed to be http only
	cookie.Expires = jwt.ExpAt // we want the cookie to expire at the same time as the token
	cookie.SameSite = getCookieSameSiteMode()
	cookie.Path = "/"             // Send cookie in every sub path request
	cookie.Secure = !IsLocalEnv() // in production allow https only
	c.SetCookie(cookie)

	return nil
}

func SetRefreshTokenCookie(c echo.Context, refresh service.RefreshTokenResponse) error {
	cookie := new(http.Cookie)
	cookie.Name = "refresh_token"
	cookie.Value = refresh.Token
	cookie.HttpOnly = true
	cookie.Expires = refresh.ExpAt // we want the cookie to expire at the same time as the token
	cookie.SameSite = getCookieSameSiteMode()
	cookie.Path = "/login/"       // Send cookie only in /login path request
	cookie.Secure = !IsLocalEnv() // in production allow https only
	c.SetCookie(cookie)

	return nil
}

func SetAuthCookies(c echo.Context, jwt service.JWT) error {
	err := SetJWTCookie(c, jwt)
	if err != nil {
		return err
	}

	err = SetRefreshTokenCookie(c, jwt.RefreshToken)
	if err != nil {
		return err
	}

	return nil
}

func DeleteAuthCookies(c echo.Context) error {
	// in order to delete a cookie we need to set it with an expired date
	cookie := new(http.Cookie)
	cookie.Name = "StringJWT"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.SameSite = getCookieSameSiteMode()
	cookie.Path = "/" // Send cookie in every sub path request
	cookie.Secure = !IsLocalEnv()
	c.SetCookie(cookie)

	cookie = new(http.Cookie)
	cookie.Name = "refresh_token"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.SameSite = getCookieSameSiteMode()
	cookie.Path = "/login/" // Send cookie only in refresh path request
	cookie.Secure = !IsLocalEnv()
	c.SetCookie(cookie)

	return nil
}

func IsLocalEnv() bool {
	return os.Getenv("ENV") == "local"
}

func getCookieSameSiteMode() http.SameSite {
	sameSiteMode := http.SameSiteNoneMode // allow cors
	if IsLocalEnv() {
		sameSiteMode = http.SameSiteLaxMode // because SameSiteNoneMode is not allowed in localhost we use lax mode
	}
	return sameSiteMode
}

func DefaultErrorHandler(c echo.Context, err error) error {
	if err == nil {
		return nil
	}

	if serrors.ErrorIs(err, serrors.ERR_NOT_FOUND) {
		return httperror.NotFoundError(c)
	}

	if serrors.ErrorIs(err, serrors.ERR_FORBIDDEN) {
		return httperror.ForbiddenError(c, "Invoking member lacks authority")
	}

	if serrors.ErrorIs(err, serrors.ERR_INVALID_RESET_TOKEN) {
		return httperror.BadRequestError(c, "Invalid password reset token")
	}

	if serrors.ErrorIs(err, serrors.ERR_INVALID_PASSWORD) {
		return httperror.BadRequestError(c, "Invalid password")
	}

	if serrors.ErrorIs(err, serrors.ERR_ALREADY_IN_USE) {
		return httperror.ConflictError(c, "Already in use")
	}

	return httperror.InternalError(c)
}

func BindAndValidateBody[T any](c echo.Context, body *T) error {
	err := c.Bind(&body)
	if err != nil {
		return httperror.BadRequestError(c, "invalid payload", "invalid payload", "invite")
	}

	if err := c.Validate(body); err != nil {
		return httperror.InvalidPayloadError(c, err)
	}

	return nil
}
