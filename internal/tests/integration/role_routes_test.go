package integration_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"fdlp-standard-api/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestIntegration_CreateRole_Returns201(t *testing.T) {
	e := newTestServer(newOKRoleRepoMock())

	req := httptest.NewRequest(http.MethodPost, "/role/create-role", strings.NewReader(`{"name":"admin","description":"Administrator"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestIntegration_CreateRole_Returns400WhenNameMissing(t *testing.T) {
	e := newTestServer(newOKRoleRepoMock())

	req := httptest.NewRequest(http.MethodPost, "/role/create-role", strings.NewReader(`{"description":"no name"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestIntegration_GetRole_Returns200(t *testing.T) {
	e := newTestServer(newOKRoleRepoMock())

	req := httptest.NewRequest(http.MethodGet, "/role?id=some-id", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestIntegration_GetRole_Returns404WhenNotFound(t *testing.T) {
	mock := newOKRoleRepoMock()
	mock.GetByIdFn = func(id string) (*models.Role, error) {
		return nil, errors.New("not found")
	}
	e := newTestServer(mock)

	req := httptest.NewRequest(http.MethodGet, "/role?id=nonexistent", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
