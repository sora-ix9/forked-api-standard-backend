package repositories

import (
	"context"
	"fdlp-standard-api/internal/models"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"fdlp-standard-api/pkg/utils"

	"github.com/sony/gobreaker"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	GetById(id string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	FindAllByFilterAndPage(filter map[string]interface{}, page, pageSize int) ([]models.User, int64, int, error)
}

type userRepository struct {
	db      *mongo.Database
	breaker *gobreaker.TwoStepCircuitBreaker
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		db:      db,
		breaker: utils.NewCircuitBreaker("UserRepository"),
	}
}

func (repo *userRepository) GetById(id string) (*models.User, error) {
	result, err := utils.ExecuteWithBreaker(repo.breaker, func() (interface{}, error) {
		var user models.User
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := repo.db.Collection("users").FindOne(ctx, bson.M{"user_id": id}).Decode(&user)
		if err != nil {
			return nil, err
		}
		return &user, nil
	})

	if err != nil {
		return nil, err
	}
	return result.(*models.User), nil
}

func (repo *userRepository) GetByEmail(email string) (*models.User, error) {
	result, err := utils.ExecuteWithBreaker(repo.breaker, func() (interface{}, error) {
		var user models.User
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := repo.db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
		if err != nil {
			return nil, err
		}
		return &user, nil
	})

	if err != nil {
		return nil, err
	}
	return result.(*models.User), nil
}

func (repo *userRepository) Create(user *models.User) error {
	_, err := utils.ExecuteWithBreaker(repo.breaker, func() (interface{}, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return repo.db.Collection("users").InsertOne(ctx, user)
	})
	return err
}

func (repo *userRepository) Update(user *models.User) error {
	_, err := utils.ExecuteWithBreaker(repo.breaker, func() (interface{}, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return repo.db.Collection("users").UpdateOne(
			ctx,
			bson.M{"user_id": user.UserID},
			bson.M{"$set": user},
		)
	})
	return err
}

func (repo *userRepository) FindAllByFilterAndPage(filter map[string]interface{}, page, pageSize int) ([]models.User, int64, int, error) {
	type findResult struct {
		users      []models.User
		totalRows  int64
		totalPages int
	}

	result, err := utils.ExecuteWithBreaker(repo.breaker, func() (interface{}, error) {
		var users []models.User
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		mongoFilter := bson.M{}
		for k, v := range filter {
			mongoFilter[k] = v
		}

		totalRows, err := repo.db.Collection("users").CountDocuments(ctx, mongoFilter)
		if err != nil {
			return nil, err
		}

		opts := options.Find().
			SetSkip(int64((page - 1) * pageSize)).
			SetLimit(int64(pageSize))

		cursor, err := repo.db.Collection("users").Find(ctx, mongoFilter, opts)
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &users); err != nil {
			return nil, err
		}

		totalPages := int(math.Ceil(float64(totalRows) / float64(pageSize)))
		return &findResult{users: users, totalRows: totalRows, totalPages: totalPages}, nil
	})

	if err != nil {
		return nil, 0, 0, err
	}

	res := result.(*findResult)
	return res.users, res.totalRows, res.totalPages, nil
}
