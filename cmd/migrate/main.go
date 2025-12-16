package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manuelmtzv/mangocatnotes-api/internal/env"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUser struct {
	ID        primitive.ObjectID `bson:"_id"`
	Email     string             `bson:"email"`
	Username  string             `bson:"username"`
	Hash      string             `bson:"hash"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type MongoNote struct {
	ID        primitive.ObjectID   `bson:"_id"`
	UserID    primitive.ObjectID   `bson:"userId"`
	Title     string               `bson:"title"`
	Content   string               `bson:"content"`
	Archived  bool                 `bson:"archived"`
	TagIDs    []primitive.ObjectID `bson:"tagIDs"`
	CreatedAt time.Time            `bson:"createdAt"`
	UpdatedAt time.Time            `bson:"updatedAt"`
}

type MongoTag struct {
	ID        primitive.ObjectID   `bson:"_id"`
	UserID    primitive.ObjectID   `bson:"userId"`
	Name      string               `bson:"name"`
	Color     string               `bson:"color"`
	NoteIDs   []primitive.ObjectID `bson:"noteIDs"`
	CreatedAt time.Time            `bson:"createdAt"`
	UpdatedAt time.Time            `bson:"updatedAt"`
}

func main() {
	env.Load()

	mongoURI := env.GetRequired("MONGO_URL")
	mongoDBName := env.GetRequired("MONGO_DB_NAME")

	pgAddr := env.GetRequired("DB_ADDR")

	dryRun := env.GetString("DRY_RUN", "false") == "true"

	if dryRun {
		log.Println("🔍 Running in DRY RUN mode - no data will be written")
	}

	log.Println("🔌 Connecting to MongoDB...")
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	mongoDB := mongoClient.Database(mongoDBName)

	log.Println("🔌 Connecting to PostgreSQL...")
	pgConfig, err := pgxpool.ParseConfig(pgAddr)
	if err != nil {
		log.Fatalf("Failed to parse PostgreSQL config: %v", err)
	}

	pgPool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgPool.Close()

	log.Println("✅ Connected to both databases")
	log.Println("")

	userIDMap := make(map[primitive.ObjectID]uuid.UUID)
	noteIDMap := make(map[primitive.ObjectID]uuid.UUID)
	tagIDMap := make(map[primitive.ObjectID]uuid.UUID)

	log.Println("👥 Migrating users...")
	if err := migrateUsers(mongoDB, pgPool, userIDMap, dryRun); err != nil {
		log.Fatalf("Failed to migrate users: %v", err)
	}

	log.Println("📝 Migrating notes...")
	if err := migrateNotes(mongoDB, pgPool, userIDMap, noteIDMap, dryRun); err != nil {
		log.Fatalf("Failed to migrate notes: %v", err)
	}

	log.Println("🏷️  Migrating tags...")
	if err := migrateTags(mongoDB, pgPool, userIDMap, tagIDMap, dryRun); err != nil {
		log.Fatalf("Failed to migrate tags: %v", err)
	}

	log.Println("🔗 Creating note-tag associations...")
	if err := migrateNoteTags(mongoDB, pgPool, noteIDMap, tagIDMap, dryRun); err != nil {
		log.Fatalf("Failed to migrate note-tag associations: %v", err)
	}

	log.Println("")
	log.Println("✨ Migration completed successfully!")
	log.Printf("📊 Migrated: %d users, %d notes, %d tags", len(userIDMap), len(noteIDMap), len(tagIDMap))
}

func migrateUsers(mongoDB *mongo.Database, pg *pgxpool.Pool, idMap map[primitive.ObjectID]uuid.UUID, dryRun bool) error {
	cursor, err := mongoDB.Collection("users").Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	count := 0
	for cursor.Next(context.Background()) {
		var user MongoUser
		if err := cursor.Decode(&user); err != nil {
			return fmt.Errorf("failed to decode user: %w", err)
		}

		newID := uuid.New()
		idMap[user.ID] = newID

		if !dryRun {
			query := `
				INSERT INTO users (id, email, username, hash, name, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`
			_, err := pg.Exec(context.Background(), query,
				newID, user.Email, user.Username, user.Hash, user.Name, user.CreatedAt, user.UpdatedAt)
			if err != nil {
				return fmt.Errorf("failed to insert user %s: %w", user.Email, err)
			}
		}
		count++
	}

	log.Printf("  ✓ Migrated %d users", count)
	return nil
}

func migrateNotes(mongoDB *mongo.Database, pg *pgxpool.Pool, userIDMap map[primitive.ObjectID]uuid.UUID, noteIDMap map[primitive.ObjectID]uuid.UUID, dryRun bool) error {
	cursor, err := mongoDB.Collection("notes").Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	count := 0
	for cursor.Next(context.Background()) {
		var note MongoNote
		if err := cursor.Decode(&note); err != nil {
			return fmt.Errorf("failed to decode note: %w", err)
		}

		newID := uuid.New()
		noteIDMap[note.ID] = newID

		pgUserID, ok := userIDMap[note.UserID]
		if !ok {
			log.Printf("  ⚠️  Skipping note %s: user not found", note.ID.Hex())
			continue
		}

		if !dryRun {
			query := `
				INSERT INTO notes (id, user_id, title, content, archived, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`
			_, err := pg.Exec(context.Background(), query,
				newID, pgUserID, note.Title, note.Content, note.Archived, note.CreatedAt, note.UpdatedAt)
			if err != nil {
				return fmt.Errorf("failed to insert note %s: %w", note.Title, err)
			}
		}
		count++
	}

	log.Printf("  ✓ Migrated %d notes", count)
	return nil
}

func migrateTags(mongoDB *mongo.Database, pg *pgxpool.Pool, userIDMap map[primitive.ObjectID]uuid.UUID, tagIDMap map[primitive.ObjectID]uuid.UUID, dryRun bool) error {
	cursor, err := mongoDB.Collection("tags").Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	count := 0
	for cursor.Next(context.Background()) {
		var tag MongoTag
		if err := cursor.Decode(&tag); err != nil {
			return fmt.Errorf("failed to decode tag: %w", err)
		}

		newID := uuid.New()
		tagIDMap[tag.ID] = newID

		pgUserID, ok := userIDMap[tag.UserID]
		if !ok {
			log.Printf("  ⚠️  Skipping tag %s: user not found", tag.Name)
			continue
		}

		if !dryRun {
			query := `
				INSERT INTO tags (id, user_id, name, color, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6)
			`
			_, err := pg.Exec(context.Background(), query,
				newID, pgUserID, tag.Name, tag.Color, tag.CreatedAt, tag.UpdatedAt)
			if err != nil {
				return fmt.Errorf("failed to insert tag %s: %w", tag.Name, err)
			}
		}
		count++
	}

	log.Printf("  ✓ Migrated %d tags", count)
	return nil
}

func migrateNoteTags(mongoDB *mongo.Database, pg *pgxpool.Pool, noteIDMap map[primitive.ObjectID]uuid.UUID, tagIDMap map[primitive.ObjectID]uuid.UUID, dryRun bool) error {
	cursor, err := mongoDB.Collection("notes").Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	count := 0
	for cursor.Next(context.Background()) {
		var note MongoNote
		if err := cursor.Decode(&note); err != nil {
			return fmt.Errorf("failed to decode note: %w", err)
		}

		pgNoteID, ok := noteIDMap[note.ID]
		if !ok {
			continue
		}

		for _, mongoTagID := range note.TagIDs {
			pgTagID, ok := tagIDMap[mongoTagID]
			if !ok {
				log.Printf("  ⚠️  Skipping tag association: tag not found")
				continue
			}

			if !dryRun {
				query := `
					INSERT INTO note_tags (note_id, tag_id, created_at)
					VALUES ($1, $2, $3)
					ON CONFLICT (note_id, tag_id) DO NOTHING
				`
				_, err := pg.Exec(context.Background(), query, pgNoteID, pgTagID, time.Now())
				if err != nil {
					return fmt.Errorf("failed to insert note-tag association: %w", err)
				}
			}
			count++
		}
	}

	log.Printf("  ✓ Created %d note-tag associations", count)
	return nil
}
