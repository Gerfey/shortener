package settings

import "time"

// ServerSettings настройки сервера
type ServerSettings struct {
	ServerRunAddress       string
	ServerShortenerAddress string
	DefaultFilePath        string
	DefaultDatabaseDSN     string
	ShutdownTimeout        time.Duration
	EnableHTTPS            bool
}

// Settings объединяет все настройки приложения
type Settings struct {
	Server ServerSettings
}

// NewSettings создает новые настройки
func NewSettings(serverSettings ServerSettings) *Settings {
	return &Settings{
		Server: ServerSettings{
			ServerRunAddress:       serverSettings.ServerRunAddress,
			ServerShortenerAddress: serverSettings.ServerShortenerAddress,
			DefaultFilePath:        serverSettings.DefaultFilePath,
			DefaultDatabaseDSN:     serverSettings.DefaultDatabaseDSN,
			ShutdownTimeout:        serverSettings.ShutdownTimeout,
			EnableHTTPS:            serverSettings.EnableHTTPS,
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
