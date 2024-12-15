package repository

import (
	"encoding/json"
	"fmt"
	"github.com/Gerfey/shortener/internal/models"
	"github.com/google/uuid"
	"os"
	"sync"
)

type FileRepository struct {
	data map[string]models.URLInfo
	Path string
	sync.Mutex
}

func NewFileRepository(path string) *FileRepository {
	return &FileRepository{
		data: make(map[string]models.URLInfo),
		Path: path,
	}
}

func (fs *FileRepository) Initialize() error {
	fs.Lock()
	defer fs.Unlock()

	file, err := os.OpenFile(fs.Path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file for reading: %v", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stats: %v", err)
	}

	if stat.Size() > 0 {
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&fs.data); err != nil {
			return fmt.Errorf("failed to decode data from file: %v", err)
		}
	}

	return nil
}

func (fs *FileRepository) Save(key, value string, userID string) (string, error) {
	fs.Lock()
	defer fs.Unlock()

	urlInfo := models.URLInfo{
		UUID:        uuid.New().String(),
		ShortURL:    key,
		OriginalURL: value,
		UserID:      userID,
	}

	fs.data[key] = urlInfo
	return key, nil
}

func (fs *FileRepository) SaveBatch(urls map[string]string, userID string) error {
	fs.Lock()
	defer fs.Unlock()

	for shortURL, originalURL := range urls {
		urlInfo := models.URLInfo{
			UUID:        uuid.New().String(),
			ShortURL:    shortURL,
			OriginalURL: originalURL,
			UserID:      userID,
		}
		fs.data[shortURL] = urlInfo
	}

	return nil
}

func (fs *FileRepository) Find(key string) (string, bool) {
	fs.Lock()
	defer fs.Unlock()

	if urlInfo, ok := fs.data[key]; ok {
		return urlInfo.OriginalURL, true
	}
	return "", false
}

func (fs *FileRepository) FindShortURL(originalURL string) (string, error) {
	fs.Lock()
	defer fs.Unlock()

	for shortURL, urlInfo := range fs.data {
		if urlInfo.OriginalURL == originalURL {
			return shortURL, nil
		}
	}
	return "", fmt.Errorf("URL not found")
}

func (fs *FileRepository) All() map[string]string {
	fs.Lock()
	defer fs.Unlock()

	result := make(map[string]string)
	for shortURL, urlInfo := range fs.data {
		result[shortURL] = urlInfo.OriginalURL
	}
	return result
}

func (fs *FileRepository) GetUserURLs(userID string) ([]models.URLPair, error) {
	fs.Lock()
	defer fs.Unlock()

	var userURLs []models.URLPair
	for _, urlInfo := range fs.data {
		if urlInfo.UserID == userID {
			userURLs = append(userURLs, models.URLPair{
				ShortURL:    urlInfo.ShortURL,
				OriginalURL: urlInfo.OriginalURL,
			})
		}
	}
	return userURLs, nil
}

func (fs *FileRepository) Ping() error {
	fs.Lock()
	defer fs.Unlock()

	file, err := os.OpenFile(fs.Path, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("failed to ping file storage: %v", err)
	}
	defer file.Close()

	return nil
}

func (fs *FileRepository) Close() error {
	fs.Lock()
	defer fs.Unlock()

	file, err := os.OpenFile(fs.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(fs.data); err != nil {
		return fmt.Errorf("failed to encode data to file: %v", err)
	}

	return nil
}
