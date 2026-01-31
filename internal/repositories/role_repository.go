package repositories

import (
	"fdlp-standard-api/internal/models"

	"gorm.io/gorm"
)

type RoleRepository interface {
	GetById(id string) (*models.Role, error)
	GetByName(name string) (*models.Role, error)
	Create(user *models.Role) error
	CreateWithTransaction(tx *gorm.DB, role *models.Role) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (repo *roleRepository) GetById(id string) (*models.Role, error) {
	var role models.Role
	result := repo.db.Where("id = ?", id).First(&role)
	return &role, result.Error
}

func (repo *roleRepository) GetByName(name string) (*models.Role, error) {
	var role models.Role
	result := repo.db.Where(&models.Role{Name: name}).First(&role)
	return &role, result.Error
}

func (repo *roleRepository) Create(role *models.Role) error {
	return repo.db.Create(role).Error
}

func (repo *roleRepository) CreateWithTransaction(tx *gorm.DB, role *models.Role) error {
	return tx.Create(role).Error
}
