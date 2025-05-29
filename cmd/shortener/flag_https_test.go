package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags_WithHTTPSFlag(t *testing.T) {
	flags := parseFlags([]string{"-s"})

	assert.True(t, flags.FlagEnableHTTPS)
	assert.Equal(t, ":443", flags.FlagServerRunAddress)
}

func TestParseFlags_WithHTTPSEnv(t *testing.T) {
	flags := parseFlags([]string{"-s"})

	assert.True(t, flags.FlagEnableHTTPS)
	assert.Equal(t, ":443", flags.FlagServerRunAddress)
}

func TestParseFlags_WithHTTPSAndServerAddrEnv(t *testing.T) {
	oldServerAddr := os.Getenv("SERVER_ADDRESS")
	oldEnableHTTPS := os.Getenv("ENABLE_HTTPS")
	defer func() {
		_ = os.Setenv("SERVER_ADDRESS", oldServerAddr)
		_ = os.Setenv("ENABLE_HTTPS", oldEnableHTTPS)
	}()

	_ = os.Setenv("ENABLE_HTTPS", "true")
	_ = os.Setenv("SERVER_ADDRESS", ":8888")

	flags := parseFlags([]string{})

	assert.True(t, flags.FlagEnableHTTPS)
	assert.Equal(t, ":8888", flags.FlagServerRunAddress)
}

func TestParseFlags_WithCertKeyFlags(t *testing.T) {
	flags := parseFlags([]string{"-cert=custom.crt", "-key=custom.key"})

	assert.Equal(t, "custom.crt", flags.FlagCertFile)
	assert.Equal(t, "custom.key", flags.FlagKeyFile)
}

func TestParseFlags_WithCertKeyEnv(t *testing.T) {
	_ = os.Setenv("CERT_FILE", "env.crt")
	_ = os.Setenv("KEY_FILE", "env.key")
	defer func() {
		_ = os.Unsetenv("CERT_FILE")
		_ = os.Unsetenv("KEY_FILE")
	}()

	flags := parseFlags([]string{})

	assert.Equal(t, "env.crt", flags.FlagCertFile)
	assert.Equal(t, "env.key", flags.FlagKeyFile)
}
