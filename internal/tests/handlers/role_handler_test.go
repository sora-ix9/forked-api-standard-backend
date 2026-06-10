package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"fdlp-standard-api/internal/dto"
	"fdlp-standard-api/internal/handler"
	"fdlp-standard-api/internal/models"
	mock_services "fdlp-standard-api/internal/tests/mock/services"
	"fdlp-standard-api/internal/types"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type customValidator struct{ v *validator.Validate }

func (cv *customValidator) Validate(i interface{}) error { return cv.v.Struct(i) }

func newEcho() *echo.Echo {
	e := echo.New()
	e.Validator = &customValidator{v: validator.New()}
	return e
}

func TestCreateRole_Returns201OnSuccess(t *testing.T) {
	mock := &mock_services.RoleServiceMock{
		CreateRoleFn: func(r dto.CreateRoleRequestBody) (*models.Role, error) {
			return &models.Role{RoleID: types.UUID(uuid.New()), Name: r.Name}, nil
		},
	}
	e := newEcho()
	req := httptest.NewRequest(http.MethodPost, "/role/create-role", strings.NewReader(`{"name":"admin","description":"Administrator"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := handler.NewRoleHandler(mock)
	assert.NoError(t, h.CreateRole(c))
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestCreateRole_Returns400WhenNameMissing(t *testing.T) {
	mock := &mock_services.RoleServiceMock{}
	e := newEcho()
	req := httptest.NewRequest(http.MethodPost, "/role/create-role", strings.NewReader(`{"description":"no name"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := handler.NewRoleHandler(mock)
	assert.NoError(t, h.CreateRole(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateRole_Returns500OnServiceError(t *testing.T) {
	mock := &mock_services.RoleServiceMock{
		CreateRoleFn: func(r dto.CreateRoleRequestBody) (*models.Role, error) {
			return nil, errors.New("db error")
		},
	}
	e := newEcho()
	req := httptest.NewRequest(http.MethodPost, "/role/create-role", strings.NewReader(`{"name":"admin"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := handler.NewRoleHandler(mock)
	assert.NoError(t, h.CreateRole(c))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetRole_Returns200OnSuccess(t *testing.T) {
	id := uuid.New().String()
	mock := &mock_services.RoleServiceMock{
		GetRoleFn: func(i string) (*models.Role, error) {
			return &models.Role{Name: "admin"}, nil
		},
	}
	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/role?id="+id, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.QueryParams().Set("id", id)

	h := handler.NewRoleHandler(mock)
	assert.NoError(t, h.GetRole(c))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetRole_Returns404WhenNotFound(t *testing.T) {
	mock := &mock_services.RoleServiceMock{
		GetRoleFn: func(i string) (*models.Role, error) {
			return nil, errors.New("not found")
		},
	}
	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/role?id=nonexistent", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := handler.NewRoleHandler(mock)
	assert.NoError(t, h.GetRole(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
