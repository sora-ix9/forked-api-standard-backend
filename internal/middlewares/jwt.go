package middlewares

import (
	"fdlp-standard-api/internal/dto"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware returns the JWT middleware with custom configuration
func JWTMiddleware(secret string) echo.MiddlewareFunc {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(dto.JwtCustomClaims)
		},
		SigningKey: []byte(secret),
	}
	return echojwt.WithConfig(config)
}
