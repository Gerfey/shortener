package strategy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStrategy_Initialize(t *testing.T) {
	strategy := NewMemoryStrategy()

	repo, err := strategy.Initialize()
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	err = strategy.Close()
	assert.NoError(t, err)
}
