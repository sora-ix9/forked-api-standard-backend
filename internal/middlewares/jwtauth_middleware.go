package middlewares

import (
	"fdlp-standard-api/pkg/config"
	"fdlp-standard-api/pkg/redisclient"
	"fmt"

	"net/http"
	"strings"

	"fdlp-standard-api/pkg/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

var redisBreaker = utils.NewCircuitBreaker("Redis-Auth")

func JWTAuthMiddleware(next echo.HandlerFunc, cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := extractToken(c) // Assume this correctly extracts the "Bearer token"
		if tokenString == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid token")
		}

		rdb := redisclient.NewClient(cfg)
		ctx := c.Request().Context() // Correctly use Echo's request context

		// Check if the token is blacklisted using Circuit Breaker
		_, err := utils.ExecuteWithBreaker(redisBreaker, func() (interface{}, error) {
			result, err := rdb.Get(ctx, tokenString).Result()
			if err == nil && result == "blacklisted" {
				return nil, fmt.Errorf("token blacklisted")
			}
			return nil, err
		})

		if err != nil && err.Error() == "token blacklisted" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}
		// If it's a circuit breaker error or other redis error, we might want to continue or fail.
		// Usually, if Redis is down, we might allow the request if it's not a critical security risk,
		// or fail if it's mandatory. Here we'll just log and continue for now if it's a breaker error.

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
