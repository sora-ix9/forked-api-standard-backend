package models

import (
	"fdlp-standard-api/internal/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        types.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key;"`
	Username  string     `gorm:"uniqueIndex;not null"`
	Email     string     `gorm:"uniqueIndex;not null"`
	Fisrtname string
	Password  string     `gorm:"not null"`
	RoleID    types.UUID `gorm:"type:uuid;not null"`
	Role      Role       `gorm:"foreignKey:RoleID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == [16]byte{} {
		u.ID = types.UUID(uuid.New())
	}
	return
}
