package strategy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileStrategy_Initialize(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "url_store.json")

	testData := `{
		"abc123": {"original_url": "https://example.com", "user_id": "user1"},
		"def456": {"original_url": "https://google.com", "user_id": "user1"}
	}`

	err := os.WriteFile(tmpFile, []byte(testData), 0644)
	assert.NoError(t, err)

	strategy := NewFileStrategy(tmpFile)

	repo, err := strategy.Initialize()
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	url, exists := repo.Find("abc123")
	assert.True(t, exists)
	assert.Equal(t, "https://example.com", url)
}

func TestFileStrategy_Close(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "url_store.json")

	strategy := NewFileStrategy(tmpFile)

	repo, err := strategy.Initialize()
	assert.NoError(t, err)

	_, err = repo.Save("abc123", "https://example.com", "user1")
	assert.NoError(t, err)

	err = strategy.Close()
	assert.NoError(t, err)

	data, err := os.ReadFile(tmpFile)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "https://example.com")
	assert.Contains(t, string(data), "user1")
}

func TestFileStrategy_InitializeError(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "nonexistent_dir", "url_store.json")

	strategy := NewFileStrategy(tmpFile)

	repo, err := strategy.Initialize()
	assert.Error(t, err)
	assert.Nil(t, repo)
}

func TestFileStrategy_CloseError(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "url_store.json")

	strategy := NewFileStrategy(tmpFile)

	err := strategy.Close()
	assert.NoError(t, err)
}
