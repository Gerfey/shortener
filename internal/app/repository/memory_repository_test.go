package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore_Get(t *testing.T) {
	type store struct {
		data map[string]string
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		store  store
		want   args
		exists bool
	}{
		{
			name:  "Find existing key",
			store: store{data: map[string]string{"key": "value"}},
			want: args{
				key:   "key",
				value: "value",
			},
			exists: true,
		},
		{
			name:  "Find not found key",
			store: store{data: map[string]string{"key": "value"}},
			want: args{
				key:   "key1",
				value: "value1",
			},
			exists: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemoryRepository{
				data: tt.store.data,
			}
			value, exists := s.Find(tt.want.key)
			if tt.exists && value != tt.want.value {
				assert.Equal(t, value, tt.want.value, "Значение в хранилище не совпадает с ожидаемым")
			}
			if exists != tt.exists {
				assert.Equal(t, exists, tt.exists, "Ответ от хранилища о существовании записи не совпадает с ожидаемым")
			}
		})
	}
}

func TestURLMemoryRepository_SaveAndFind(t *testing.T) {
	repo := NewMemoryRepository()

	originalURL := "https://example.com"
	shortID := "s65fg"

	err := repo.Save(shortID, originalURL)
	assert.NoError(t, err)

	url, found := repo.Find(shortID)
	assert.True(t, found)
	assert.Equal(t, originalURL, url)
}

func TestURLMemoryRepository_All(t *testing.T) {
	repo := NewMemoryRepository()

	_ = repo.Save("key1", "https://example1.com")
	_ = repo.Save("key2", "https://example2.com")

	all := repo.All()
	assert.Equal(t, 2, len(all))
	assert.Equal(t, "https://example1.com", all["key1"])
	assert.Equal(t, "https://example2.com", all["key2"])
}
