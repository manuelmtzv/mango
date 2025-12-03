package seed

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/manuelmtzv/mangocatnotes-api/internal/auth"
	"github.com/manuelmtzv/mangocatnotes-api/internal/models"
	"github.com/manuelmtzv/mangocatnotes-api/internal/store"
)

type Seeder struct {
	store *store.Storage
}

func NewSeeder(storage *store.Storage) *Seeder {
	return &Seeder{
		store: storage,
	}
}

func (s *Seeder) Run(ctx context.Context) error {
	log.Println("Starting database seeding...")

	if err := s.seedUsers(ctx); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

func (s *Seeder) seedUsers(ctx context.Context) error {
	log.Println("Seeding users...")

	users := []struct {
		Email    string
		Username string
		Password string
		Name     string
	}{
		{
			Email:    "test@example.com",
			Username: "testuser",
			Password: "password123",
			Name:     "Test User",
		},
		{
			Email:    "john@example.com",
			Username: "johndoe",
			Password: "password123",
			Name:     "John Doe",
		},
		{
			Email:    "jane@example.com",
			Username: "janedoe",
			Password: "password123",
			Name:     "Jane Doe",
		},
	}

	for _, userData := range users {
		existing, err := s.store.Users.GetByEmail(ctx, userData.Email)
		if err != nil {
			return err
		}
		var userID uuid.UUID
		if existing != nil {
			log.Printf("User %s already exists, seeding data...", userData.Email)
			userID = existing.ID
		} else {
			hash, err := auth.HashPassword(userData.Password)
			if err != nil {
				return err
			}

			user := &models.User{
				Email:    userData.Email,
				Username: userData.Username,
				Hash:     hash,
				Name:     userData.Name,
			}

			if err := s.store.Users.Create(ctx, user); err != nil {
				return err
			}
			userID = user.ID
			log.Printf("Created user: %s (%s)", user.Email, user.Username)
		}

		if err := s.seedTagsAndNotes(ctx, userID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Seeder) seedTagsAndNotes(ctx context.Context, userID uuid.UUID) error {
	tags := []string{"work", "personal", "ideas", "todo"}
	createdTags, err := s.store.Tags.FindOrCreate(ctx, userID, tags)
	if err != nil {
		return err
	}
	log.Printf("Created %d tags for user %s", len(createdTags), userID.String())

	notes := []struct {
		Title    string
		Content  string
		TagNames []string
	}{
		{
			Title:    "Welcome Note",
			Content:  "Welcome to MangoCatNotes! This is your first note.",
			TagNames: []string{"personal"},
		},
		{
			Title:    "Project Ideas",
			Content:  "1. Build a cool API\n2. Learn Go\n3. Profit?",
			TagNames: []string{"work", "ideas"},
		},
		{
			Title:    "Shopping List",
			Content:  "- Milk\n- Eggs\n- Coffee",
			TagNames: []string{"todo"},
		},
	}

	for _, noteData := range notes {
		var tagIDs []uuid.UUID
		for _, tagName := range noteData.TagNames {
			for _, tag := range createdTags {
				if tag.Name == tagName {
					tagIDs = append(tagIDs, tag.ID)
					break
				}
			}
		}

		note := &models.Note{
			UserID:  userID,
			Title:   noteData.Title,
			Content: noteData.Content,
		}

		if err := s.store.Notes.Create(ctx, note); err != nil {
			return err
		}
	}
	log.Printf("Created %d notes for user %s", len(notes), userID.String())

	return nil
}
