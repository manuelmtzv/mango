package models

import (
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"userId"`
	Name      string    `db:"name" json:"name"`
	Color     string    `db:"color" json:"color"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type PaginatedTagsResponse struct {
	Data []Tag              `json:"data"`
	Meta PaginationMetadata `json:"meta"`
}
