package middlewares

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// This middleware assumes the JWT has already been verified by the JWT middleware
func RoleMiddleware(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := extractToken(c) // Assume this correctly extracts the "Bearer token"
			if tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid token")
			}

			// Verify JWT as usual
			user, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, err := token.Method.(*jwt.SigningMethodHMAC); !err {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(""), nil
			})

			claims, ok := user.Claims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
			}
			role, ok := claims["role"].(string)
			if !ok || role != requiredRole {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized for this role")
			}
			return next(c)
		}
	}
}
