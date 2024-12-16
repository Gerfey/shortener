package repository

import (
	"context"
	"fmt"
	"github.com/Gerfey/shortener/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func NewPostgresRepository(ctx context.Context, pool *pgxpool.Pool) (*PostgresRepository, error) {
	repo := &PostgresRepository{
		pool: pool,
		ctx:  ctx,
	}

	_, err := pool.Exec(ctx, `
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
	rows, err := r.pool.Query(r.ctx, "SELECT short_url, original_url FROM urls")
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

	err := r.pool.QueryRow(r.ctx, "SELECT original_url, is_deleted FROM urls WHERE short_url = $1", key).Scan(&originalURL, &isDeleted)
	if err != nil {
		return "", false, false
	}

	return originalURL, true, isDeleted
}

func (r *PostgresRepository) FindShortURL(originalURL string) (string, error) {
	var shortURL string
	err := r.pool.QueryRow(r.ctx, "SELECT short_url FROM urls WHERE original_url = $1", originalURL).Scan(&shortURL)
	if err != nil {
		return "", fmt.Errorf("original URL not found")
	}
	return shortURL, nil
}

func (r *PostgresRepository) Save(key, value string, userID string) (string, error) {
	_, err := r.pool.Exec(r.ctx, `
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
	tx, err := r.pool.Begin(r.ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(r.ctx); rollbackErr != nil && err == nil {
			err = fmt.Errorf("failed to rollback transaction: %w", rollbackErr)
		}
	}()

	for shortURL, originalURL := range urls {
		_, err = tx.Exec(r.ctx, `
			INSERT INTO urls (short_url, original_url, user_id)
			VALUES ($1, $2, $3)
			ON CONFLICT (short_url) DO NOTHING
		`, shortURL, originalURL, userID)
		if err != nil {
			return fmt.Errorf("failed to save URL in batch: %w", err)
		}
	}

	if err := tx.Commit(r.ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresRepository) DeleteUserURLsBatch(shortURLs []string, userID string) error {
	tx, err := r.pool.Begin(r.ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(r.ctx); rollbackErr != nil && err == nil {
			err = fmt.Errorf("failed to rollback transaction: %w", rollbackErr)
		}
	}()

	for _, shortURL := range shortURLs {
		_, err = tx.Exec(r.ctx, `
			UPDATE urls 
			SET is_deleted = true 
			WHERE short_url = $1 AND user_id = $2
		`, shortURL, userID)
		if err != nil {
			return fmt.Errorf("failed to mark URLs as deleted: %w", err)
		}
	}

	if err := tx.Commit(r.ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetUserURLs(userID string) ([]models.URLPair, error) {
	rows, err := r.pool.Query(r.ctx, `
		SELECT short_url, original_url 
		FROM urls 
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user URLs: %w", err)
	}
	defer rows.Close()

	var urls []models.URLPair
	for rows.Next() {
		var pair models.URLPair
		if err := rows.Scan(&pair.ShortURL, &pair.OriginalURL); err != nil {
			return nil, fmt.Errorf("failed to scan URL pair: %w", err)
		}
		urls = append(urls, pair)
	}

	return urls, nil
}

func (r *PostgresRepository) Ping() error {
	return r.pool.Ping(r.ctx)
}

func (r *PostgresRepository) Close() error {
	if r.pool != nil {
		r.pool.Close()
	}
	return nil
}
