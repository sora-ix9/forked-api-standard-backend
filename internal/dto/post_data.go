package dto

type CreatePostRequestBody struct {
	Title   string `json:"title"   validate:"required"`
	Content string `json:"content" validate:"required"`
	Author  string `json:"author"  validate:"required"`
}
