package repositories

import (
	"fdlp-standard-api/internal/models"
	"math"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetById(id string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	CreateWithTransaction(tx *gorm.DB, user *models.User) error
	FindAllByFilterAndPage(filter map[string]interface{}, page, pageSize int) ([]models.User, int64, int, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (repo *userRepository) GetById(id string) (*models.User, error) {
	var user models.User
	result := repo.db.Preload("Role").First(&user, id)
	return &user, result.Error
}

func (repo *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := repo.db.Where(&models.User{Email: email}).Preload("Role").First(&user)
	return &user, result.Error
}

func (repo *userRepository) Create(user *models.User) error {
	return repo.db.Create(user).Error
}

func (repo *userRepository) Update(user *models.User) error {
	return repo.db.Model(models.User{}).Where("id = ?", user.ID).Updates(user).Error
}

func (repo *userRepository) CreateWithTransaction(tx *gorm.DB, user *models.User) error {
	if len(user.RoleID) > 0 {
		var role models.Role
		if err := tx.Model(&models.Role{}).Where("id = ?", user.RoleID).First(&role).Error; err != nil {
			return err
		}
		user.Role = role
	}
	return tx.Create(user).Error
}

func (repo *userRepository) FindAllByFilterAndPage(filter map[string]interface{}, page, pageSize int) ([]models.User, int64, int, error) {
	var users []models.User
	var totalRows int64

	query := repo.db.Model(&models.User{})
	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	// Count total results
	err := query.Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, err
	}

	// Retrieve paginated results
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Preload("Role").Find(&users).Error
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(pageSize)))
	return users, totalRows, totalPages, nil
}
