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

type PostgresTagStore struct {
	pool *pgxpool.Pool
}

func NewTagStore(pool *pgxpool.Pool) *PostgresTagStore {
	return &PostgresTagStore{
		pool: pool,
	}
}

func (s *PostgresTagStore) Create(ctx context.Context, tag *models.Tag) error {
	query := `
		INSERT INTO tags (user_id, name, color, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	now := time.Now()
	tag.CreatedAt = now
	tag.UpdatedAt = now

	err := s.pool.QueryRow(ctx, query,
		tag.UserID,
		tag.Name,
		tag.Color,
		tag.CreatedAt,
		tag.UpdatedAt,
	).Scan(&tag.ID)

	return err
}

func (s *PostgresTagStore) GetByID(ctx context.Context, id uuid.UUID) (*models.Tag, error) {
	query := `
		SELECT id, user_id, name, color, created_at, updated_at
		FROM tags
		WHERE id = $1
	`
	var tag models.Tag
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&tag.ID,
		&tag.UserID,
		&tag.Name,
		&tag.Color,
		&tag.CreatedAt,
		&tag.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *PostgresTagStore) GetAll(ctx context.Context, userID uuid.UUID, page, limit int64) ([]models.Tag, int64, error) {
	offset := (page - 1) * limit

	countQuery := `
		SELECT COUNT(*)
		FROM tags
		WHERE user_id = $1
	`
	var total int64
	err := s.pool.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	dataQuery := `
		SELECT id, user_id, name, color, created_at, updated_at
		FROM tags
		WHERE user_id = $1
		ORDER BY name
		LIMIT $2 OFFSET $3
	`

	rows, err := s.pool.Query(ctx, dataQuery, userID, limit, offset)
	if err != nil {
		return nil, 0, err
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
			return nil, 0, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}

func (s *PostgresTagStore) Update(ctx context.Context, tag *models.Tag) error {
	tag.UpdatedAt = time.Now()
	query := `
		UPDATE tags
		SET name = $1, color = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := s.pool.Exec(ctx, query,
		tag.Name,
		tag.Color,
		tag.UpdatedAt,
		tag.ID,
	)
	return err
}

func (s *PostgresTagStore) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tags WHERE id = $1`
	_, err := s.pool.Exec(ctx, query, id)
	return err
}

func (s *PostgresTagStore) FindOrCreate(ctx context.Context, userID uuid.UUID, names []string) ([]models.Tag, error) {
	if len(names) == 0 {
		return []models.Tag{}, nil
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var tags []models.Tag
	now := time.Now()

	for _, name := range names {
		query := `
			INSERT INTO tags (user_id, name, color, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (user_id, name) DO UPDATE SET updated_at = EXCLUDED.updated_at
			RETURNING id, user_id, name, color, created_at, updated_at
		`

		var tag models.Tag
		err := tx.QueryRow(ctx, query, userID, name, "#fff", now, now).Scan(
			&tag.ID,
			&tag.UserID,
			&tag.Name,
			&tag.Color,
			&tag.CreatedAt,
			&tag.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to find or create tag %s: %w", name, err)
		}
		tags = append(tags, tag)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return tags, nil
}
