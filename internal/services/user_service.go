package services

import (
	"fdlp-standard-api/internal/dto"
	"fdlp-standard-api/internal/models"
	"fdlp-standard-api/internal/repositories"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	GetUser(id string) (*dto.UserGetByIdResponse, error)
	CreateUser(user dto.CreateUserRequestBody) (*models.User, error)
	CreateUserWithRole(userWithRole dto.CreateUserWithRoleRequestBody) (*models.User, error)
	LoginUser(loginInfo dto.LoginUserRequestBody) (*string, error)
	GetUsersWithPagination(filter map[string]interface{}, page, pageSize int) ([]models.User, int64, int, error)
}

type userService struct {
	repo      repositories.UserRepository
	repoRole  repositories.RoleRepository
	jwtSecret string
	db        *gorm.DB
}

func NewUserService(repo repositories.UserRepository, repoRole repositories.RoleRepository, jwtSecret string, db *gorm.DB) UserService {
	return &userService{repo: repo, repoRole: repoRole, jwtSecret: jwtSecret, db: db}
}

func (s *userService) GetUser(id string) (*dto.UserGetByIdResponse, error) {
	user, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}

	result := dto.UserGetByIdResponse{
		Username:        user.Username,
		Email:           user.Email,
		RoleName:        user.Role.Name,
		RoleDescription: user.Role.Description,
	}
	return &result, nil
}

func (s *userService) CreateUser(data dto.CreateUserRequestBody) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err // Handle password encryption failure
	}

	// Get Role
	role, errRole := s.repoRole.GetByName("guest")
	if errRole != nil {
		return nil, errRole
	}

	user := models.User{
		Username: data.Username,
		Email:    data.Email,
		Password: string(hashedPassword),
		RoleID:   role.ID,
	}

	// Assuming you have a repositories method to save the User model
	if err := s.repo.Create(&user); err != nil {
		return nil, err // Handle user creation failure
	}

	return &user, nil
}

func (s *userService) LoginUser(data dto.LoginUserRequestBody) (*string, error) {
	// Get User
	user, err := s.repo.GetByEmail(data.Email)
	if err != nil {
		return nil, err
	}

	// Compare password
	byteHash := []byte(user.Password)
	plainPwdByte := []byte(data.Password)

	// Compare the stored hashed password, with the hashed version of the password that was received
	err = bcrypt.CompareHashAndPassword(byteHash, plainPwdByte)
	if err != nil {
		fmt.Println("Passwords do not match")
		return nil, err
	}

	// Set custom claims
	claims := &dto.JwtCustomClaims{
		Name:     user.Username,
		Role:     user.Role.Name,
		RoleName: user.Role.Description,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response
	secret := []byte(s.jwtSecret)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		fmt.Println("Sign invalid")
		return nil, err
	}

	return &signedToken, nil
}

func (s *userService) CreateUserWithRole(data dto.CreateUserWithRoleRequestBody) (*models.User, error) {
	var user models.User

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	errTransaction := s.db.Transaction(func(tx *gorm.DB) error {
		role, errRole := s.repoRole.GetByName(data.RoleName)
		if role == nil || errRole != nil {
			newRole := models.Role{Name: data.RoleName}
			if err := s.repoRole.CreateWithTransaction(tx, &newRole); err != nil {
				return err
			}
			role = &newRole
		}

		user = models.User{
			Username: data.Username,
			Email:    data.Email,
			Password: string(hashedPassword),
			RoleID:   role.ID,
		}
		return s.repo.CreateWithTransaction(tx, &user)
	})

	if errTransaction != nil {
		return nil, errTransaction
	}

	return &user, nil
}

func (s *userService) GetUsersWithPagination(filter map[string]interface{}, page, pageSize int) ([]models.User, int64, int, error) {
	return s.repo.FindAllByFilterAndPage(filter, page, pageSize)
}
