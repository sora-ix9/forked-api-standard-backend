package services

import (
	"fdlp-standard-api/internal/dto"
	"fdlp-standard-api/internal/models"
	"fdlp-standard-api/internal/repositories"
)

type RoleService interface {
	GetRole(id string) (*models.Role, error)
	CreateRole(role dto.CreateRoleRequestBody) (*models.Role, error)
}

type roleService struct {
	repo repositories.RoleRepository
}

func NewRoleService(repo repositories.RoleRepository) RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) GetRole(id string) (*models.Role, error) {
	return s.repo.GetById(id)
}

func (s *roleService) CreateRole(data dto.CreateRoleRequestBody) (*models.Role, error) {
	role := models.Role{
		Name:        data.Name,
		Description: data.Description,
	}

	if err := s.repo.Create(&role); err != nil {
		return nil, err
	}

	return &role, nil
}
