package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestInitializeStorageStrategy(t *testing.T) {
	tests := []struct {
		name         string
		databaseDSN  string
		filePath     string
		expectedType string
	}{
		{
			name:         "Postgres Strategy",
			databaseDSN:  "postgres://user:pass@localhost:5432/db",
			filePath:     "",
			expectedType: "*strategy.PostgresStrategy",
		},
		{
			name:         "File Strategy",
			databaseDSN:  "",
			filePath:     "/tmp/test.json",
			expectedType: "*strategy.FileStrategy",
		},
		{
			name:         "Memory Strategy",
			databaseDSN:  "",
			filePath:     "",
			expectedType: "*strategy.MemoryStrategy",
		},
		{
			name:         "Postgres Priority",
			databaseDSN:  "postgres://user:pass@localhost:5432/db",
			filePath:     "/tmp/test.json",
			expectedType: "*strategy.PostgresStrategy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origDBDSN := os.Getenv("DATABASE_DSN")
			origFilePath := os.Getenv("FILE_STORAGE_PATH")
			defer func() {
				os.Setenv("DATABASE_DSN", origDBDSN)
				os.Setenv("FILE_STORAGE_PATH", origFilePath)
			}()

			if tt.databaseDSN != "" {
				os.Setenv("DATABASE_DSN", tt.databaseDSN)
			} else {
				os.Unsetenv("DATABASE_DSN")
			}
			if tt.filePath != "" {
				os.Setenv("FILE_STORAGE_PATH", tt.filePath)
			} else {
				os.Unsetenv("FILE_STORAGE_PATH")
			}

			origArgs := os.Args
			defer func() { os.Args = origArgs }()

			os.Args = []string{"cmd"}

			flags := parseFlags(os.Args[1:])

			if tt.databaseDSN != "" {
				assert.Equal(t, tt.databaseDSN, flags.FlagDefaultDatabaseDSN)
			}
			if tt.filePath != "" {
				assert.Equal(t, tt.filePath, flags.FlagDefaultFilePath)
			}
		})
	}
}

func TestMainWithDifferentFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		envVars  map[string]string
		checkEnv bool
	}{
		{
			name: "Command Line Arguments",
			args: []string{
				"-a=:8081",
				"-b=http://localhost:8081",
				"-f=/tmp/urls.json",
				"-d=postgres://test:test@localhost:5432/test",
			},
			checkEnv: false,
		},
		{
			name: "Environment Variables",
			envVars: map[string]string{
				"SERVER_ADDRESS":    ":8082",
				"BASE_URL":          "http://localhost:8082",
				"FILE_STORAGE_PATH": "/tmp/storage.json",
				"DATABASE_DSN":      "postgres://prod:prod@localhost:5432/prod",
			},
			checkEnv: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origArgs := os.Args
			origEnvVars := make(map[string]string)

			if tt.checkEnv {
				for key := range tt.envVars {
					origEnvVars[key] = os.Getenv(key)
				}
			}

			defer func() {
				os.Args = origArgs
				if tt.checkEnv {
					for key, value := range origEnvVars {
						if value != "" {
							os.Setenv(key, value)
						} else {
							os.Unsetenv(key)
						}
					}
				}
			}()

			if tt.checkEnv {
				for key, value := range tt.envVars {
					os.Setenv(key, value)
				}
				os.Args = []string{"cmd"}
			} else {
				os.Args = append([]string{"cmd"}, tt.args...)
			}

			flags := parseFlags(os.Args[1:])

			if tt.checkEnv {
				assert.Equal(t, tt.envVars["SERVER_ADDRESS"], flags.FlagServerRunAddress)
				assert.Equal(t, tt.envVars["BASE_URL"], flags.FlagServerShortenerAddress)
				assert.Equal(t, tt.envVars["FILE_STORAGE_PATH"], flags.FlagDefaultFilePath)
				assert.Equal(t, tt.envVars["DATABASE_DSN"], flags.FlagDefaultDatabaseDSN)
			} else {
				assert.Equal(t, ":8081", flags.FlagServerRunAddress)
				assert.Equal(t, "http://localhost:8081", flags.FlagServerShortenerAddress)
				assert.Equal(t, "/tmp/urls.json", flags.FlagDefaultFilePath)
				assert.Equal(t, "postgres://test:test@localhost:5432/test", flags.FlagDefaultDatabaseDSN)
			}
		})
	}
}

func TestMain(t *testing.T) {
	origArgs := os.Args
	origTestMode := testMode
	defer func() {
		os.Args = origArgs
		testMode = origTestMode
	}()

	testCases := []struct {
		name    string
		args    []string
		envVars map[string]string
	}{
		{
			name: "Default configuration",
			args: []string{"shortener"},
			envVars: map[string]string{
				"SERVER_ADDRESS":    ":8080",
				"BASE_URL":          "http://localhost:8080",
				"FILE_STORAGE_PATH": "",
				"DATABASE_DSN":      "",
			},
		},
		{
			name: "With file storage",
			args: []string{"shortener", "-f", "/tmp/test.json"},
			envVars: map[string]string{
				"SERVER_ADDRESS":    ":8081",
				"BASE_URL":          "http://localhost:8081",
				"FILE_STORAGE_PATH": "/tmp/test.json",
				"DATABASE_DSN":      "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testDoneCh = make(chan struct{})
			testMode = true

			for k, v := range tc.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			os.Args = tc.args

			go main()

			time.Sleep(100 * time.Millisecond)

			close(testDoneCh)

			time.Sleep(50 * time.Millisecond)
		})
	}
}
