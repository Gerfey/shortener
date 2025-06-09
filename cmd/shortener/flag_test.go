package main

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlagParsing(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	var aValue = "localhost:8081"
	var bValue = "http://localhost:8082"
	var dValue = "host=localhost port=5432 user=shortener password=shortener dbname=shortener sslmode=disable"

	os.Args = []string{"cmd", "-a=" + aValue, "-b=" + bValue, "-d=" + dValue}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.CommandLine = flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	flag.CommandLine = flag.NewFlagSet(os.Args[2], flag.ExitOnError)

	var a string
	flag.StringVar(&a, "a", "", "")

	var b string
	flag.StringVar(&b, "b", "", "")

	var d string
	flag.StringVar(&d, "d", "", "")

	flag.Parse()

	if a != aValue {
		t.Errorf("expected %v, got %v", aValue, a)
	}

	if b != bValue {
		t.Errorf("expected %v, got %v", bValue, b)
	}

	if d != dValue {
		t.Errorf("expected %v, got %v", dValue, d)
	}
}

func TestParseFlags_Defaults(t *testing.T) {
	flags := parseFlags([]string{})

	assert.Equal(t, ":8081", flags.FlagServerRunAddress)
	assert.Equal(t, "http://localhost:8080", flags.FlagServerShortenerAddress)
	assert.Equal(t, "", flags.FlagDefaultFilePath)
	assert.Equal(t, "", flags.FlagDefaultDatabaseDSN)
}

func TestParseFlags_WithEnv(t *testing.T) {
	_ = os.Setenv("SERVER_ADDRESS", ":9090")
	_ = os.Setenv("BASE_URL", "http://example.com")
	_ = os.Setenv("FILE_STORAGE_PATH", "data.json")
	_ = os.Setenv("DATABASE_DSN", "postgresql://example:example@localhost:5432/example")

	defer func() {
		_ = os.Unsetenv("SERVER_ADDRESS")
		_ = os.Unsetenv("BASE_URL")
		_ = os.Unsetenv("FILE_STORAGE_PATH")
		_ = os.Unsetenv("DATABASE_DSN")
	}()

	flags := parseFlags([]string{})

	assert.Equal(t, ":9090", flags.FlagServerRunAddress)
	assert.Equal(t, "http://example.com", flags.FlagServerShortenerAddress)
	assert.Equal(t, "data.json", flags.FlagDefaultFilePath)
	assert.Equal(t, "postgresql://example:example@localhost:5432/example", flags.FlagDefaultDatabaseDSN)
	assert.False(t, flags.FlagEnableHTTPS)
}
