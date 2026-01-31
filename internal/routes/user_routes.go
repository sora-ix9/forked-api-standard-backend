package routes

import (
	"fdlp-standard-api/internal/handler"
	"fdlp-standard-api/internal/middlewares"

	"github.com/labstack/echo/v4"
)

func RegisterUserRoutes(e *echo.Echo, r *echo.Group, userHandler handler.UserHandler) {
	// Public Routes
	e.POST("/user/create-user", userHandler.CreateUser)
	e.POST("/user/login", userHandler.LoginUser)

	// Protected Routes (Authentication Required)
	// Apply RoleMiddleware if needed to specific routes
	r.GET("/user/:id", userHandler.GetUser, middlewares.RoleMiddleware("admin"))
	r.GET("/users", userHandler.GetUsers, middlewares.RoleMiddleware("admin"))
	r.POST("/user/create", userHandler.CreateUserWithRole, middlewares.RoleMiddleware("admin"))
}
