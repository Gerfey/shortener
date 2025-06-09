package main

import (
	"cmp"
	"encoding/json"
	"flag"
	"os"
)

// Config структура для конфигурационного файла JSON
type Config struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
}

// Flags содержит флаги командной строки
type Flags struct {
	FlagServerRunAddress       string
	FlagServerShortenerAddress string
	FlagDefaultFilePath        string
	FlagDefaultDatabaseDSN     string
	FlagEnableHTTPS            bool
	FlagConfigFile             string
}

func parseFlags(args []string) Flags {
	const (
		defaultServerAddress = ":8080"
		defaultBaseURL       = "http://localhost:8080"
		httpsServerAddress   = ":443"
	)

	var flagServerRunAddress, flagServerShortenerAddress, flagDefaultFilePath, flagDefaultDatabaseDSN, flagConfigFile string
	var flagEnableHTTPS bool

	envConfigFile := os.Getenv("CONFIG")

	fs := flag.NewFlagSet("shortener", flag.ExitOnError)

	fs.StringVar(&flagServerRunAddress, "a", defaultServerAddress, "Run server address and port")
	fs.StringVar(&flagServerShortenerAddress, "b", defaultBaseURL, "Base URL for shortened URLs")
	fs.StringVar(&flagDefaultFilePath, "f", "", "Path to the file where URLs are stored")
	fs.StringVar(&flagDefaultDatabaseDSN, "d", "", "Database connection DSN")
	fs.BoolVar(&flagEnableHTTPS, "s", false, "Enable HTTPS")
	fs.StringVar(&flagConfigFile, "c", "", "Path to configuration file")
	fs.StringVar(&flagConfigFile, "config", "", "Path to configuration file (shorthand for -c)")

	_ = fs.Parse(args)

	var configServerAddress, configBaseURL, configFileStoragePath, configDatabaseDSN string
	var configEnableHTTPS bool

	configPath := cmp.Or(envConfigFile, flagConfigFile)

	if configPath != "" {
		var config Config
		configFile, err := os.ReadFile(configPath)
		if err == nil {
			if err = json.Unmarshal(configFile, &config); err == nil {
				configServerAddress = config.ServerAddress
				configBaseURL = config.BaseURL
				configFileStoragePath = config.FileStoragePath
				configDatabaseDSN = config.DatabaseDSN
				configEnableHTTPS = config.EnableHTTPS
			}
		}
	}

	envServerAddress := os.Getenv("SERVER_ADDRESS")
	envBaseURL := os.Getenv("BASE_URL")
	envFilePath := os.Getenv("FILE_STORAGE_PATH")
	envDatabaseDSN := os.Getenv("DATABASE_DSN")
	envEnableHTTPS := os.Getenv("ENABLE_HTTPS") == "true"

	serverRunAddress := cmp.Or(envServerAddress, configServerAddress, flagServerRunAddress, defaultServerAddress)
	serverShortenerAddress := cmp.Or(envBaseURL, configBaseURL, flagServerShortenerAddress, defaultBaseURL)
	defaultFilePath := cmp.Or(envFilePath, configFileStoragePath, flagDefaultFilePath)
	defaultDatabaseDSN := cmp.Or(envDatabaseDSN, configDatabaseDSN, flagDefaultDatabaseDSN)
	enableHTTPS := cmp.Or(envEnableHTTPS, configEnableHTTPS, flagEnableHTTPS)

	if enableHTTPS && (serverRunAddress == defaultServerAddress) {
		serverRunAddress = httpsServerAddress
	}

	return Flags{
		FlagServerRunAddress:       serverRunAddress,
		FlagServerShortenerAddress: serverShortenerAddress,
		FlagDefaultFilePath:        defaultFilePath,
		FlagDefaultDatabaseDSN:     defaultDatabaseDSN,
		FlagEnableHTTPS:            enableHTTPS,
		FlagConfigFile:             configPath,
	}
}
