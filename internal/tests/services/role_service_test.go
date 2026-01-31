package services_test

import (
	"errors"
	"fdlp-standard-api/internal/dto"
	"fdlp-standard-api/internal/models"
	"fdlp-standard-api/internal/services"
	mock_repositories "fdlp-standard-api/internal/tests/mock/repositories"
	"fdlp-standard-api/internal/types"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRoleService_GetRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockRoleRepository(ctrl)
	svc := services.NewRoleService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		roleIDStr := "123e4567-e89b-12d3-a456-426614174000"
		roleID := types.UUID(uuid.MustParse(roleIDStr))

		expectedRole := &models.Role{
			ID:          roleID,
			Name:        "admin",
			Description: "Administrator",
		}

		mockRepo.EXPECT().GetById(roleIDStr).Return(expectedRole, nil)

		result, err := svc.GetRole(roleIDStr)

		assert.NoError(t, err)
		assert.Equal(t, expectedRole, result)
	})

	t.Run("Error", func(t *testing.T) {
		roleIDStr := "123e4567-e89b-12d3-a456-426614174000"
		expectedError := errors.New("role not found")

		mockRepo.EXPECT().GetById(roleIDStr).Return(nil, expectedError)

		result, err := svc.GetRole(roleIDStr)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
	})
}

func TestRoleService_CreateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repositories.NewMockRoleRepository(ctrl)
	svc := services.NewRoleService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		req := dto.CreateRoleRequestBody{
			Name:        "moderator",
			Description: "Forum Moderator",
		}

		mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(role *models.Role) error {
			assert.Equal(t, req.Name, role.Name)
			assert.Equal(t, req.Description, role.Description)
			return nil
		})

		role, err := svc.CreateRole(req)

		assert.NoError(t, err)
		assert.NotNil(t, role)
		assert.Equal(t, req.Name, role.Name)
		assert.Equal(t, req.Description, role.Description)
	})

	t.Run("Error", func(t *testing.T) {
		req := dto.CreateRoleRequestBody{
			Name:        "moderator",
			Description: "Forum Moderator",
		}

		expectedError := errors.New("failed to create role")
		mockRepo.EXPECT().Create(gomock.Any()).Return(expectedError)

		role, err := svc.CreateRole(req)

		assert.Error(t, err)
		assert.Nil(t, role)
		assert.Equal(t, expectedError, err)
	})
}
