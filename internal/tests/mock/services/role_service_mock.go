package mock_services

import (
	"fdlp-standard-api/internal/dto"
	"fdlp-standard-api/internal/models"
)

type RoleServiceMock struct {
	GetRoleFn    func(id string) (*models.Role, error)
	CreateRoleFn func(role dto.CreateRoleRequestBody) (*models.Role, error)
}

func (m *RoleServiceMock) GetRole(id string) (*models.Role, error) {
	return m.GetRoleFn(id)
}

func (m *RoleServiceMock) CreateRole(role dto.CreateRoleRequestBody) (*models.Role, error) {
	return m.CreateRoleFn(role)
}
