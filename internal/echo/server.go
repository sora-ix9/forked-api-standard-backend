package echo

import (
	"context"
	"fdlp-standard-api/internal/handler"
	"fdlp-standard-api/internal/middlewares"
	"fdlp-standard-api/internal/repositories"
	"fdlp-standard-api/internal/routes"
	"fdlp-standard-api/internal/services"
	"fdlp-standard-api/pkg/config"
	"fdlp-standard-api/pkg/providers/mongodb"
	"fdlp-standard-api/pkg/redisclient"
	"fdlp-standard-api/pkg/utils"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// InitServer initializes and starts the Echo web server
func InitServer(cfg *config.Config) {
	// Echo
	e := echo.New()

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"}, // Allow frontend server
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
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

	// Initialize MongoDB
	mongoClient := mongodb.NewClient(cfg)
	database := mongoClient.DB

	// Initialize Redis
	redisClient := redisclient.NewClient(cfg)

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

	// Routes
	e.GET("/", helpCheckHandler)
	e.GET("/health-mongo", func(c echo.Context) error {
		_, err := utils.ExecuteWithBreaker(mongoBreaker, func() (interface{}, error) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			return nil, mongoClient.Client.Ping(ctx, readpref.Primary())
		})

		if err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "fail", "error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	e.GET("/health-redis", func(c echo.Context) error {
		_, err := utils.ExecuteWithBreaker(redisBreaker, func() (interface{}, error) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			return nil, redisClient.Ping(ctx).Err()
		})

		if err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "fail", "error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Documentation routes
	middlewares.RegisterDocRoutes(e)
	// Routes Authen
	r := e.Group("/v1")

	// Configure middleware with the custom claims type
	r.Use(middlewares.JWTMiddleware(cfg.JWTSecret))

	// Register Routes
	routes.RegisterUserRoutes(e, r, userHandler)
	routes.RegisterRoleRoutes(e, roleHandler)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func helpCheckHandler(c echo.Context) error {
	return c.String(http.StatusOK, "fdlp Standard API Status:Online")
}
