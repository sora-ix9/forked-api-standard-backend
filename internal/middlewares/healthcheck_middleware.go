package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"fdlp-standard-api/pkg/utils"

	"github.com/labstack/echo/v4"
	"github.com/sony/gobreaker"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service,omitempty"`
	Error   string `json:"error,omitempty"`
}

// --- MongoDB Health Check ---

type MongoPingFunc func(ctx context.Context, rp *readpref.ReadPref) error

type MongoHealthCheckConfig struct {
	PingFunc  MongoPingFunc
	Breaker   *gobreaker.TwoStepCircuitBreaker
	SkipPaths map[string]bool
	Timeout   time.Duration
	Endpoint  string
}

func MongoHealthCheckMiddleware(config MongoHealthCheckConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Register endpoint
			if c.Request().URL.Path == config.Endpoint {
				_, err := utils.ExecuteWithBreaker(config.Breaker, func() (interface{}, error) {
					ctx, cancel := context.WithTimeout(c.Request().Context(), config.Timeout)
					defer cancel()
					return nil, config.PingFunc(ctx, readpref.Primary())
				})
				if err != nil {
					return c.JSON(http.StatusServiceUnavailable, healthResponse{
						Status:  "fail",
						Service: "mongo",
						Error:   err.Error(),
					})
				}
				return c.JSON(http.StatusOK, healthResponse{Status: "ok", Service: "mongo"})
			}

			// Skip health check for specified paths
			if config.SkipPaths[c.Request().URL.Path] {
				return next(c)
			}

			// Check health on every request
			_, err := utils.ExecuteWithBreaker(config.Breaker, func() (interface{}, error) {
				ctx, cancel := context.WithTimeout(c.Request().Context(), config.Timeout)
				defer cancel()
				return nil, config.PingFunc(ctx, readpref.Primary())
			})
			if err != nil {
				c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c.Response().WriteHeader(http.StatusServiceUnavailable)
				return json.NewEncoder(c.Response().Writer).Encode(healthResponse{
					Status:  "fail",
					Service: "mongo",
					Error:   err.Error(),
				})
			}

			return next(c)
		}
	}
}

// --- Redis Health Check ---

type RedisPingFunc func(ctx context.Context) error

type RedisHealthCheckConfig struct {
	PingFunc  RedisPingFunc
	Breaker   *gobreaker.TwoStepCircuitBreaker
	SkipPaths map[string]bool
	Timeout   time.Duration
	Endpoint  string
}

func RedisHealthCheckMiddleware(config RedisHealthCheckConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Register endpoint
			if c.Request().URL.Path == config.Endpoint {
				_, err := utils.ExecuteWithBreaker(config.Breaker, func() (interface{}, error) {
					ctx, cancel := context.WithTimeout(c.Request().Context(), config.Timeout)
					defer cancel()
					return nil, config.PingFunc(ctx)
				})
				if err != nil {
					return c.JSON(http.StatusServiceUnavailable, healthResponse{
						Status:  "fail",
						Service: "redis",
						Error:   err.Error(),
					})
				}
				return c.JSON(http.StatusOK, healthResponse{Status: "ok", Service: "redis"})
			}

			// Skip health check for specified paths
			if config.SkipPaths[c.Request().URL.Path] {
				return next(c)
			}

			// Check health on every request
			_, err := utils.ExecuteWithBreaker(config.Breaker, func() (interface{}, error) {
				ctx, cancel := context.WithTimeout(c.Request().Context(), config.Timeout)
				defer cancel()
				return nil, config.PingFunc(ctx)
			})
			if err != nil {
				c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c.Response().WriteHeader(http.StatusServiceUnavailable)
				return json.NewEncoder(c.Response().Writer).Encode(healthResponse{
					Status:  "fail",
					Service: "redis",
					Error:   err.Error(),
				})
			}

			return next(c)
		}
	}
}
