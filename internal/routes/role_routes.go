package routes

import (
	"fdlp-standard-api/internal/handler"

	"github.com/labstack/echo/v4"
)

func RegisterRoleRoutes(e *echo.Echo, roleHandler handler.RoleHandler) {
	e.GET("/role", roleHandler.GetRole)
	e.POST("/role/create-role", roleHandler.CreateRole)
}
