package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFlags_WithConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configData := Config{
		ServerAddress:   "localhost:8081",
		BaseURL:         "http://localhost:8081",
		FileStoragePath: "/tmp/short-url-db.json",
		DatabaseDSN:     "postgres://user:password@localhost:5432/shortener",
		EnableHTTPS:     true,
	}

	configJSON, err := json.Marshal(configData)
	require.NoError(t, err)

	err = os.WriteFile(configPath, configJSON, 0644)
	require.NoError(t, err)

	flags := parseFlags([]string{"-c=" + configPath})

	assert.Equal(t, "localhost:8081", flags.FlagServerRunAddress)
	assert.Equal(t, "http://localhost:8081", flags.FlagServerShortenerAddress)
	assert.Equal(t, "/tmp/short-url-db.json", flags.FlagDefaultFilePath)
	assert.Equal(t, "postgres://user:password@localhost:5432/shortener", flags.FlagDefaultDatabaseDSN)
	assert.True(t, flags.FlagEnableHTTPS)
}

func TestParseFlags_ConfigPriority(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configData := Config{
		ServerAddress:   "localhost:8081",
		BaseURL:         "http://localhost:8081",
		FileStoragePath: "/tmp/short-url-db.json",
		DatabaseDSN:     "postgres://user:password@localhost:5432/shortener",
		EnableHTTPS:     true,
	}

	configJSON, err := json.Marshal(configData)
	require.NoError(t, err)

	err = os.WriteFile(configPath, configJSON, 0644)
	require.NoError(t, err)

	oldServerAddr := os.Getenv("SERVER_ADDRESS")
	oldBaseURL := os.Getenv("BASE_URL")
	defer func() {
		_ = os.Setenv("SERVER_ADDRESS", oldServerAddr)
		_ = os.Setenv("BASE_URL", oldBaseURL)
	}()

	_ = os.Setenv("SERVER_ADDRESS", ":9090")
	_ = os.Setenv("BASE_URL", "http://example.com")

	flags := parseFlags([]string{
		"-c=" + configPath,
	})

	assert.Equal(t, ":9090", flags.FlagServerRunAddress)
	assert.Equal(t, "http://example.com", flags.FlagServerShortenerAddress)
	assert.Equal(t, "/tmp/short-url-db.json", flags.FlagDefaultFilePath)
	assert.Equal(t, "postgres://user:password@localhost:5432/shortener", flags.FlagDefaultDatabaseDSN)
	assert.True(t, flags.FlagEnableHTTPS)
}

func TestParseFlags_ConfigEnv(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configData := Config{
		ServerAddress:   "localhost:8081",
		BaseURL:         "http://localhost:8081",
		FileStoragePath: "/tmp/short-url-db.json",
		DatabaseDSN:     "postgres://user:password@localhost:5432/shortener",
		EnableHTTPS:     true,
	}

	configJSON, err := json.Marshal(configData)
	require.NoError(t, err)

	err = os.WriteFile(configPath, configJSON, 0644)
	require.NoError(t, err)

	_ = os.Setenv("CONFIG", configPath)
	defer func() {
		_ = os.Unsetenv("CONFIG")
	}()

	flags := parseFlags([]string{})

	assert.Equal(t, "localhost:8081", flags.FlagServerRunAddress)
	assert.Equal(t, "http://localhost:8081", flags.FlagServerShortenerAddress)
	assert.Equal(t, "/tmp/short-url-db.json", flags.FlagDefaultFilePath)
	assert.Equal(t, "postgres://user:password@localhost:5432/shortener", flags.FlagDefaultDatabaseDSN)
	assert.True(t, flags.FlagEnableHTTPS)
}
