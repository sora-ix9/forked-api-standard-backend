package echo

import (
	"fdlp-standard-api/internal/dto"
	"fdlp-standard-api/internal/handler"
	"fdlp-standard-api/internal/middlewares"
	"fdlp-standard-api/internal/repositories"
	"fdlp-standard-api/internal/routes"
	"fdlp-standard-api/internal/services"
	"fdlp-standard-api/internal/utils"
	"fdlp-standard-api/pkg/config"
	"fdlp-standard-api/pkg/db"
	"net/http"
	"os"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/watchakorn-18k/scalar-go"

	"log"

	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// InitServer initializes and starts the Echo web server
func InitServer() {
	// Loading .env from the root directory
	err := godotenv.Load(".env") // Adjust the path according to your project structure
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	utils.SetGlobalTimezone()
	log.Println("Global timezone set to Asia/Bangkok (UTC+07:00)")

	// Initialize Configuration
	cfg := config.New()

	// Echo
	e := echo.New()

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"}, // Allow frontend server
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	// Middleware
	// e.Use(middleware.Logger())
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

	// Initialize Provider

	// Initialize Database with Config
	database := db.InitializeDB(cfg)
	defer db.CloseDB(database)

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

	// Routes
	e.GET("/", helpCheckHandler)

	// Documentation routes
	statusMode := os.Getenv("BACKEND_MODE")
	if statusMode == "" {
		statusMode = "production"
	}

	if statusMode == "development" {
		e.GET("/docs", func(c echo.Context) error {
			htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
				SpecURL: "./docs/swagger.yaml",
				CustomOptions: scalar.CustomOptions{
					PageTitle: "PINTO API SPEC",
				},
				DarkMode: true,
			})

			if err != nil {
				return err
			}
			c.Response().Header().Set("Content-Type", "text/html")
			return c.HTML(http.StatusOK, htmlContent)
		})
		e.GET("/swagger.yaml", func(c echo.Context) error {
			return c.File("./docs/swagger.yaml")
		})
		// Documentation routes
		e.Static("/docs", "docs")
		e.GET("/schema", func(c echo.Context) error {
			htmlContent, err := os.ReadFile("./docs/database-schema.html")
			if err != nil {
				return err
			}
			c.Response().Header().Set("Content-Type", "text/html")
			return c.HTML(http.StatusOK, string(htmlContent))
		})
	}
	// Routes Authen
	r := e.Group("/v1")

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(dto.JwtCustomClaims)
		},
		SigningKey: []byte(cfg.JWTSecret),
	}
	r.Use(echojwt.WithConfig(config))

	// Register Routes
	routes.RegisterUserRoutes(e, r, userHandler)
	routes.RegisterRoleRoutes(e, roleHandler)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func helpCheckHandler(c echo.Context) error {
	return c.String(http.StatusOK, "fdlp Standard API Status:Online")
}
