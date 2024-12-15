package repository

import (
	"context"
	"fmt"
	"github.com/Gerfey/shortener/internal/models"
	"github.com/jackc/pgx"
)

type PostgresRepository struct {
	conn *pgx.Conn
}

func NewPostgresRepository(conn *pgx.Conn) (*PostgresRepository, error) {
	repo := &PostgresRepository{
		conn: conn,
	}

	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			short_url VARCHAR(255) UNIQUE NOT NULL,
			original_url TEXT NOT NULL,
			user_id VARCHAR(255),
			is_deleted BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return repo, nil
}

func (r *PostgresRepository) All() map[string]string {
	urls := make(map[string]string)
	rows, err := r.conn.Query("SELECT short_url, original_url FROM urls")
	if err != nil {
		return urls
	}
	defer rows.Close()

	for rows.Next() {
		var shortURL, originalURL string
		if err := rows.Scan(&shortURL, &originalURL); err != nil {
			continue
		}
		urls[shortURL] = originalURL
	}

	return urls
}

func (r *PostgresRepository) Find(key string) (string, bool, bool) {
	var originalURL string
	var isDeleted bool
	
	err := r.conn.QueryRow("SELECT original_url, is_deleted FROM urls WHERE short_url = $1", key).Scan(&originalURL, &isDeleted)
	if err != nil {
		return "", false, false
	}
	
	return originalURL, true, isDeleted
}

func (r *PostgresRepository) FindShortURL(originalURL string) (string, error) {
	var shortURL string
	err := r.conn.QueryRow("SELECT short_url FROM urls WHERE original_url = $1", originalURL).Scan(&shortURL)
	if err != nil {
		return "", fmt.Errorf("original URL not found")
	}
	return shortURL, nil
}

func (r *PostgresRepository) Save(key, value string, userID string) (string, error) {
	_, err := r.conn.Exec(`
		INSERT INTO urls (short_url, original_url, user_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (short_url) DO NOTHING
	`, key, value, userID)
	if err != nil {
		return "", fmt.Errorf("failed to save URL: %w", err)
	}
	return key, nil
}

func (r *PostgresRepository) SaveBatch(urls map[string]string, userID string) error {
	tx, err := r.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx *pgx.Tx) {
		_ = tx.Rollback()
	}(tx)

	for shortURL, originalURL := range urls {
		_, err = tx.Exec(`
			INSERT INTO urls (short_url, original_url, user_id)
			VALUES ($1, $2, $3)
			ON CONFLICT (short_url) DO NOTHING
		`, shortURL, originalURL, userID)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) DeleteUserURLsBatch(shortURLs []string, userID string) error {
	query := `
		UPDATE urls 
		SET is_deleted = TRUE 
		WHERE short_url = ANY($1) AND user_id = $2
	`
	
	_, err := r.conn.Exec(query, shortURLs, userID)
	if err != nil {
		return fmt.Errorf("failed to mark URLs as deleted: %w", err)
	}
	
	return nil
}

func (r *PostgresRepository) GetUserURLs(userID string) ([]models.URLPair, error) {
	rows, err := r.conn.Query("SELECT short_url, original_url FROM urls WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user URLs: %w", err)
	}
	defer rows.Close()

	var userURLs []models.URLPair
	for rows.Next() {
		var pair models.URLPair
		if err := rows.Scan(&pair.ShortURL, &pair.OriginalURL); err != nil {
			return nil, fmt.Errorf("failed to scan URL pair: %w", err)
		}
		userURLs = append(userURLs, pair)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return userURLs, nil
}

func (r *PostgresRepository) Ping() error {
	ctx := context.Background()
	return r.conn.Ping(ctx)
}

func (r *PostgresRepository) Close() error {
	return r.conn.Close()
}
