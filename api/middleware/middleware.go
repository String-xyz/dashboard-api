package middleware

import (
	"os"
	"strings"

	"github.com/String-xyz/platform-admin-api/api/handler"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

func BearerAuth() echo.MiddlewareFunc {
	config := echoMiddleware.JWTConfig{
		TokenLookup: "header:Authorization,cookie:StringJWT",
		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
			var claims = &service.JWTClaims{}
			t, err := jwt.ParseWithClaims(auth, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})

			c.Set("memberId", claims.MemberId)
			c.Set("platformId", claims.PlatformId)
			return t, err
		},
		SigningKey: []byte(os.Getenv("JWT_SECRET_KEY")),
		ErrorHandlerWithContext: func(err error, c echo.Context) error {
			if strings.Contains(err.Error(), "token is expired") {
				return handler.TokenExpired(c)
			}

			if strings.Contains(errors.Cause(err).Error(), "missing or malformed jwt") {
				return handler.MissingToken(c)
			}

			return handler.Unauthorized(c)
		},
	}
	return echoMiddleware.JWTWithConfig(config)
}
