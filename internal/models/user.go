package models

import (
	"fdlp-standard-api/internal/types"
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID    types.UUID `bson:"user_id,omitempty"`
	Username  string     `bson:"username"`
	Email     string     `bson:"email"`
	Fisrtname string     `bson:"firstname"`
	Password  string     `bson:"password"`
	RoleID    types.UUID `bson:"role_id"`
	Role      Role       `bson:"role,omitempty"`
	CreatedAt time.Time  `bson:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at"`
}

func (u *User) BeforeCreate() {
	if u.UserID == [16]byte{} {
		u.UserID = types.UUID(uuid.New())
	}
}
