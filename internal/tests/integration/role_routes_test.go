package integration_test

import (
	"bytes"
	"fdlp-standard-api/internal/dto"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRoleRoutes(t *testing.T) {
	e, _ := SetupIntegrationServer()

	t.Run("CreateRole", func(t *testing.T) {
		roleReq := dto.CreateRoleRequestBody{
			Name:        "moderator",
			Description: "Forum Moderator",
		}
		roleBody, _ := json.Marshal(roleReq)
		req := httptest.NewRequest(http.MethodPost, "/role/create-role", bytes.NewReader(roleBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "moderator")
	})

	t.Run("GetRole", func(t *testing.T) {
		// Create a role first
		roleReq := dto.CreateRoleRequestBody{
			Name:        "editor",
			Description: "Content Editor",
		}
		roleBody, _ := json.Marshal(roleReq)
		reqCreate := httptest.NewRequest(http.MethodPost, "/role/create-role", bytes.NewReader(roleBody))
		reqCreate.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		recCreate := httptest.NewRecorder()
		e.ServeHTTP(recCreate, reqCreate)
		assert.Equal(t, http.StatusCreated, recCreate.Code)

		// Parse ID from response
		var createResp struct {
			Data struct {
				ID string `json:"ID"`
			} `json:"data"`
		}
		err := json.Unmarshal(recCreate.Body.Bytes(), &createResp)
		assert.NoError(t, err)
		roleID := createResp.Data.ID
		assert.NotEmpty(t, roleID)

		// Get Role
		reqGet := httptest.NewRequest(http.MethodGet, "/role?id="+roleID, nil)
		recGet := httptest.NewRecorder()
		e.ServeHTTP(recGet, reqGet)

		assert.Equal(t, http.StatusOK, recGet.Code)
		assert.Contains(t, recGet.Body.String(), "editor")
	})
}
