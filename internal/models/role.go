package models

import (
	"time"

	"fdlp-standard-api/internal/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          types.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key;"`
	Name        string     `gorm:"not null;uniqueIndex"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == [16]byte{} {
		r.ID = types.UUID(uuid.New())
	}
	return
}
