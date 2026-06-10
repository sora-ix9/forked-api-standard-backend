package repositories

import (
	"context"
	"fdlp-standard-api/internal/models"
	"fdlp-standard-api/pkg/utils"
	"time"

	"github.com/sony/gobreaker"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepository interface {
	GetById(id string) (*models.Post, error)
	Create(post *models.Post) error
}

type postRepository struct {
	db      *mongo.Database
	breaker *gobreaker.TwoStepCircuitBreaker
}

func NewPostRepository(db *mongo.Database) PostRepository {
	return &postRepository{
		db:      db,
		breaker: utils.NewCircuitBreaker("PostRepository"),
	}
}

func (repo *postRepository) GetById(id string) (*models.Post, error) {
	result, err := utils.ExecuteWithBreaker(repo.breaker, func() (interface{}, error) {
		var post models.Post
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := repo.db.Collection("posts").FindOne(ctx, bson.M{"post_id": id}).Decode(&post)
		if err != nil {
			return nil, err
		}
		return &post, nil
	})

	if err != nil {
		return nil, err
	}
	return result.(*models.Post), nil
}

func (repo *postRepository) Create(post *models.Post) error {
	_, err := utils.ExecuteWithBreaker(repo.breaker, func() (interface{}, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return repo.db.Collection("posts").InsertOne(ctx, post)
	})
	return err
}
