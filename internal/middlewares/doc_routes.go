package middlewares

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/watchakorn-18k/scalar-go"
)

// RegisterDocRoutes registers documentation routes if the application is in development mode
func RegisterDocRoutes(e *echo.Echo) {
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

		// Documentation static files
		e.Static("/docs-static", "docs")

		e.GET("/schema", func(c echo.Context) error {
			htmlContent, err := os.ReadFile("./docs/database-schema.html")
			if err != nil {
				return err
			}
			c.Response().Header().Set("Content-Type", "text/html")
			return c.HTML(http.StatusOK, string(htmlContent))
		})
	}
}
