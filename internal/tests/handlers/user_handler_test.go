package handlers_test

import (
	"errors"
	"fdlp-standard-api/internal/dto"
	"fdlp-standard-api/internal/handler"
	"fdlp-standard-api/internal/models"
	mock_services "fdlp-standard-api/internal/tests/mock/services"
	"fdlp-standard-api/internal/types"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func TestUserHandler_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_services.NewMockUserService(ctrl)
	mockWsService := mock_services.NewMockWebSocketService(ctrl)

	h := handler.NewUserHandler(mockUserService, mockWsService)
	e := echo.New()

	t.Run("Success", func(t *testing.T) {
		userIDStr := "123e4567-e89b-12d3-a456-426614174000"
		req := httptest.NewRequest(http.MethodGet, "/users/"+userIDStr, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(userIDStr)

		expectedResponse := &dto.UserGetByIdResponse{
			Username:        "testuser",
			Email:           "test@example.com",
			RoleName:        "user",
			RoleDescription: "desc",
		}

		mockUserService.EXPECT().GetUser(userIDStr).Return(expectedResponse, nil)

		if assert.NoError(t, h.GetUser(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		userIDStr := "123e4567-e89b-12d3-a456-426614174000"
		req := httptest.NewRequest(http.MethodGet, "/users/"+userIDStr, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(userIDStr)

		mockUserService.EXPECT().GetUser(userIDStr).Return(nil, errors.New("not found"))

		if assert.NoError(t, h.GetUser(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})
}

func TestUserHandler_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_services.NewMockUserService(ctrl)
	mockWsService := mock_services.NewMockWebSocketService(ctrl)

	h := handler.NewUserHandler(mockUserService, mockWsService)

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	t.Run("Success", func(t *testing.T) {
		reqBody := `{"username":"newuser", "email":"new@example.com", "password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expectedUser := &models.User{
			Username: "newuser",
			Email:    "new@example.com",
			ID:       types.UUID(uuid.New()),
		}

		mockUserService.EXPECT().CreateUser(gomock.Any()).Return(expectedUser, nil)

		if assert.NoError(t, h.CreateUser(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})

	t.Run("InvalidInput_Bind", func(t *testing.T) {
		reqBody := `{"username":`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h.CreateUser(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("InvalidInput_Validate", func(t *testing.T) {
		// Missing fields
		reqBody := `{"username":""}`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h.CreateUser(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("ServiceError", func(t *testing.T) {
		reqBody := `{"username":"newuser", "email":"new@example.com", "password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expectedError := errors.New("service error")
		mockUserService.EXPECT().CreateUser(gomock.Any()).Return(nil, expectedError)

		if assert.NoError(t, h.CreateUser(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestUserHandler_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_services.NewMockUserService(ctrl)
	mockWsService := mock_services.NewMockWebSocketService(ctrl)

	h := handler.NewUserHandler(mockUserService, mockWsService)
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	t.Run("Success", func(t *testing.T) {
		reqBody := `{"email":"test@example.com", "password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		token := "mock_token_string"
		mockUserService.EXPECT().LoginUser(gomock.Any()).Return(&token, nil)
		mockWsService.EXPECT().BroadcastMessage(gomock.Any()).Return(nil)

		if assert.NoError(t, h.LoginUser(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		reqBody := `{"email":"test@example.com", "password":"wrong"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserService.EXPECT().LoginUser(gomock.Any()).Return(nil, errors.New("invalid"))

		if assert.NoError(t, h.LoginUser(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("InvalidInput_Bind", func(t *testing.T) {
		reqBody := `{"email":`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h.LoginUser(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("InvalidInput_Validate", func(t *testing.T) {
		reqBody := `{"email":""}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h.LoginUser(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}

func TestUserHandler_GetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_services.NewMockUserService(ctrl)
	mockWsService := mock_services.NewMockWebSocketService(ctrl)

	h := handler.NewUserHandler(mockUserService, mockWsService)
	e := echo.New()

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users?page=1&pageSize=10", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		users := []models.User{}
		mockUserService.EXPECT().GetUsersWithPagination(gomock.Any(), 1, 10).Return(users, int64(0), 0, nil)

		if assert.NoError(t, h.GetUsers(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Success_WithFilter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users?page=1&pageSize=10&filterBy=username&filterValue=test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		users := []models.User{}
		// Check that map contains username:test
		mockUserService.EXPECT().GetUsersWithPagination(gomock.Any(), 1, 10).DoAndReturn(func(filter map[string]interface{}, page, pageSize int) ([]models.User, int64, int, error) {
			assert.Equal(t, "test", filter["username"])
			return users, int64(0), 0, nil
		})

		if assert.NoError(t, h.GetUsers(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("ServiceError", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserService.EXPECT().GetUsersWithPagination(gomock.Any(), 0, 0).Return(nil, int64(0), 0, errors.New("db error"))

		if assert.NoError(t, h.GetUsers(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestUserHandler_CreateUserWithRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_services.NewMockUserService(ctrl)
	mockWsService := mock_services.NewMockWebSocketService(ctrl)

	h := handler.NewUserHandler(mockUserService, mockWsService)
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	t.Run("Success", func(t *testing.T) {
		reqBody := `{"username":"newuser", "email":"new@example.com", "password":"password123", "role_name":"admin"}`
		req := httptest.NewRequest(http.MethodPost, "/users/role", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expectedUser := &models.User{Username: "newuser"}
		mockUserService.EXPECT().CreateUserWithRole(gomock.Any()).Return(expectedUser, nil)

		if assert.NoError(t, h.CreateUserWithRole(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})

	t.Run("InvalidInput_Bind", func(t *testing.T) {
		reqBody := `{"username":`
		req := httptest.NewRequest(http.MethodPost, "/users/role", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h.CreateUserWithRole(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("InvalidInput_Validate", func(t *testing.T) {
		reqBody := `{"username":""}`
		req := httptest.NewRequest(http.MethodPost, "/users/role", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, h.CreateUserWithRole(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("ServiceError", func(t *testing.T) {
		reqBody := `{"username":"newuser", "email":"new@example.com", "password":"password123", "role_name":"admin"}`
		req := httptest.NewRequest(http.MethodPost, "/users/role", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserService.EXPECT().CreateUserWithRole(gomock.Any()).Return(nil, errors.New("error"))

		if assert.NoError(t, h.CreateUserWithRole(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}
