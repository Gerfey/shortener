package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDatabase_Success(t *testing.T) {
	db, err := NewDatabase("postgresql://user:password@localhost:5432/testdb")
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func TestNewDatabase_Failure(t *testing.T) {
	_, err := NewDatabase("invalid_dsn")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse DSN")
}
