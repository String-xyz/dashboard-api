package middleware

import (
	"strings"

	validator "github.com/String-xyz/go-lib/v2/validator"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	httperror "github.com/String-xyz/go-lib/v2/httperror"
	"github.com/String-xyz/platform-admin-api/config"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/pkg/errors"
)

func JWT(auth service.Auth) echo.MiddlewareFunc {
	config := echoMiddleware.JWTConfig{
		TokenLookup: "cookie:StringAdminJWT",
		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
			var claims = &service.JWTClaims{}
			t, err := jwt.ParseWithClaims(auth, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(config.Var.JWT_SECRET_KEY), nil
			})

			// validate claims
			if !validator.IsUUID(claims.MemberId, claims.OrganizationId) {
				return nil, errors.New("missing or malformed jwt")
			}
			c.Set("organizationId", claims.OrganizationId)
			c.Set("memberId", claims.MemberId)

			return t, err
		},
		SigningKey: []byte(config.Var.JWT_SECRET_KEY),
		ErrorHandlerWithContext: func(err error, c echo.Context) error {
			if strings.Contains(err.Error(), "token is expired") || strings.Contains(err.Error(), "missing or malformed jwt") {
				return httperror.Unauthorized401(c)
			}

			// If member is denied, do not honor JWT
			var claims = &service.JWTClaims{}
			denied, err := auth.IsDenied(claims.MemberId)
			if err != nil {
				return httperror.Internal500(c)
			}
			if denied {
				return httperror.Unauthorized401(c)
			}

			return httperror.Unauthorized401(c)
		},
	}
	return echoMiddleware.JWTWithConfig(config)
}

func APIKeySecretAuth(service service.Auth) echo.MiddlewareFunc {
	config := echoMiddleware.KeyAuthConfig{
		KeyLookup: "header:X-Api-Key",
		Validator: func(auth string, c echo.Context) (bool, error) {
			apikey, err := service.ValidateAPIKeySecret(c.Request().Context(), auth)
			if err != nil {
				libcommon.LogStringError(c, err, "Error in APIKeySecretAuth middleware")
				return false, err
			}

			c.Set("organizationId", apikey.OrganizationId)
			c.Set("memberId", apikey.CreatedBy)

			return true, nil
		},
	}
	return echoMiddleware.KeyAuthWithConfig(config)
}

func APIKeyOrJWTAuth(service service.Auth) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiKeyAuth := APIKeySecretAuth(service)(next)
			jwtAuth := JWT(service)(next)

			err := apiKeyAuth(c)
			if err == nil {
				return nil
			}

			err = jwtAuth(c)
			if err == nil {
				return nil
			}

			libcommon.LogStringError(c, err, "Error in APIKeyOrJWTAuth middleware")
			return httperror.Unauthorized401(c)
		}
	}
}
