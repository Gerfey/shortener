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
