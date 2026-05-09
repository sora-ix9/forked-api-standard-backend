package repositories

import (
	"context"
	"fdlp-standard-api/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoleRepository interface {
	GetById(id string) (*models.Role, error)
	GetByName(name string) (*models.Role, error)
	Create(role *models.Role) error
}

type roleRepository struct {
	db *mongo.Database
}

func NewRoleRepository(db *mongo.Database) RoleRepository {
	return &roleRepository{db: db}
}

func (repo *roleRepository) GetById(id string) (*models.Role, error) {
	var role models.Role
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.db.Collection("roles").FindOne(ctx, bson.M{"role_id": id}).Decode(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (repo *roleRepository) GetByName(name string) (*models.Role, error) {
	var role models.Role
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.db.Collection("roles").FindOne(ctx, bson.M{"name": name}).Decode(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (repo *roleRepository) Create(role *models.Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := repo.db.Collection("roles").InsertOne(ctx, role)
	return err
}
