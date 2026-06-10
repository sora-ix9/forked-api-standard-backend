package services

import (
	"fdlp-standard-api/internal/dto"
	"fdlp-standard-api/internal/models"
	"fdlp-standard-api/internal/repositories"
	"time"
)

type PostService interface {
	GetPost(id string) (*models.Post, error)
	CreatePost(data dto.CreatePostRequestBody) (*models.Post, error)
}

type postService struct {
	repo repositories.PostRepository
}

func NewPostService(repo repositories.PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) GetPost(id string) (*models.Post, error) {
	return s.repo.GetById(id)
}

func (s *postService) CreatePost(data dto.CreatePostRequestBody) (*models.Post, error) {
	now := time.Now()
	post := models.Post{
		Title:     data.Title,
		Content:   data.Content,
		Author:    data.Author,
		CreatedAt: now,
		UpdatedAt: now,
	}
	post.BeforeCreate()

	if err := s.repo.Create(&post); err != nil {
		return nil, err
	}

	return &post, nil
}
