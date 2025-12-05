package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/manuelmtzv/mangocatnotes-api/internal/models"
)

type UserStorage interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmailOrUsername(ctx context.Context, identifier string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type NoteStorage interface {
	Create(ctx context.Context, note *models.Note) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Note, error)
	GetAll(ctx context.Context, userID uuid.UUID, page, limit int64, search string, tags []string) ([]models.Note, int64, error)
	Update(ctx context.Context, note *models.Note) error
	Delete(ctx context.Context, id uuid.UUID) error
	AttachTags(ctx context.Context, noteID uuid.UUID, tagIDs []uuid.UUID) error
	DetachTag(ctx context.Context, noteID uuid.UUID, tagID uuid.UUID) error
	GetTags(ctx context.Context, noteID uuid.UUID) ([]models.Tag, error)
}

type TagStorage interface {
	Create(ctx context.Context, tag *models.Tag) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Tag, error)
	GetAll(ctx context.Context, userID uuid.UUID, page, limit int64) ([]models.Tag, int64, error)
	Update(ctx context.Context, tag *models.Tag) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindOrCreate(ctx context.Context, userID uuid.UUID, names []string) ([]models.Tag, error)
}
