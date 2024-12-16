package strategy

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryStrategy_Initialize(t *testing.T) {
	strategy := NewMemoryStrategy()

	repo, err := strategy.Initialize(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	err = strategy.Close()
	assert.NoError(t, err)
}
