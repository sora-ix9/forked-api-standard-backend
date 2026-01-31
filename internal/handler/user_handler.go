package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"fdlp-standard-api/internal/dto"
	"fdlp-standard-api/internal/services"
	"fdlp-standard-api/internal/utils"

	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	GetUser(c echo.Context) error
	GetUsers(c echo.Context) error
	CreateUser(c echo.Context) error
	LoginUser(c echo.Context) error
	CreateUserWithRole(c echo.Context) error
}

type userHandler struct {
	services  services.UserService
	wsService services.WebSocketService
}

func NewUserHandler(services services.UserService, wsService services.WebSocketService) UserHandler {
	return &userHandler{services: services, wsService: wsService}
}

func (h *userHandler) GetUser(c echo.Context) error {
	id := c.Param("id")
	user, err := h.services.GetUser(id)
	if err != nil {
		return utils.JsonResponse(c, http.StatusNotFound, false, "User not found")
	}
	return utils.JsonResponse(c, http.StatusOK, true, user)
}

func (h *userHandler) CreateUser(c echo.Context) error {
	var request dto.CreateUserRequestBody
	if err := c.Bind(&request); err != nil {
		return utils.JsonResponse(c, http.StatusBadRequest, false, "Invalid input")
	}
	if err := c.Validate(request); err != nil {
		return utils.JsonResponse(c, http.StatusBadRequest, false, err.Error())
	}

	user, err := h.services.CreateUser(request)
	if err != nil {
		return utils.JsonResponse(c, http.StatusInternalServerError, false, err.Error())
	}
	return utils.JsonResponse(c, http.StatusCreated, true, user)
}

func (h *userHandler) CreateUserWithRole(c echo.Context) error {
	var request dto.CreateUserWithRoleRequestBody
	if err := c.Bind(&request); err != nil {
		return utils.JsonResponse(c, http.StatusBadRequest, false, "Invalid input")
	}
	if err := c.Validate(request); err != nil {
		return utils.JsonResponse(c, http.StatusBadRequest, false, err.Error())
	}

	user, err := h.services.CreateUserWithRole(request)
	if err != nil {
		return utils.JsonResponse(c, http.StatusInternalServerError, false, err.Error())
	}
	return utils.JsonResponse(c, http.StatusCreated, true, user)
}

func (h *userHandler) LoginUser(c echo.Context) error {
	var request dto.LoginUserRequestBody
	if err := c.Bind(&request); err != nil {
		return utils.JsonResponse(c, http.StatusBadRequest, false, "Invalid input")
	}
	if err := c.Validate(request); err != nil {
		return utils.JsonResponse(c, http.StatusBadRequest, false, err.Error())
	}

	// Authentication logic goes here
	token, err := h.services.LoginUser(request)
	if err != nil {
		return utils.JsonResponse(c, http.StatusInternalServerError, false, err.Error())
	}

	h.wsService.BroadcastMessage(fmt.Sprintf("User:%s logged in.", request.Email))

	return utils.JsonResponse(c, http.StatusOK, true, token)
}

func (h *userHandler) GetUsers(c echo.Context) error {
	page := c.QueryParam("page")
	pageSize := c.QueryParam("pageSize")
	filterBy := c.QueryParam("filterBy")
	filterValue := c.QueryParam("filterValue")
	filter := make(map[string]interface{})

	if filterBy != "" && filterValue != "" {
		filter[filterBy] = filterValue
	}

	// Convert pagination params to integers
	pageNum, _ := strconv.Atoi(page)
	pageSizeNum, _ := strconv.Atoi(pageSize)

	users, totalRows, totalPages, err := h.services.GetUsersWithPagination(filter, pageNum, pageSizeNum)
	if err != nil {
		return utils.JsonResponse(c, http.StatusInternalServerError, false, "Failed to fetch users: "+err.Error())
	}

	return utils.JsonResponse(c, http.StatusOK, true, echo.Map{
		"users":      users,
		"totalRows":  totalRows,
		"totalPages": totalPages,
	})
}
