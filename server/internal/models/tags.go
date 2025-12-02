package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID   `bson:"userId" json:"userId"`
	Name      string               `bson:"name" json:"name"`
	Color     string               `bson:"color" json:"color"`
	NoteIDs   []primitive.ObjectID `bson:"noteIDs" json:"noteIDs"`
	CreatedAt time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time            `bson:"updatedAt" json:"updatedAt"`
}

type PaginatedTagsResponse struct {
	Data []Tag              `json:"data"`
	Meta PaginationMetadata `json:"meta"`
}
