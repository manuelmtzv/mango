package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID   `bson:"userId" json:"userId"`
	Title     string               `bson:"title" json:"title"`
	Content   string               `bson:"content" json:"content"`
	Archived  bool                 `bson:"archived" json:"archived"`
	TagIDs    []primitive.ObjectID `bson:"tagIDs" json:"tagIDs"`
	CreatedAt time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time            `bson:"updatedAt" json:"updatedAt"`

	Tags []Tag `bson:"-" json:"tags"`
}

type PaginatedNotesResponse struct {
	Data []Note             `json:"data"`
	Meta PaginationMetadata `json:"meta"`
}
