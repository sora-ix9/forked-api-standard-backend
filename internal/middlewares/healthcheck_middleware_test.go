package middlewares

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoPingRecorder struct {
	calls int
	err   error
}

func (r *mongoPingRecorder) ping(ctx context.Context, rp *readpref.ReadPref) error {
	r.calls++
	if ctx == nil {
		return errors.New("missing context")
	}
	if rp == nil {
		return errors.New("missing read preference")
	}

	return r.err
}

type redisPingRecorder struct {
	calls int
	err   error
}

func (r *redisPingRecorder) ping(ctx context.Context) error {
	r.calls++
	if ctx == nil {
		return errors.New("missing context")
	}

	return r.err
}

func TestMongoHealthCheckMiddleware(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		skipPaths   map[string]bool
		pingErr     error
		wantStatus  int
		wantBody    string
		wantCalls   int
		wantNextRun bool
	}{
		{
			name:       "endpoint healthy",
			path:       "/health-mongo",
			wantStatus: http.StatusOK,
			wantBody:   `"service":"mongo"`,
			wantCalls:  1,
		},
		{
			name:       "endpoint unhealthy",
			path:       "/health-mongo",
			pingErr:    errors.New("mongo down"),
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   "mongo down",
			wantCalls:  1,
		},
		{
			name:        "request healthy",
			path:        "/v1/users",
			wantStatus:  http.StatusOK,
			wantBody:    "ok",
			wantCalls:   1,
			wantNextRun: true,
		},
		{
			name:       "request unhealthy",
			path:       "/v1/users",
			pingErr:    errors.New("mongo down"),
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   `"service":"mongo"`,
			wantCalls:  1,
		},
		{
			name:        "skip path",
			path:        "/",
			skipPaths:   map[string]bool{"/": true},
			wantStatus:  http.StatusOK,
			wantBody:    "ok",
			wantNextRun: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			pinger := &mongoPingRecorder{err: tt.pingErr}
			nextRun := false

			middleware := MongoHealthCheckMiddleware(MongoHealthCheckConfig{
				PingFunc:  pinger.ping,
				Breaker:   gobreaker.NewTwoStepCircuitBreaker(gobreaker.Settings{Name: "Test-Mongo-" + tt.name}),
				SkipPaths: tt.skipPaths,
				Timeout:   time.Second,
				Endpoint:  "/health-mongo",
			})

			err := middleware(func(c echo.Context) error {
				nextRun = true
				return c.String(http.StatusOK, "ok")
			})(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.wantBody)
			assert.Equal(t, tt.wantCalls, pinger.calls)
			assert.Equal(t, tt.wantNextRun, nextRun)
		})
	}
}

func TestRedisHealthCheckMiddleware(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		skipPaths   map[string]bool
		pingErr     error
		wantStatus  int
		wantBody    string
		wantCalls   int
		wantNextRun bool
	}{
		{
			name:       "endpoint healthy",
			path:       "/health-redis",
			wantStatus: http.StatusOK,
			wantBody:   `"service":"redis"`,
			wantCalls:  1,
		},
		{
			name:       "endpoint unhealthy",
			path:       "/health-redis",
			pingErr:    errors.New("redis down"),
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   "redis down",
			wantCalls:  1,
		},
		{
			name:        "request healthy",
			path:        "/v1/users",
			wantStatus:  http.StatusOK,
			wantBody:    "ok",
			wantCalls:   1,
			wantNextRun: true,
		},
		{
			name:       "request unhealthy",
			path:       "/v1/users",
			pingErr:    errors.New("redis down"),
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   `"service":"redis"`,
			wantCalls:  1,
		},
		{
			name:        "skip path",
			path:        "/",
			skipPaths:   map[string]bool{"/": true},
			wantStatus:  http.StatusOK,
			wantBody:    "ok",
			wantNextRun: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			pinger := &redisPingRecorder{err: tt.pingErr}
			nextRun := false

			middleware := RedisHealthCheckMiddleware(RedisHealthCheckConfig{
				PingFunc:  pinger.ping,
				Breaker:   gobreaker.NewTwoStepCircuitBreaker(gobreaker.Settings{Name: "Test-Redis-" + tt.name}),
				SkipPaths: tt.skipPaths,
				Timeout:   time.Second,
				Endpoint:  "/health-redis",
			})

			err := middleware(func(c echo.Context) error {
				nextRun = true
				return c.String(http.StatusOK, "ok")
			})(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.wantBody)
			assert.Equal(t, tt.wantCalls, pinger.calls)
			assert.Equal(t, tt.wantNextRun, nextRun)
		})
	}
}
