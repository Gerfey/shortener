package strategy

import (
	"fmt"
	"github.com/Gerfey/shortener/internal/app/repository"
	"github.com/Gerfey/shortener/internal/models"
	"github.com/jackc/pgx"
)

type PostgresStrategy struct {
	DSN        string
	connection *pgx.Conn
}

func NewPostgresStrategy(dsn string) *PostgresStrategy {
	return &PostgresStrategy{DSN: dsn}
}

func (s *PostgresStrategy) Initialize() (models.Repository, error) {
	config, err := pgx.ParseConnectionString(s.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	conn, err := pgx.Connect(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	s.connection = conn

	return repository.NewPostgresRepository(conn)
}

func (s *PostgresStrategy) Close() error {
	err := s.connection.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}
