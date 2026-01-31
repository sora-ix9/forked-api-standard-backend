package middlewares

import (
	"fdlp-standard-api/pkg/config"
	"fdlp-standard-api/pkg/redisclient"
	"fmt"

	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTAuthMiddleware(next echo.HandlerFunc, cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := extractToken(c) // Assume this correctly extracts the "Bearer token"
		if tokenString == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid token")
		}

		rdb := redisclient.NewClient(cfg)
		ctx := c.Request().Context() // Correctly use Echo's request context

		// Check if the token is blacklisted
		result, err := rdb.Get(ctx, tokenString).Result()
		if err == nil && result == "blacklisted" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		// Verify JWT as usual
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		// Your usual token validation logic
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user", claims)
			return next(c)
		}

		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}
}

func extractToken(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Split the header to separate "Bearer" from the "<token>" part
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
