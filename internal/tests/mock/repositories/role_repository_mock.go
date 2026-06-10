package mock_repositories

import "fdlp-standard-api/internal/models"

type RoleRepositoryMock struct {
	CreateFn    func(role *models.Role) error
	GetByIdFn   func(id string) (*models.Role, error)
	GetByNameFn func(name string) (*models.Role, error)
}

func (m *RoleRepositoryMock) Create(role *models.Role) error {
	return m.CreateFn(role)
}

func (m *RoleRepositoryMock) GetById(id string) (*models.Role, error) {
	return m.GetByIdFn(id)
}

func (m *RoleRepositoryMock) GetByName(name string) (*models.Role, error) {
	return m.GetByNameFn(name)
}
