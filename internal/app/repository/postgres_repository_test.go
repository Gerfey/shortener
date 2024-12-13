package repository

import (
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDatabase_Success(t *testing.T) {
	config, _ := pgx.ParseConnectionString("postgresql://user:password@localhost:5432/testdb")
	conn, _ := pgx.Connect(config)

	db, err := NewPostgresRepository(conn)
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func TestNewDatabase_Failure(t *testing.T) {
	_, err := pgx.ParseConnectionString("invalid_dsn")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid dsn")
}
