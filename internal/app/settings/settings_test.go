package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerAddress(t *testing.T) {
	testCases := []struct {
		name                           string
		serverSettings                 ServerSettings
		expectedServerAddress          string
		expectedShortenerServerAddress string
	}{
		{
			"Empty",
			ServerSettings{ServerRunAddress: "", ServerShortenerAddress: ""},
			"",
			"",
		},
		{
			"Localhost server address",
			ServerSettings{ServerRunAddress: "localhost", ServerShortenerAddress: "https://localhost"},
			"localhost",
			"https://localhost",
		},
		{
			"Domain name server address",
			ServerSettings{ServerRunAddress: "www.example.com", ServerShortenerAddress: "https://www.example.com"},
			"www.example.com",
			"https://www.example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			settings := NewSettings(tc.serverSettings)

			serverAddress := settings.ServerAddress()
			shortenerServerAddress := settings.ShortenerServerAddress()

			assert.Equal(t, tc.expectedServerAddress, serverAddress)
			assert.Equal(t, tc.expectedShortenerServerAddress, shortenerServerAddress)
		})
	}
}
