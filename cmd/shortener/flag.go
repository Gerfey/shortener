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
	TrustedSubnet   string `json:"trusted_subnet"`
}

// Flags содержит флаги командной строки
type Flags struct {
	FlagServerRunAddress       string
	FlagServerShortenerAddress string
	FlagGRPCRunAddress         string
	FlagDefaultFilePath        string
	FlagDefaultDatabaseDSN     string
	FlagEnableHTTPS            bool
	FlagConfigFile             string
	FlagTrustedSubnet          string
}

func parseFlags(args []string) Flags {
	const (
		defaultServerAddress = ":8081"
		defaultBaseURL       = "http://localhost:8080"
		defaultGRPCAddress   = ":50051"
		httpsServerAddress   = ":443"
	)

	var flagServerRunAddress, flagServerShortenerAddress, flagGRPCRunAddress, flagDefaultFilePath, flagDefaultDatabaseDSN, flagConfigFile, flagTrustedSubnet string
	var flagEnableHTTPS bool

	envConfigFile := os.Getenv("CONFIG")
	envTrustedSubnet := os.Getenv("TRUSTED_SUBNET")

	fs := flag.NewFlagSet("shortener", flag.ExitOnError)

	fs.StringVar(&flagServerRunAddress, "a", defaultServerAddress, "Run server address and port")
	fs.StringVar(&flagServerShortenerAddress, "b", defaultBaseURL, "Base URL for shortened URLs")
	fs.StringVar(&flagGRPCRunAddress, "g", defaultGRPCAddress, "Run gRPC server address and port")
	fs.StringVar(&flagDefaultFilePath, "f", "", "Path to the file where URLs are stored")
	fs.StringVar(&flagDefaultDatabaseDSN, "d", "", "Database connection DSN")
	fs.BoolVar(&flagEnableHTTPS, "s", false, "Enable HTTPS")
	fs.StringVar(&flagConfigFile, "c", "", "Path to configuration file")
	fs.StringVar(&flagConfigFile, "config", "", "Path to configuration file (shorthand for -c)")
	fs.StringVar(&flagTrustedSubnet, "t", "", "Trusted subnet in CIDR notation for stats access")

	_ = fs.Parse(args)

	var configServerAddress, configBaseURL, configFileStoragePath, configDatabaseDSN, configTrustedSubnet string
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
				configTrustedSubnet = config.TrustedSubnet
			}
		}
	}

	envServerAddress := os.Getenv("SERVER_ADDRESS")
	envBaseURL := os.Getenv("BASE_URL")
	envFilePath := os.Getenv("FILE_STORAGE_PATH")
	envDatabaseDSN := os.Getenv("DATABASE_DSN")
	envEnableHTTPS := os.Getenv("ENABLE_HTTPS") == "true"

	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		flagServerRunAddress = envServerAddress
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		flagServerShortenerAddress = envBaseURL
	}

	if envGRPCAddress := os.Getenv("GRPC_ADDRESS"); envGRPCAddress != "" {
		flagGRPCRunAddress = envGRPCAddress
	}

	serverRunAddress := cmp.Or(envServerAddress, configServerAddress, flagServerRunAddress, defaultServerAddress)
	baseURL := cmp.Or(envBaseURL, configBaseURL, flagServerShortenerAddress, defaultBaseURL)
	fileStoragePath := cmp.Or(envFilePath, configFileStoragePath, flagDefaultFilePath)
	databaseDSN := cmp.Or(envDatabaseDSN, configDatabaseDSN, flagDefaultDatabaseDSN)
	enableHTTPS := cmp.Or(envEnableHTTPS, configEnableHTTPS, flagEnableHTTPS)
	trustedSubnet := cmp.Or(envTrustedSubnet, configTrustedSubnet, flagTrustedSubnet)

	if enableHTTPS && (serverRunAddress == defaultServerAddress) {
		serverRunAddress = httpsServerAddress
	}

	return Flags{
		FlagServerRunAddress:       serverRunAddress,
		FlagServerShortenerAddress: baseURL,
		FlagGRPCRunAddress:         flagGRPCRunAddress,
		FlagDefaultFilePath:        fileStoragePath,
		FlagDefaultDatabaseDSN:     databaseDSN,
		FlagEnableHTTPS:            enableHTTPS,
		FlagConfigFile:             configPath,
		FlagTrustedSubnet:          trustedSubnet,
	}
}
