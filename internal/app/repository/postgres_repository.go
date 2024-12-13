package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
)

type PostgresRepository struct {
	connection *pgx.Conn
}

func NewPostgresRepository(c *pgx.Conn) (*PostgresRepository, error) {
	return &PostgresRepository{connection: c}, nil
}

func (r *PostgresRepository) FindShortURL(originalURL string) (string, error) {
	var shortURL string
	query := `SELECT short_url FROM urls WHERE original_url = $1`
	err := r.connection.QueryRow(query, originalURL).Scan(&shortURL)
	if err != nil {
		return "", fmt.Errorf("failed to find short URL: %w", err)
	}
	return shortURL, nil
}

func (r *PostgresRepository) SaveBatch(urls map[string]string) error {
	tx, err := r.connection.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx *pgx.Tx) {
		_ = tx.Rollback()
	}(tx)

	for shortURL, originalURL := range urls {
		_, err := tx.Exec("INSERT INTO urls (short_url, original_url) VALUES ($1, $2) ON CONFLICT (original_url) DO NOTHING",
			shortURL, originalURL)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				continue
			}
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresRepository) Save(shortURL, originalURL string) (string, error) {
	query := `INSERT INTO urls (short_url, original_url) VALUES ($1, $2) ON CONFLICT (original_url) DO NOTHING RETURNING short_url`

	var resultShortURL string
	err := r.connection.QueryRow(query, shortURL, originalURL).Scan(&resultShortURL)
	if err != nil {
		if err.Error() == "no rows in result set" {
			query = `SELECT short_url FROM urls WHERE original_url = $1`
			newErr := r.connection.QueryRow(query, originalURL).Scan(&resultShortURL)
			if newErr != nil {
				return "", fmt.Errorf("failed to fetch existing short_url: %w", err)
			}
			return resultShortURL, err
		}
		return "", fmt.Errorf("failed to save URL: %w", err)
	}

	return resultShortURL, nil
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
