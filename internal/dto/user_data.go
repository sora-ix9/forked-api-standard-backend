package dto

type LoginUserRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateUserRequestBody struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateUserWithRoleRequestBody struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	RoleName string `json:"role_name" validate:"required"`
}

type UserGetByIdResponse struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	RoleName        string `json:"role_name"`
	RoleDescription string `json:"role_description"`
}
