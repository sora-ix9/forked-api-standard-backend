package models

import (
	"time"

	"fdlp-standard-api/internal/types"

	"github.com/google/uuid"
)

type Role struct {
	RoleID      types.UUID `bson:"role_id,omitempty"`
	Name        string     `bson:"name"`
	Description string     `bson:"description"`
	CreatedAt   time.Time  `bson:"created_at"`
	UpdatedAt   time.Time  `bson:"updated_at"`
}

func (r *Role) BeforeCreate() {
	if r.RoleID == [16]byte{} {
		r.RoleID = types.UUID(uuid.New())
	}
}
