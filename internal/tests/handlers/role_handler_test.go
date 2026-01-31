package handlers_test

import (
	"errors"
	"fdlp-standard-api/internal/handler"
	"fdlp-standard-api/internal/models"
	mock_services "fdlp-standard-api/internal/tests/mock/services"
	"fdlp-standard-api/internal/types"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/goccy/go-json"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRoleHandler_GetRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_services.NewMockRoleService(ctrl)
	h := handler.NewRoleHandler(mockService)

	e := echo.New()

	t.Run("Success", func(t *testing.T) {
		roleIDStr := "123e4567-e89b-12d3-a456-426614174000"
		// Use query param
		req := httptest.NewRequest(http.MethodGet, "/role?id="+roleIDStr, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expectedRole := &models.Role{
			ID:          types.UUID(uuid.MustParse(roleIDStr)),
			Name:        "admin",
			Description: "Administrator",
		}

		mockService.EXPECT().GetRole(roleIDStr).Return(expectedRole, nil)

		if assert.NoError(t, h.GetRole(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, true, resp["ok"]) // Changed from success to ok

			data := resp["data"].(map[string]interface{})
			assert.Equal(t, expectedRole.Name, data["Name"])
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		roleIDStr := "123e4567-e89b-12d3-a456-426614174000"
		req := httptest.NewRequest(http.MethodGet, "/role?id="+roleIDStr, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().GetRole(roleIDStr).Return(nil, errors.New("not found"))

		if assert.NoError(t, h.GetRole(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})
}

func TestRoleHandler_CreateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_services.NewMockRoleService(ctrl)
	h := handler.NewRoleHandler(mockService)

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	t.Run("Success", func(t *testing.T) {
		reqBody := `{"name":"moderator", "description":"Forum Moderator"}`
		req := httptest.NewRequest(http.MethodPost, "/role", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expectedRole := &models.Role{
			Name:        "moderator",
			Description: "Forum Moderator",
			ID:          types.UUID(uuid.New()),
		}

		mockService.EXPECT().CreateRole(gomock.Any()).Return(expectedRole, nil)

		if assert.NoError(t, h.CreateRole(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})

	t.Run("InvalidInput_Bind", func(t *testing.T) {
		reqBody := `{"name":` // Invalid JSON
		req := httptest.NewRequest(http.MethodPost, "/role", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h.CreateRole(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("InvalidInput_Validate", func(t *testing.T) {
		reqBody := `{"name":""}` // Empty name
		req := httptest.NewRequest(http.MethodPost, "/role", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h.CreateRole(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("ServiceError", func(t *testing.T) {
		reqBody := `{"name":"moderator", "description":"Forum Moderator"}`
		req := httptest.NewRequest(http.MethodPost, "/role", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.EXPECT().CreateRole(gomock.Any()).Return(nil, errors.New("db error"))

		if assert.NoError(t, h.CreateRole(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}
