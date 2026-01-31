package utils

import (
	"fdlp-standard-api/internal/dto"

	"github.com/labstack/echo/v4"
)

// JsonResponse sends a standardized JSON response with a message and optional data.
func JsonResponse(c echo.Context, statusCode int, ok bool, data interface{}) error {
	resp := dto.StandardReponseBody{
		Ok: ok,
	}

	if ok {
		resp.Data = data
	} else {
		resp.Error = data
	}

	return c.JSON(statusCode, resp)
}
