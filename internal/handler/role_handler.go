package handler

import (
	"net/http"

	"fdlp-standard-api/internal/dto"
	"fdlp-standard-api/internal/services"
	"fdlp-standard-api/internal/utils"

	"github.com/labstack/echo/v4"
)

type RoleHandler interface {
	GetRole(c echo.Context) error
	CreateRole(c echo.Context) error
}

type roleHandler struct {
	services services.RoleService
}

func NewRoleHandler(services services.RoleService) RoleHandler {
	return &roleHandler{services: services}
}

func (h *roleHandler) GetRole(c echo.Context) error {
	id := c.QueryParam("id")
	user, err := h.services.GetRole(id)
	if err != nil {
		return utils.JsonResponse(c, http.StatusNotFound, false, "Role not found")
	}
	return utils.JsonResponse(c, http.StatusOK, true, user)
}

func (h *roleHandler) CreateRole(c echo.Context) error {
	var request dto.CreateRoleRequestBody
	if err := c.Bind(&request); err != nil {
		return utils.JsonResponse(c, http.StatusBadRequest, false, "Invalid input")
	}
	if err := c.Validate(request); err != nil {
		return utils.JsonResponse(c, http.StatusBadRequest, false, err.Error())
	}

	role, err := h.services.CreateRole(request)
	if err != nil {
		return utils.JsonResponse(c, http.StatusInternalServerError, false, err.Error())
	}
	return utils.JsonResponse(c, http.StatusCreated, true, role)
}
