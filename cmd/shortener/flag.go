package main

import (
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
	CertFile        string `json:"cert_file,omitempty"`
	KeyFile         string `json:"key_file,omitempty"`
}

// Flags содержит флаги командной строки
type Flags struct {
	FlagServerRunAddress       string
	FlagServerShortenerAddress string
	FlagDefaultFilePath        string
	FlagDefaultDatabaseDSN     string
	FlagEnableHTTPS            bool
	FlagCertFile               string
	FlagKeyFile                string
	FlagConfigFile             string
}

func parseFlags(args []string) Flags {
	var flagServerRunAddress, flagServerShortenerAddress, flagDefaultFilePath, flagDefaultDatabaseDSN, flagCertFile, flagKeyFile, flagConfigFile string
	var flagEnableHTTPS bool

	fs := flag.NewFlagSet("shortener", flag.ExitOnError)

	fs.StringVar(&flagServerRunAddress, "a", ":8080", "Run server address and port")
	fs.StringVar(&flagServerShortenerAddress, "b", "http://localhost:8080", "Run server address and port")
	fs.StringVar(&flagDefaultFilePath, "f", "", "Path to the file where URLs are stored")
	fs.StringVar(&flagDefaultDatabaseDSN, "d", "", "Database connection DSN")
	fs.BoolVar(&flagEnableHTTPS, "s", false, "Enable HTTPS")
	fs.StringVar(&flagCertFile, "cert", "server.crt", "Path to TLS certificate file")
	fs.StringVar(&flagKeyFile, "key", "server.key", "Path to TLS key file")
	fs.StringVar(&flagConfigFile, "c", "", "Path to configuration file")
	fs.StringVar(&flagConfigFile, "config", "", "Path to configuration file (shorthand for -c)")

	if err := fs.Parse(args); err != nil {
		return Flags{}
	}

	if envConfigFile := os.Getenv("CONFIG"); envConfigFile != "" {
		flagConfigFile = envConfigFile
	}

	if flagConfigFile != "" {
		configData, err := os.ReadFile(flagConfigFile)
		if err == nil {
			var config Config
			if err := json.Unmarshal(configData, &config); err == nil {
				if flagServerRunAddress == ":8080" && config.ServerAddress != "" {
					flagServerRunAddress = config.ServerAddress
				}
				if flagServerShortenerAddress == "http://localhost:8080" && config.BaseURL != "" {
					flagServerShortenerAddress = config.BaseURL
				}
				if flagDefaultFilePath == "" && config.FileStoragePath != "" {
					flagDefaultFilePath = config.FileStoragePath
				}
				if flagDefaultDatabaseDSN == "" && config.DatabaseDSN != "" {
					flagDefaultDatabaseDSN = config.DatabaseDSN
				}
				if !flagEnableHTTPS && config.EnableHTTPS {
					flagEnableHTTPS = config.EnableHTTPS
				}
				if flagCertFile == "server.crt" && config.CertFile != "" {
					flagCertFile = config.CertFile
				}
				if flagKeyFile == "server.key" && config.KeyFile != "" {
					flagKeyFile = config.KeyFile
				}
			}
		}
	}

	if flagEnableHTTPS {
		if flagServerRunAddress == ":8080" {
			flagServerRunAddress = ":443"
		}
	}

	if envServerRunAddress := os.Getenv("SERVER_ADDRESS"); envServerRunAddress != "" {
		flagServerRunAddress = envServerRunAddress
	}

	if envServerShortenerAddress := os.Getenv("BASE_URL"); envServerShortenerAddress != "" {
		flagServerShortenerAddress = envServerShortenerAddress
	}

	if envDefaultFilePath := os.Getenv("FILE_STORAGE_PATH"); envDefaultFilePath != "" {
		flagDefaultFilePath = envDefaultFilePath
	}

	if envDefaultDatabaseDSN := os.Getenv("DATABASE_DSN"); envDefaultDatabaseDSN != "" {
		flagDefaultDatabaseDSN = envDefaultDatabaseDSN
	}

	if os.Getenv("ENABLE_HTTPS") == "true" {
		flagEnableHTTPS = true
	}

	if envCertFile := os.Getenv("CERT_FILE"); envCertFile != "" {
		flagCertFile = envCertFile
	}

	if envKeyFile := os.Getenv("KEY_FILE"); envKeyFile != "" {
		flagKeyFile = envKeyFile
	}

	return Flags{flagServerRunAddress, flagServerShortenerAddress, flagDefaultFilePath, flagDefaultDatabaseDSN, flagEnableHTTPS, flagCertFile, flagKeyFile, flagConfigFile}
}
