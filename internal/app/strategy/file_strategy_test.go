package strategy

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestFileStrategy_Initialize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	strategy := NewFileStrategy("testdata/url_store.json")

	repo, err := strategy.Initialize()
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	err = strategy.Close()
	assert.NoError(t, err)
}
