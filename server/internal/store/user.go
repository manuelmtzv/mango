package store

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/manuelmtzv/mangocatnotes-api/internal/models"
)

type PostgresUserStore struct {
	pool *pgxpool.Pool
}

func NewUserStore(pool *pgxpool.Pool) *PostgresUserStore {
	return &PostgresUserStore{
		pool: pool,
	}
}

func (s *PostgresUserStore) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, username, hash, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := s.pool.QueryRow(ctx, query,
		user.Email,
		user.Username,
		user.Hash,
		user.Name,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	return err
}

func (s *PostgresUserStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, username, hash, name, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var user models.User
	err := s.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Hash,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *PostgresUserStore) GetByEmailOrUsername(ctx context.Context, identifier string) (*models.User, error) {
	query := `
		SELECT id, email, username, hash, name, created_at, updated_at
		FROM users
		WHERE email = $1 OR username = $1
	`
	var user models.User
	err := s.pool.QueryRow(ctx, query, identifier).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Hash,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *PostgresUserStore) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, username, hash, name, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user models.User
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Hash,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *PostgresUserStore) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	query := `
		UPDATE users
		SET email = $1, username = $2, name = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := s.pool.Exec(ctx, query,
		user.Email,
		user.Username,
		user.Name,
		user.UpdatedAt,
		user.ID,
	)
	return err
}

func (s *PostgresUserStore) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := s.pool.Exec(ctx, query, id)
	return err
}
