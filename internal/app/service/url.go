package service

import (
	"fmt"
	"net/url"

	"github.com/Gerfey/shortener/internal/app/settings"
)

type URLService struct {
	settings *settings.Settings
}

func NewURLService(s *settings.Settings) *URLService {
	return &URLService{settings: s}
}

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
