package middlewares

import (
	"bytes"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func RequestLogMiddleware(skipLogPaths map[string]bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// Check if we should skip body logging (essential for file uploads to allow streaming)
			if skipLogPaths[c.Path()] {
				logrus.WithFields(logrus.Fields{
					"method": c.Request().Method,
					"path":   c.Path(),
					"params": c.QueryParams(),
					"info":   "Body logging skipped for performance",
				}).Info("Handled request")
				return next(c)
			}

			// Read body for logging (only for non-skipped paths)
			bodyBytes, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return err
			}
			// Restore the body for further use
			c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// Log body and params
			logrus.WithFields(logrus.Fields{
				"method": c.Request().Method,
				"path":   c.Path(),
				"params": c.QueryParams(),
				"body":   string(bodyBytes),
			}).Info("Handled request")

			return next(c)
		}
	}
}
