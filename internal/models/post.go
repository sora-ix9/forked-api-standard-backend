package models

import (
	"time"

	"fdlp-standard-api/internal/types"

	"github.com/google/uuid"
)

type Post struct {
	PostID    types.UUID `bson:"post_id,omitempty"`
	Title     string     `bson:"title"`
	Content   string     `bson:"content"`
	Author    string     `bson:"author"`
	CreatedAt time.Time  `bson:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at"`
}

func (p *Post) BeforeCreate() {
	if p.PostID == [16]byte{} {
		p.PostID = types.UUID(uuid.New())
	}
}
