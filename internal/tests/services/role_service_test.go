package services_test

import (
	"errors"
	"testing"
	"time"

	mock_repositories "fdlp-standard-api/internal/tests/mock/repositories"
	"fdlp-standard-api/internal/dto"
	"fdlp-standard-api/internal/models"
	"fdlp-standard-api/internal/services"
	"fdlp-standard-api/internal/types"

	"github.com/stretchr/testify/assert"
)

func newOKRoleMock() *mock_repositories.RoleRepositoryMock {
	return &mock_repositories.RoleRepositoryMock{
		CreateFn:    func(role *models.Role) error { return nil },
		GetByIdFn:   func(id string) (*models.Role, error) { return &models.Role{}, nil },
		GetByNameFn: func(name string) (*models.Role, error) { return &models.Role{}, nil },
	}
}

func TestCreateRole_SetsRoleID(t *testing.T) {
	svc := services.NewRoleService(newOKRoleMock())

	role, err := svc.CreateRole(dto.CreateRoleRequestBody{Name: "admin", Description: "Administrator"})

	assert.NoError(t, err)
	assert.NotEqual(t, types.UUID([16]byte{}), role.RoleID, "RoleID should not be zero")
}

func TestCreateRole_SetsTimestamps(t *testing.T) {
	before := time.Now().Add(-time.Second)
	svc := services.NewRoleService(newOKRoleMock())

	role, err := svc.CreateRole(dto.CreateRoleRequestBody{Name: "admin", Description: "Administrator"})

	assert.NoError(t, err)
	assert.True(t, role.CreatedAt.After(before), "CreatedAt should be set")
	assert.True(t, role.UpdatedAt.After(before), "UpdatedAt should be set")
}

func TestCreateRole_ReturnsErrorOnRepoFailure(t *testing.T) {
	mock := newOKRoleMock()
	mock.CreateFn = func(role *models.Role) error { return errors.New("db error") }

	svc := services.NewRoleService(mock)
	_, err := svc.CreateRole(dto.CreateRoleRequestBody{Name: "admin"})

	assert.Error(t, err)
}
