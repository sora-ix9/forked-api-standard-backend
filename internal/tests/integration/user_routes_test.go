package integration_test

import (
	"bytes"
	"fdlp-standard-api/internal/dto"
	"fdlp-standard-api/internal/models"
	"fdlp-standard-api/internal/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserRoutes(t *testing.T) {
	e, db := SetupIntegrationServer()

	// 1. Seed 'admin' Role and User
	adminRoleID := uuid.New()
	adminRole := models.Role{
		ID:          types.UUID(adminRoleID),
		Name:        "admin",
		Description: "Administrator",
	}
	if err := db.Create(&adminRole).Error; err != nil {
		t.Fatalf("Failed to seed role: %v", err)
	}

	pwd, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	adminUser := models.User{
		ID:       types.UUID(uuid.New()),
		Username: "superadmin",
		Email:    "admin@test.com",
		Password: string(pwd),
		RoleID:   types.UUID(adminRoleID),
	}
	if err := db.Create(&adminUser).Error; err != nil {
		t.Fatalf("Failed to seed user: %v", err)
	}

	t.Run("LoginAndAccessProtectedResource", func(t *testing.T) {
		// Login
		loginReq := dto.LoginUserRequestBody{
			Email:    "admin@test.com",
			Password: "password123",
		}
		loginBody, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewReader(loginBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if !assert.Equal(t, http.StatusOK, rec.Code) {
			t.Logf("Login Response: %s", rec.Body.String())
		}

		var loginResp struct {
			Ok   bool   `json:"ok"`
			Data string `json:"data"` // Token
		}
		err := json.Unmarshal(rec.Body.Bytes(), &loginResp)
		assert.NoError(t, err)
		assert.True(t, loginResp.Ok)
		token := loginResp.Data
		assert.NotEmpty(t, token)

		// Access Protected Route (GET /v1/users)
		reqProtected := httptest.NewRequest(http.MethodGet, "/v1/users?page=1&pageSize=10", nil)
		reqProtected.Header.Set("Authorization", "Bearer "+token)
		recProtected := httptest.NewRecorder()
		e.ServeHTTP(recProtected, reqProtected)

		if !assert.Equal(t, http.StatusOK, recProtected.Code) {
			t.Logf("Protected Response: %s", recProtected.Body.String())
		}
		assert.Contains(t, recProtected.Body.String(), "superadmin")
	})

	t.Run("CreateUser_Public", func(t *testing.T) {
		// Need 'guest' role for public create user
		guestRole := models.Role{
			Name: "guest",
		}
		db.Create(&guestRole)

		userReq := dto.CreateUserRequestBody{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
		}
		userBody, _ := json.Marshal(userReq)
		req := httptest.NewRequest(http.MethodPost, "/user/create-user", bytes.NewReader(userBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "newuser")
	})
}
