package dto

import "github.com/golang-jwt/jwt/v5"

type JwtCustomClaims struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	RoleName string `json:"role_name"`
	jwt.RegisteredClaims
}
