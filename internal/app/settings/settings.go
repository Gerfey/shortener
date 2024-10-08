package settings

type ServerSettings struct {
	ServerRunAddress       string
	ServerShortenerAddress string
}

type Settings struct {
	Server ServerSettings
}

func NewSettings(serverSettings ServerSettings) *Settings {
	return &Settings{
		Server: ServerSettings{
			ServerRunAddress:       serverSettings.ServerRunAddress,
			ServerShortenerAddress: serverSettings.ServerShortenerAddress,
		},
	}
}

func (c *Settings) ServerAddress() string {
	return c.Server.ServerRunAddress
}

func (c *Settings) ShortenerServerAddress() string {
	return c.Server.ServerShortenerAddress
}
