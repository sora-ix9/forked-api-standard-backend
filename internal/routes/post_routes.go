package routes

import (
	"fdlp-standard-api/internal/handler"

	"github.com/labstack/echo/v4"
)

func RegisterPostRoutes(e *echo.Echo, postHandler handler.PostHandler) {
	e.GET("/post", postHandler.GetPost)
	e.POST("/post/create-post", postHandler.CreatePost)
}
