package strategy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresStrategy_Initialize(t *testing.T) {
	tests := []struct {
		name          string
		dsn           string
		expectedError bool
	}{
		{
			name:          "Invalid DSN",
			dsn:           "invalid-dsn",
			expectedError: true,
		},
		{
			name:          "Empty DSN",
			dsn:           "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := NewPostgresStrategy(tt.dsn)
			repo, err := strategy.Initialize(context.Background())

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, repo)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, repo)
			}
		})
	}
}

func TestPostgresStrategy_Close(t *testing.T) {
	strategy := &PostgresStrategy{}
	err := strategy.Close()
	assert.NoError(t, err)
}
