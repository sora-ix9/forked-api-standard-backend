package echo

import (
	"context"
	"errors"
	"fdlp-standard-api/internal/handler"
	"fdlp-standard-api/internal/middlewares"
	"fdlp-standard-api/internal/repositories"
	"fdlp-standard-api/internal/routes"
	"fdlp-standard-api/internal/services"
	"fdlp-standard-api/pkg/config"
	"fdlp-standard-api/pkg/db"
	"fdlp-standard-api/pkg/utils"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func StartServer(cfg *config.Config, mongodb *db.MongoDB, redisClient *redis.Client) {
	e := newServer(cfg, mongodb, redisClient)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer(cfg *config.Config, mongodb *db.MongoDB, redisClient *redis.Client) *echo.Echo {
	// Echo
	e := echo.New()

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"}, // Allow frontend server
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
	}))

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339} ${method} ${uri} ${status} ${latency_human}` + "\n",
	}))
	e.Use(middleware.Recover())

	// Ignore Log payload
	skipLogPaths := map[string]bool{
		"/v1/users/update-image-profile": true,
	}
	e.Use(middlewares.RequestLogMiddleware(skipLogPaths))

	// Websocket Config
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Adjust this to ensure a proper origin check
		},
	}

	var database *mongo.Database
	if mongodb != nil {
		database = mongodb.DB
	}

	// Initialize Repo
	userRepo := repositories.NewUserRepository(database)
	roleRepo := repositories.NewRoleRepository(database)

	// Initialize Service
	userService := services.NewUserService(userRepo, roleRepo, cfg.JWTSecret, database)
	roleService := services.NewRoleService(roleRepo)
	websocketService := services.NewWebSocketService()

	// Initialize Handler
	userHandler := handler.NewUserHandler(userService, websocketService)
	roleHandler := handler.NewRoleHandler(roleService)
	websocketHandler := handler.NewWebsocketHandler(websocketService, upgrader)

	// Route Websocket
	e.GET("/ws", websocketHandler.WebSocketInit)

	// Global Breakers for health checks
	mongoBreaker := utils.NewCircuitBreaker("Health-Mongo")
	redisBreaker := utils.NewCircuitBreaker("Health-Redis")

	// Health check middlewares
	e.Use(middlewares.MongoHealthCheckMiddleware(middlewares.MongoHealthCheckConfig{
		PingFunc: mongoPingFunc(mongodb),
		Breaker:  mongoBreaker,
		SkipPaths: map[string]bool{
			"/":             true,
			"/ws":           true,
			"/health-mongo": true,
			"/health-redis": true,
		},
		Timeout:  2 * time.Second,
		Endpoint: "/health-mongo",
	}))

	e.Use(middlewares.RedisHealthCheckMiddleware(middlewares.RedisHealthCheckConfig{
		PingFunc: redisPingFunc(redisClient),
		Breaker:  redisBreaker,
		SkipPaths: map[string]bool{
			"/":             true,
			"/ws":           true,
			"/health-mongo": true,
			"/health-redis": true,
		},
		Timeout:  2 * time.Second,
		Endpoint: "/health-redis",
	}))
	// Documentation routes
	middlewares.RegisterDocRoutes(e)

	// Routes
	e.GET("/", helpCheckHandler)

	// Routes Authen
	r := e.Group("/v1")

	// Configure middleware with the custom claims type
	r.Use(middlewares.JWTMiddleware(cfg.JWTSecret))

	// Register Routes
	routes.RegisterUserRoutes(e, r, userHandler)
	routes.RegisterRoleRoutes(e, roleHandler)

	return e
}

func helpCheckHandler(c echo.Context) error {
	return c.String(http.StatusOK, "fdlp Standard API Status:Online")
}

func mongoPingFunc(mongodb *db.MongoDB) middlewares.MongoPingFunc {
	return func(ctx context.Context, rp *readpref.ReadPref) error {
		if mongodb == nil || mongodb.Client == nil {
			return errors.New("mongo client is not configured")
		}

		return mongodb.Client.Ping(ctx, rp)
	}
}

func redisPingFunc(redisClient *redis.Client) middlewares.RedisPingFunc {
	return func(ctx context.Context) error {
		if redisClient == nil {
			return errors.New("redis client is not configured")
		}

		return redisClient.Ping(ctx).Err()
	}
}
