package repository

import (
	"context"
	"github.com/Gerfey/shortener/internal/models"
	"github.com/jackc/pgx/v5"
	pgxmock "github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestPostgresRepository_All(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &PostgresRepository{pool: mock}

	rows := mock.NewRows([]string{"short_url", "original_url"}).
		AddRow("abc123", "https://example.com").
		AddRow("def456", "https://example.org")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT short_url, original_url FROM urls`)).
		WillReturnRows(rows)

	urls := repo.All(context.Background())

	assert.Equal(t, map[string]string{
		"abc123": "https://example.com",
		"def456": "https://example.org",
	}, urls)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_Find(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &PostgresRepository{pool: mock}

	t.Run("URL Found", func(t *testing.T) {
		rows := mock.NewRows([]string{"original_url", "is_deleted"}).
			AddRow("https://example.com", false)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT original_url, is_deleted FROM urls WHERE short_url = $1`)).
			WithArgs("abc123").
			WillReturnRows(rows)

		url, found, isDeleted := repo.Find(context.Background(), "abc123")
		assert.True(t, found)
		assert.False(t, isDeleted)
		assert.Equal(t, "https://example.com", url)
	})

	t.Run("URL Not Found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT original_url, is_deleted FROM urls WHERE short_url = $1`)).
			WithArgs("notfound").
			WillReturnError(pgx.ErrNoRows)

		url, found, isDeleted := repo.Find(context.Background(), "notfound")
		assert.False(t, found)
		assert.False(t, isDeleted)
		assert.Empty(t, url)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_FindShortURL(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &PostgresRepository{pool: mock}

	t.Run("URL Found", func(t *testing.T) {
		rows := mock.NewRows([]string{"short_url"}).
			AddRow("abc123")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT short_url FROM urls WHERE original_url = $1`)).
			WithArgs("https://example.com").
			WillReturnRows(rows)

		shortURL, err := repo.FindShortURL(context.Background(), "https://example.com")
		assert.NoError(t, err)
		assert.Equal(t, "abc123", shortURL)
	})

	t.Run("URL Not Found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT short_url FROM urls WHERE original_url = $1`)).
			WithArgs("https://notfound.com").
			WillReturnError(pgx.ErrNoRows)

		shortURL, err := repo.FindShortURL(context.Background(), "https://notfound.com")
		assert.Error(t, err)
		assert.Empty(t, shortURL)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_Save(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &PostgresRepository{pool: mock}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO urls (short_url, original_url, user_id) VALUES ($1, $2, $3) ON CONFLICT (short_url) DO NOTHING`)).
		WithArgs("abc123", "https://example.com", "user1").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	shortURL, err := repo.Save(context.Background(), "abc123", "https://example.com", "user1")
	assert.NoError(t, err)
	assert.Equal(t, "abc123", shortURL)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_SaveBatch(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &PostgresRepository{pool: mock}

	urls := map[string]string{
		"abc123": "https://example.com",
		"def456": "https://example.org",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO urls (short_url, original_url, user_id) VALUES ($1, $2, $3) ON CONFLICT (short_url) DO NOTHING`)).
		WithArgs("abc123", "https://example.com", "user1").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO urls (short_url, original_url, user_id) VALUES ($1, $2, $3) ON CONFLICT (short_url) DO NOTHING`)).
		WithArgs("def456", "https://example.org", "user1").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	err = repo.SaveBatch(context.Background(), urls, "user1")
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_GetUserURLs(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &PostgresRepository{pool: mock}

	rows := mock.NewRows([]string{"short_url", "original_url"}).
		AddRow("abc123", "https://example.com").
		AddRow("def456", "https://example.org")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT short_url, original_url FROM urls WHERE user_id = $1`)).
		WithArgs("user1").
		WillReturnRows(rows)

	urls, err := repo.GetUserURLs(context.Background(), "user1")
	assert.NoError(t, err)
	assert.Equal(t, []models.URLPair{
		{ShortURL: "abc123", OriginalURL: "https://example.com"},
		{ShortURL: "def456", OriginalURL: "https://example.org"},
	}, urls)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_DeleteUserURLsBatch(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &PostgresRepository{pool: mock}

	shortURLs := []string{"abc123", "def456"}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE urls SET is_deleted = true WHERE short_url = $1 AND user_id = $2`)).
		WithArgs("abc123", "user1").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE urls SET is_deleted = true WHERE short_url = $1 AND user_id = $2`)).
		WithArgs("def456", "user1").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()

	err = repo.DeleteUserURLsBatch(context.Background(), shortURLs, "user1")
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_Ping(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &PostgresRepository{pool: mock}

	mock.ExpectPing()

	err = repo.Ping(context.Background())
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepository_Close(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)

	repo := &PostgresRepository{pool: mock}
	assert.NoError(t, repo.Close())
}
