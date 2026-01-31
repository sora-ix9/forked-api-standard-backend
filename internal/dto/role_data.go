package dto

type CreateRoleRequestBody struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}
