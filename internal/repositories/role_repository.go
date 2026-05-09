package repositories

import (
	"context"
	"fdlp-standard-api/internal/models"
	"time"

	"fdlp-standard-api/pkg/utils"

	"github.com/sony/gobreaker"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoleRepository interface {
	GetById(id string) (*models.Role, error)
	GetByName(name string) (*models.Role, error)
	Create(role *models.Role) error
}

type roleRepository struct {
	db      *mongo.Database
	breaker *gobreaker.TwoStepCircuitBreaker
}

func NewRoleRepository(db *mongo.Database) RoleRepository {
	return &roleRepository{
		db:      db,
		breaker: utils.NewCircuitBreaker("RoleRepository"),
	}
}

func (repo *roleRepository) GetById(id string) (*models.Role, error) {
	result, err := utils.ExecuteWithBreaker(repo.breaker, func() (interface{}, error) {
		var role models.Role
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := repo.db.Collection("roles").FindOne(ctx, bson.M{"role_id": id}).Decode(&role)
		if err != nil {
			return nil, err
		}
		return &role, nil
	})

	if err != nil {
		return nil, err
	}
	return result.(*models.Role), nil
}

func (repo *roleRepository) GetByName(name string) (*models.Role, error) {
	result, err := utils.ExecuteWithBreaker(repo.breaker, func() (interface{}, error) {
		var role models.Role
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := repo.db.Collection("roles").FindOne(ctx, bson.M{"name": name}).Decode(&role)
		if err != nil {
			return nil, err
		}
		return &role, nil
	})

	if err != nil {
		return nil, err
	}
	return result.(*models.Role), nil
}

func (repo *roleRepository) Create(role *models.Role) error {
	_, err := utils.ExecuteWithBreaker(repo.breaker, func() (interface{}, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return repo.db.Collection("roles").InsertOne(ctx, role)
	})
	return err
}
