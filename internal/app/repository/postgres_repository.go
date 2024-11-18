package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx"
)

type PostgresRepository struct {
	connection *pgx.Conn
}

func NewPostgresRepository(c *pgx.Conn) (*PostgresRepository, error) {
	return &PostgresRepository{connection: c}, nil
}

func (r *PostgresRepository) Save(shortURL, originalURL string) error {
	query := `INSERT INTO urls (short_url, original_url) VALUES ($1, $2) ON CONFLICT (short_url) DO NOTHING`
	_, err := r.connection.Exec(query, shortURL, originalURL)
	if err != nil {
		return fmt.Errorf("failed to save URL: %w", err)
	}
	return nil
}

func (r *PostgresRepository) Find(shortURL string) (string, bool) {
	query := `SELECT original_url FROM urls WHERE short_url = $1`
	row := r.connection.QueryRow(query, shortURL)

	var originalURL string
	err := row.Scan(&originalURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false
	}
	if err != nil {
		return "", false
	}
	return originalURL, true
}

func (r *PostgresRepository) All() map[string]string {
	query := `SELECT short_url, original_url FROM urls`
	rows, err := r.connection.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var shortURL, originalURL string
		if err := rows.Scan(&shortURL, &originalURL); err == nil {
			result[shortURL] = originalURL
		}
	}

	return result
}

func (r *PostgresRepository) Close() error {
	err := r.connection.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}

func (r *PostgresRepository) Pint(ctx context.Context) error {
	return r.connection.Ping(ctx)
}
