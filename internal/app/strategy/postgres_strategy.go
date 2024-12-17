package strategy

import (
	"context"
	"fmt"
	"github.com/Gerfey/shortener/internal/app/repository"
	"github.com/Gerfey/shortener/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStrategy struct {
	DSN  string
	pool *pgxpool.Pool
}

func NewPostgresStrategy(dsn string) *PostgresStrategy {
	return &PostgresStrategy{DSN: dsn}
}

func (s *PostgresStrategy) Initialize() (models.Repository, error) {
	config, err := pgxpool.ParseConfig(s.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	config.MaxConns = 50
	config.MinConns = 10

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	s.pool = pool

	return repository.NewPostgresRepository(pool)
}

func (s *PostgresStrategy) Close() error {
	if s.pool == nil {
		return nil
	}

	s.pool.Close()
	return nil
}
