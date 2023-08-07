package controller

import (
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Permission int

const (
	PermissionNone Permission = iota
	PermissionAdmin
	PermissionVisitor
)

type ContextKey string

const (
	ContextKeyPermission ContextKey = "permission"
)

func CheckPermissionMiddleware(publicKeySrc []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authz := c.Request().Header.Get("Authorization")

			auth := func() Permission {
				if authz == "" || !strings.HasPrefix(authz, "Bearer ") {
					return PermissionVisitor
				}

				access_token_str := strings.TrimPrefix(authz, "Bearer ")

				// parse access token
				access_token, err := jwt.Parse(access_token_str, func(token *jwt.Token) (interface{}, error) {
					return jwt.ParseRSAPublicKeyFromPEM(publicKeySrc)
				})
				if err != nil {
					return PermissionVisitor
				}

				// check token claims
				claims, ok := access_token.Claims.(jwt.MapClaims)
				if !ok || !access_token.Valid {
					return PermissionVisitor
				}

				// check token type (is access token?)
				token_type, _ := claims["type"].(string)
				if token_type != "access" {
					return PermissionVisitor
				}

				return PermissionAdmin
			}()

			c.Set(string(ContextKeyPermission), auth)

			return next(c)
		}
	}
}
