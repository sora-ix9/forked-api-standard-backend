package handler

import (
	"net/http"

	"fdlp-standard-api/internal/dto"
	"fdlp-standard-api/internal/services"
	"fdlp-standard-api/internal/utils"

	"github.com/labstack/echo/v4"
)

type PostHandler interface {
	GetPost(c echo.Context) error
	CreatePost(c echo.Context) error
}

type postHandler struct {
	services services.PostService
}

func NewPostHandler(services services.PostService) PostHandler {
	return &postHandler{services: services}
}

func (h *postHandler) GetPost(c echo.Context) error {
	id := c.QueryParam("id")
	post, err := h.services.GetPost(id)
	if err != nil {
		return utils.JsonResponse(c, http.StatusNotFound, false, "Post not found")
	}
	return utils.JsonResponse(c, http.StatusOK, true, post)
}

func (h *postHandler) CreatePost(c echo.Context) error {
	var request dto.CreatePostRequestBody
	if err := c.Bind(&request); err != nil {
		return utils.JsonResponse(c, http.StatusBadRequest, false, "Invalid input")
	}
	if err := c.Validate(request); err != nil {
		return utils.JsonResponse(c, http.StatusBadRequest, false, err.Error())
	}

	post, err := h.services.CreatePost(request)
	if err != nil {
		return utils.JsonResponse(c, http.StatusInternalServerError, false, err.Error())
	}
	return utils.JsonResponse(c, http.StatusCreated, true, post)
}
