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

	return fmt.Sprintf("%v/%v", urlFormat, shortenerID), nil
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
	urlParsed, err := url.Parse(URL)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", urlParsed.String()), err
}
