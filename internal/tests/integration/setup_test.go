package integration_test

import (
	"fdlp-standard-api/internal/handler"
	"fdlp-standard-api/internal/repositories"
	"fdlp-standard-api/internal/routes"
	"fdlp-standard-api/internal/services"
	mock_repositories "fdlp-standard-api/internal/tests/mock/repositories"
	"fdlp-standard-api/internal/models"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type customValidator struct{ v *validator.Validate }

func (cv *customValidator) Validate(i interface{}) error { return cv.v.Struct(i) }

func newTestServer(roleRepo repositories.RoleRepository) *echo.Echo {
	e := echo.New()
	e.Validator = &customValidator{v: validator.New()}
	e.Use(middleware.Recover())

	roleSvc := services.NewRoleService(roleRepo)
	roleHandler := handler.NewRoleHandler(roleSvc)
	routes.RegisterRoleRoutes(e, roleHandler)

	return e
}

func newOKRoleRepoMock() *mock_repositories.RoleRepositoryMock {
	return &mock_repositories.RoleRepositoryMock{
		CreateFn:    func(role *models.Role) error { return nil },
		GetByIdFn:   func(id string) (*models.Role, error) { return &models.Role{Name: "admin"}, nil },
		GetByNameFn: func(name string) (*models.Role, error) { return &models.Role{Name: name}, nil },
	}
}
