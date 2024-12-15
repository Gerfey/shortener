package settings

import "time"

type ServerSettings struct {
	ServerRunAddress       string
	ServerShortenerAddress string
	DefaultFilePath        string
	DefaultDatabaseDSN     string
	ShutdownTimeout        time.Duration
}

type Settings struct {
	Server ServerSettings
}

func NewSettings(serverSettings ServerSettings) *Settings {
	return &Settings{
		Server: ServerSettings{
			ServerRunAddress:       serverSettings.ServerRunAddress,
			ServerShortenerAddress: serverSettings.ServerShortenerAddress,
			DefaultFilePath:        serverSettings.DefaultFilePath,
			DefaultDatabaseDSN:     serverSettings.DefaultDatabaseDSN,
			ShutdownTimeout:        serverSettings.ShutdownTimeout,
		},
	}
}

func (c *Settings) ServerAddress() string {
	return c.Server.ServerRunAddress
}

func (c *Settings) ShortenerServerAddress() string {
	return c.Server.ServerShortenerAddress
}

func (c *Settings) ShutdownTimeout() time.Duration {
	return c.Server.ShutdownTimeout
}
