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
			s := &URLMemoryRepository{
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
