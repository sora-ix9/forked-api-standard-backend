package services_test

import (
	"errors"
	"fdlp-standard-api/internal/dto"
	"fdlp-standard-api/internal/models"
	"fdlp-standard-api/internal/services"
	mock_repositories "fdlp-standard-api/internal/tests/mock/repositories"
	"fdlp-standard-api/internal/types"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserService_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repositories.NewMockUserRepository(ctrl)
	mockRoleRepo := mock_repositories.NewMockRoleRepository(ctrl)

	svc := services.NewUserService(mockUserRepo, mockRoleRepo, "secret", nil)

	t.Run("Success", func(t *testing.T) {
		userIDStr := "123e4567-e89b-12d3-a456-426614174000"
		userID := types.UUID(uuid.MustParse(userIDStr))

		expectedUser := &models.User{
			ID:       userID,
			Username: "testuser",
			Email:    "test@example.com",
			Role: models.Role{
				Name:        "user",
				Description: "Standard User",
			},
		}

		mockUserRepo.EXPECT().GetById(userIDStr).Return(expectedUser, nil)

		result, err := svc.GetUser(userIDStr)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedUser.Username, result.Username)
		assert.Equal(t, expectedUser.Email, result.Email)
		assert.Equal(t, expectedUser.Role.Name, result.RoleName)
		assert.Equal(t, expectedUser.Role.Description, result.RoleDescription)
	})

	t.Run("Error", func(t *testing.T) {
		userIDStr := "123e4567-e89b-12d3-a456-426614174000"
		expectedError := errors.New("user not found")

		mockUserRepo.EXPECT().GetById(userIDStr).Return(nil, expectedError)

		result, err := svc.GetUser(userIDStr)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})
}

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repositories.NewMockUserRepository(ctrl)
	mockRoleRepo := mock_repositories.NewMockRoleRepository(ctrl)

	// In create user, jwt secret and db is not used directly in this method flow logic (except transaction in CreateUserWithRole)
	svc := services.NewUserService(mockUserRepo, mockRoleRepo, "secret", nil)

	t.Run("Success", func(t *testing.T) {
		req := dto.CreateUserRequestBody{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
		}

		roleIDStr := "123e4567-e89b-12d3-a456-426614174001"
		guestRole := &models.Role{
			ID:   types.UUID(uuid.MustParse(roleIDStr)),
			Name: "guest",
		}

		mockRoleRepo.EXPECT().GetByName("guest").Return(guestRole, nil)
		mockUserRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(user *models.User) error {
			assert.Equal(t, req.Username, user.Username)
			assert.Equal(t, req.Email, user.Email)
			assert.NotEmpty(t, user.Password) // Password should be hashed
			assert.Equal(t, guestRole.ID, user.RoleID)
			return nil
		})

		user, err := svc.CreateUser(req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Username, user.Username)
	})

	t.Run("RoleNotFound", func(t *testing.T) {
		req := dto.CreateUserRequestBody{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
		}

		expectedError := errors.New("role not found")
		mockRoleRepo.EXPECT().GetByName("guest").Return(nil, expectedError)

		user, err := svc.CreateUser(req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, expectedError, err)
	})

	t.Run("PasswordError", func(t *testing.T) {
		// Mock bcrypt behavior via interface injection? No, bcrypt in standard lib.
		// However, bcrypt fails if password is too long (>72 chars)
		longPassword := make([]byte, 73)
		for i := range longPassword {
			longPassword[i] = 'a'
		}

		req := dto.CreateUserRequestBody{
			Username: "newuser",
			Email:    "new@example.com",
			Password: string(longPassword),
		}

		user, err := svc.CreateUser(req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "bcrypt: password length exceeds 72 bytes")
	})

	t.Run("CreateRepoError", func(t *testing.T) {
		req := dto.CreateUserRequestBody{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
		}

		roleIDStr := "123e4567-e89b-12d3-a456-426614174001"
		guestRole := &models.Role{
			ID:   types.UUID(uuid.MustParse(roleIDStr)),
			Name: "guest",
		}

		expectedError := errors.New("db error")

		mockRoleRepo.EXPECT().GetByName("guest").Return(guestRole, nil)
		mockUserRepo.EXPECT().Create(gomock.Any()).Return(expectedError)

		user, err := svc.CreateUser(req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, expectedError, err)
	})
}

func TestUserService_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repositories.NewMockUserRepository(ctrl)
	svc := services.NewUserService(mockUserRepo, nil, "secret", nil)

	t.Run("Success", func(t *testing.T) {
		req := dto.LoginUserRequestBody{
			Email:    "test@example.com",
			Password: "password123",
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		expectedUser := &models.User{
			Email:    req.Email,
			Password: string(hashedPassword),
			Role: models.Role{
				Name: "user",
			},
		}

		mockUserRepo.EXPECT().GetByEmail(req.Email).Return(expectedUser, nil)

		token, err := svc.LoginUser(req)

		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.NotEmpty(t, *token)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		req := dto.LoginUserRequestBody{
			Email:    "test@example.com",
			Password: "password123",
		}

		mockUserRepo.EXPECT().GetByEmail(req.Email).Return(nil, errors.New("user not found"))

		token, err := svc.LoginUser(req)

		assert.Error(t, err)
		assert.Nil(t, token)
	})

	t.Run("WrongPassword", func(t *testing.T) {
		req := dto.LoginUserRequestBody{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
		expectedUser := &models.User{
			Email:    req.Email,
			Password: string(hashedPassword),
		}

		mockUserRepo.EXPECT().GetByEmail(req.Email).Return(expectedUser, nil)

		token, err := svc.LoginUser(req)

		assert.Error(t, err)
		assert.Nil(t, token)
		assert.Contains(t, err.Error(), "hashedPassword is not the hash of the given password")
	})

	// Testing signing error is difficult because it depends on jwt library internals or key type.
	// Since we pass a string secret, it uses HS256 which supports any byte slice.
	// It's hard to trigger "Sign invalid" without modifying the code or using invalid key type which we can't easily do via NewUserService.
	// We'll skip forcing signing error as it's unlikely with string secret.
}

func TestUserService_CreateUserWithRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repositories.NewMockUserRepository(ctrl)
	mockRoleRepo := mock_repositories.NewMockRoleRepository(ctrl)

	// Setup sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// GORM with sqlmock
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	svc := services.NewUserService(mockUserRepo, mockRoleRepo, "secret", gormDB)

	t.Run("Success_ExistingRole", func(t *testing.T) {
		req := dto.CreateUserWithRoleRequestBody{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
			RoleName: "admin",
		}

		roleIDStr := "123e4567-e89b-12d3-a456-426614174002"
		existingRole := &models.Role{
			ID:   types.UUID(uuid.MustParse(roleIDStr)),
			Name: "admin",
		}

		// Expectations
		mock.ExpectBegin()

		mockRoleRepo.EXPECT().GetByName("admin").Return(existingRole, nil)
		mockUserRepo.EXPECT().CreateWithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(func(tx *gorm.DB, user *models.User) error {
			assert.Equal(t, req.Username, user.Username)
			assert.Equal(t, req.Email, user.Email)
			assert.Equal(t, existingRole.ID, user.RoleID)
			return nil
		})

		mock.ExpectCommit()

		user, err := svc.CreateUserWithRole(req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Username, user.Username)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success_NewRole", func(t *testing.T) {
		req := dto.CreateUserWithRoleRequestBody{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
			RoleName: "custom_role",
		}

		// Expectations
		mock.ExpectBegin()

		// Return nil (not found) for GetByName
		mockRoleRepo.EXPECT().GetByName("custom_role").Return(nil, gorm.ErrRecordNotFound)

		// Expect CreateWithTransaction for Role
		mockRoleRepo.EXPECT().CreateWithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(func(tx *gorm.DB, role *models.Role) error {
			assert.Equal(t, "custom_role", role.Name)
			// Assign ID so User can use it
			role.ID = types.UUID(uuid.New())
			return nil
		})

		// Expect CreateWithTransaction for User
		mockUserRepo.EXPECT().CreateWithTransaction(gomock.Any(), gomock.Any()).Return(nil)

		mock.ExpectCommit()

		user, err := svc.CreateUserWithRole(req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("TransactionError", func(t *testing.T) {
		req := dto.CreateUserWithRoleRequestBody{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
			RoleName: "admin",
		}

		expectedError := errors.New("create error")

		// Expectations
		mock.ExpectBegin()

		mockRoleRepo.EXPECT().GetByName("admin").Return(nil, gorm.ErrRecordNotFound)
		mockRoleRepo.EXPECT().CreateWithTransaction(gomock.Any(), gomock.Any()).Return(expectedError)

		mock.ExpectRollback()

		user, err := svc.CreateUserWithRole(req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, expectedError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("PasswordError", func(t *testing.T) {
		longPassword := make([]byte, 73)
		for i := range longPassword {
			longPassword[i] = 'a'
		}
		req := dto.CreateUserWithRoleRequestBody{
			Username: "newuser",
			Email:    "new@example.com",
			Password: string(longPassword),
			RoleName: "admin",
		}

		user, err := svc.CreateUserWithRole(req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "bcrypt: password length exceeds 72 bytes")
		// No DB interaction expected
	})
}

func TestUserService_GetUsersWithPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repositories.NewMockUserRepository(ctrl)
	svc := services.NewUserService(mockUserRepo, nil, "secret", nil)

	t.Run("Success", func(t *testing.T) {
		filter := map[string]interface{}{}
		page := 1
		pageSize := 10

		expectedUsers := []models.User{
			{Username: "user1"},
			{Username: "user2"},
		}
		var totalRows int64 = 2
		totalPages := 1

		mockUserRepo.EXPECT().FindAllByFilterAndPage(filter, page, pageSize).Return(expectedUsers, totalRows, totalPages, nil)

		users, count, pages, err := svc.GetUsersWithPagination(filter, page, pageSize)

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Equal(t, totalRows, count)
		assert.Equal(t, totalPages, pages)
	})
}
