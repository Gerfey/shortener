package settings

import "time"

type ServerSettings struct {
	ServerRunAddress       string
	ServerShortenerAddress string
	DefaultFilePath        string
	DefaultDatabaseDSN     string
	ShutdownTimeout        time.Duration
}

// Settings объединяет все настройки приложения
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

// ServerAddress возвращает адрес для запуска HTTP-сервера
func (c *Settings) ServerAddress() string {
	return c.Server.ServerRunAddress
}

// ShortenerServerAddress возвращает базовый URL для сокращенных ссылок
func (c *Settings) ShortenerServerAddress() string {
	return c.Server.ServerShortenerAddress
}

// ShutdownTimeout возвращает таймаут для корректного завершения работы сервера
func (c *Settings) ShutdownTimeout() time.Duration {
	return c.Server.ShutdownTimeout
}
