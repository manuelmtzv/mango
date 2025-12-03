package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*DB, error) {
	config, err := pgxpool.ParseConfig(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	config.MaxConns = int32(maxOpenConns)
	config.MinConns = int32(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse max idle time: %w", err)
	}
	config.MaxConnIdleTime = duration

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	log.Println("Connected to PostgreSQL")
	return &DB{
		Pool: pool,
	}, nil
}

func (d *DB) Close() {
	if d.Pool != nil {
		d.Pool.Close()
		log.Println("PostgreSQL connection closed")
	}
}
