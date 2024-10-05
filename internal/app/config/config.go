package config

type ServerConfig struct {
	RunServerAddress       string
	ShortenerServerAddress string
}

type Config struct {
	Server ServerConfig
}

func NewConfig(Server ServerConfig) *Config {
	return &Config{
		Server: ServerConfig{
			RunServerAddress:       Server.RunServerAddress,
			ShortenerServerAddress: Server.ShortenerServerAddress,
		},
	}
}

func (c *Config) GetServerAddress() string {
	return c.Server.RunServerAddress
}

func (c *Config) GetShortenerServerAddress() string {
	return c.Server.ShortenerServerAddress
}
