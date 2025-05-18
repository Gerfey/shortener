package service

import (
	"fmt"
	"net/url"

	"github.com/Gerfey/shortener/internal/app/settings"
)

// URLService сервис для работы с URL
type URLService struct {
	settings *settings.Settings
}

// NewURLService создает новый сервис URL
func NewURLService(s *settings.Settings) *URLService {
	return &URLService{settings: s}
}

// ShortenerURL формирует полный сокращенный URL
func (us *URLService) ShortenerURL(shortenerID string) (string, error) {
	urlFormat, err := formatURL(us.settings.ShortenerServerAddress())
	if err != nil {
		return "", err
	}

	baseURL, err := url.Parse(urlFormat)
	if err != nil {
		return "", err
	}

	baseURL.RawQuery = ""

	return fmt.Sprintf("%v/%v", baseURL.String(), shortenerID), nil
}

// IsValidURL проверяет валидность URL
func (us *URLService) IsValidURL(URL string) bool {
	if URL == "" {
		return false
	}

	parsedURL, err := url.Parse(URL)
	if err != nil {
		return false
	}

	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

// formatURL форматирует URL в правильный формат
func formatURL(URL string) (string, error) {
	if URL == "" {
		return "", fmt.Errorf("empty URL")
	}

	urlParsed, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	return urlParsed.String(), nil
}
