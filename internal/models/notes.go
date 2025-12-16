package models

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"userId"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	Archived  bool      `db:"archived" json:"archived"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`

	Tags []Tag `db:"-" json:"tags"`
}

type PaginatedNotesResponse struct {
	Data []Note             `json:"data"`
	Meta PaginationMetadata `json:"meta"`
}
