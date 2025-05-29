package main

import (
	"flag"
	"os"
)

// Flags содержит флаги командной строки
type Flags struct {
	FlagServerRunAddress       string
	FlagServerShortenerAddress string
	FlagDefaultFilePath        string
	FlagDefaultDatabaseDSN     string
	FlagEnableHTTPS            bool
	FlagCertFile               string
	FlagKeyFile                string
}

func parseFlags(args []string) Flags {
	var flagServerRunAddress, flagServerShortenerAddress, flagDefaultFilePath, flagDefaultDatabaseDSN, flagCertFile, flagKeyFile string
	var flagEnableHTTPS bool

	fs := flag.NewFlagSet("shortener", flag.ExitOnError)

	fs.StringVar(&flagServerRunAddress, "a", ":8080", "Run server address and port")
	fs.StringVar(&flagServerShortenerAddress, "b", "http://localhost:8080", "Run server address and port")
	fs.StringVar(&flagDefaultFilePath, "f", "", "Path to the file where URLs are stored")
	fs.StringVar(&flagDefaultDatabaseDSN, "d", "", "Database connection DSN")
	fs.BoolVar(&flagEnableHTTPS, "s", false, "Enable HTTPS")
	fs.StringVar(&flagCertFile, "cert", "server.crt", "Path to TLS certificate file")
	fs.StringVar(&flagKeyFile, "key", "server.key", "Path to TLS key file")

	err := fs.Parse(args)
	if err != nil {
		return Flags{}
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

	return Flags{flagServerRunAddress, flagServerShortenerAddress, flagDefaultFilePath, flagDefaultDatabaseDSN, flagEnableHTTPS, flagCertFile, flagKeyFile}
}
