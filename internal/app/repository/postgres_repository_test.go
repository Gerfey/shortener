package repository

import (
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDatabase_Success(t *testing.T) {
	t.Skip("Skipping database test - requires actual database connection")
}

func TestNewDatabase_Failure(t *testing.T) {
	_, err := pgx.ParseConnectionString("invalid_dsn")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid dsn")
}
