package app

import (
	"bytes"
	"encoding/json"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/app/strategy"
	"github.com/Gerfey/shortener/internal/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestNewShortenerApp(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "url_store.json")

	tests := []struct {
		name          string
		storageType   string
		filePath      string
		expectedError bool
	}{
		{
			name:          "Memory storage",
			storageType:   "memory",
			filePath:      "",
			expectedError: false,
		},
		{
			name:          "File storage",
			storageType:   "file",
			filePath:      tmpFile,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("STORAGE_TYPE", tt.storageType)
			os.Setenv("FILE_STORAGE_PATH", tt.filePath)
			defer os.Unsetenv("STORAGE_TYPE")
			defer os.Unsetenv("FILE_STORAGE_PATH")

			config := settings.NewSettings(settings.ServerSettings{
				ServerShortenerAddress: "http://localhost:8080",
			})

			var stg models.StorageStrategy
			if tt.storageType == "memory" {
				stg = strategy.NewMemoryStrategy()
			} else {
				stg = strategy.NewFileStrategy(tt.filePath)
			}

			app, err := NewShortenerApp(config, stg)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, app)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, app)
			}
		})
	}
}

func TestShortenerApp_Run(t *testing.T) {
	config := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
	})

	stg := strategy.NewMemoryStrategy()
	app, err := NewShortenerApp(config, stg)
	assert.NoError(t, err)

	app.configureRouter()
	server := httptest.NewServer(app.router)
	defer server.Close()

	tests := []struct {
		name           string
		endpoint       string
		method         string
		body           interface{}
		expectedStatus int
		withCookie     bool
	}{
		{
			name:     "Shorten JSON URL",
			endpoint: "/api/shorten",
			method:   http.MethodPost,
			body: map[string]string{
				"url": "https://example.com",
			},
			expectedStatus: http.StatusCreated,
			withCookie:     true,
		},
		{
			name:     "Shorten Batch URLs",
			endpoint: "/api/shorten/batch",
			method:   http.MethodPost,
			body: []map[string]string{
				{
					"correlation_id": "1",
					"original_url":   "https://example1.com",
				},
				{
					"correlation_id": "2",
					"original_url":   "https://example2.com",
				},
			},
			expectedStatus: http.StatusCreated,
			withCookie:     true,
		},
		{
			name:           "Ping",
			endpoint:       "/ping",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			withCookie:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBytes []byte
			if tt.body != nil {
				reqBytes, _ = json.Marshal(tt.body)
			}

			req, err := http.NewRequest(tt.method, server.URL+tt.endpoint, bytes.NewReader(reqBytes))
			assert.NoError(t, err)

			if tt.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}

			if tt.withCookie {
				req.AddCookie(&http.Cookie{
					Name:  "user_id",
					Value: "test_user",
				})
			}

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if resp.StatusCode == http.StatusCreated {
				var respBody interface{}
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)
			}
		})
	}
}
