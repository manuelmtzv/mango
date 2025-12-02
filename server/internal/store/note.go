package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manuelmtzv/mangocatnotes-api/internal/models"
)

type PostgresNoteStore struct {
	pool *pgxpool.Pool
}

func NewNoteStore(pool *pgxpool.Pool) *PostgresNoteStore {
	return &PostgresNoteStore{
		pool: pool,
	}
}

func (s *PostgresNoteStore) Create(ctx context.Context, note *models.Note) error {
	query := `
		INSERT INTO notes (user_id, title, content, archived, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	now := time.Now()
	note.CreatedAt = now
	note.UpdatedAt = now

	err := s.pool.QueryRow(ctx, query,
		note.UserID,
		note.Title,
		note.Content,
		note.Archived,
		note.CreatedAt,
		note.UpdatedAt,
	).Scan(&note.ID)

	return err
}

func (s *PostgresNoteStore) AttachTags(ctx context.Context, noteID uuid.UUID, tagIDs []uuid.UUID) error {
	if len(tagIDs) == 0 {
		return nil
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, tagID := range tagIDs {
		query := `
			INSERT INTO note_tags (note_id, tag_id, created_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (note_id, tag_id) DO NOTHING
		`
		_, err := tx.Exec(ctx, query, noteID, tagID, time.Now())
		if err != nil {
			return err
		}
	}

	query := `UPDATE notes SET updated_at = $1 WHERE id = $2`
	_, err = tx.Exec(ctx, query, time.Now(), noteID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *PostgresNoteStore) DetachTag(ctx context.Context, noteID uuid.UUID, tagID uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `DELETE FROM note_tags WHERE note_id = $1 AND tag_id = $2`
	_, err = tx.Exec(ctx, query, noteID, tagID)
	if err != nil {
		return err
	}

	updateQuery := `UPDATE notes SET updated_at = $1 WHERE id = $2`
	_, err = tx.Exec(ctx, updateQuery, time.Now(), noteID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *PostgresNoteStore) GetByID(ctx context.Context, id uuid.UUID) (*models.Note, error) {
	query := `
		SELECT id, user_id, title, content, archived, created_at, updated_at
		FROM notes
		WHERE id = $1
	`
	var note models.Note
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&note.ID,
		&note.UserID,
		&note.Title,
		&note.Content,
		&note.Archived,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (s *PostgresNoteStore) GetAll(ctx context.Context, userID uuid.UUID, page, limit int64, search string, tags []string) ([]models.Note, int64, error) {
	offset := (page - 1) * limit

	baseQuery := `
		FROM notes n
		WHERE n.user_id = $1
	`
	args := []interface{}{userID}
	argCount := 1

	if search != "" {
		argCount++
		baseQuery += fmt.Sprintf(" AND n.title ILIKE $%d", argCount)
		args = append(args, "%"+search+"%")
	}

	if len(tags) > 0 {
		argCount++
		baseQuery += fmt.Sprintf(`
			AND EXISTS (
				SELECT 1 FROM note_tags nt
				JOIN tags t ON t.id = nt.tag_id
				WHERE nt.note_id = n.id
				AND t.name = ANY($%d)
				AND t.user_id = $1
			)
		`, argCount)
		args = append(args, tags)
	}

	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int64
	err := s.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	dataQuery := `
		SELECT n.id, n.user_id, n.title, n.content, n.archived, n.created_at, n.updated_at
	` + baseQuery + `
		ORDER BY n.updated_at DESC
		LIMIT $` + fmt.Sprintf("%d", argCount+1) + ` OFFSET $` + fmt.Sprintf("%d", argCount+2)

	args = append(args, limit, offset)

	rows, err := s.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
			&note.Content,
			&note.Archived,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}

func (s *PostgresNoteStore) Update(ctx context.Context, note *models.Note) error {
	note.UpdatedAt = time.Now()
	query := `
		UPDATE notes
		SET title = $1, content = $2, archived = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := s.pool.Exec(ctx, query,
		note.Title,
		note.Content,
		note.Archived,
		note.UpdatedAt,
		note.ID,
	)
	return err
}

func (s *PostgresNoteStore) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM notes WHERE id = $1`
	_, err := s.pool.Exec(ctx, query, id)
	return err
}

func (s *PostgresNoteStore) GetTags(ctx context.Context, noteID uuid.UUID) ([]models.Tag, error) {
	query := `
		SELECT t.id, t.user_id, t.name, t.color, t.created_at, t.updated_at
		FROM tags t
		INNER JOIN note_tags nt ON t.id = nt.tag_id
		WHERE nt.note_id = $1
		ORDER BY t.name
	`

	rows, err := s.pool.Query(ctx, query, noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(
			&tag.ID,
			&tag.UserID,
			&tag.Name,
			&tag.Color,
			&tag.CreatedAt,
			&tag.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}
