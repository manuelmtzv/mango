package store

import (
	"context"

	"github.com/manuelmtzv/mangocatnotes-api/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStorage interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByEmailOrUsername(ctx context.Context, identifier string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type NoteStorage interface {
	Create(ctx context.Context, note *models.Note) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Note, error)
	GetAll(ctx context.Context, userID primitive.ObjectID, page, limit int64, search string, tags []string) ([]models.Note, int64, error)
	Update(ctx context.Context, note *models.Note) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	AttachTags(ctx context.Context, noteID primitive.ObjectID, tagIDs []primitive.ObjectID) error
	DetachTag(ctx context.Context, noteID primitive.ObjectID, tagID primitive.ObjectID) error
	GetTags(ctx context.Context, noteID primitive.ObjectID) ([]models.Tag, error)
}

type TagStorage interface {
	Create(ctx context.Context, tag *models.Tag) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Tag, error)
	GetAll(ctx context.Context, userID primitive.ObjectID, page, limit int64) ([]models.Tag, int64, error)
	Update(ctx context.Context, tag *models.Tag) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindOrCreate(ctx context.Context, userID primitive.ObjectID, names []string) ([]models.Tag, error)
}
