package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/Gerfey/shortener/internal/models"
	"github.com/google/uuid"
)

// FileRepository хранилище URL в файле
type FileRepository struct {
	data map[string]models.URLInfo
	Path string
	sync.Mutex
}

// NewFileRepository создает новое файловое хранилище
func NewFileRepository(path string) *FileRepository {
	return &FileRepository{
		data: make(map[string]models.URLInfo),
		Path: path,
	}
}

// Initialize инициализирует хранилище
func (fs *FileRepository) Initialize() error {
	fs.Mutex.Lock()
	defer fs.Mutex.Unlock()

	file, err := os.OpenFile(fs.Path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file for reading: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("error closing file: %v\n", err)
		}
	}()

	stat, statErr := file.Stat()
	if statErr != nil {
		return fmt.Errorf("failed to get file stats: %v", statErr)
	}

	if stat.Size() > 0 {
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&fs.data); err != nil {
			return fmt.Errorf("failed to decode data from file: %v", err)
		}
	}

	return nil
}

// Save сохраняет URL в хранилище
func (fs *FileRepository) Save(ctx context.Context, key, value string, userID string) (string, error) {
	fs.Mutex.Lock()
	defer fs.Mutex.Unlock()

	urlInfo := models.URLInfo{
		UUID:        uuid.New().String(),
		ShortURL:    key,
		OriginalURL: value,
		UserID:      userID,
	}

	fs.data[key] = urlInfo
	return key, nil
}

// SaveBatch сохраняет пакет URL
func (fs *FileRepository) SaveBatch(ctx context.Context, urls map[string]string, userID string) error {
	fs.Mutex.Lock()
	defer fs.Mutex.Unlock()

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

// Find ищет URL по ключу
func (fs *FileRepository) Find(ctx context.Context, key string) (string, bool, bool) {
	fs.Mutex.Lock()
	defer fs.Mutex.Unlock()

	if urlInfo, ok := fs.data[key]; ok {
		return urlInfo.OriginalURL, true, urlInfo.IsDeleted
	}
	return "", false, false
}

// FindShortURL ищет короткий URL
func (fs *FileRepository) FindShortURL(ctx context.Context, originalURL string) (string, error) {
	fs.Mutex.Lock()
	defer fs.Mutex.Unlock()

	for shortURL, urlInfo := range fs.data {
		if urlInfo.OriginalURL == originalURL {
			return shortURL, nil
		}
	}
	return "", fmt.Errorf("URL not found")
}

// All возвращает все URL
func (fs *FileRepository) All(ctx context.Context) map[string]string {
	fs.Mutex.Lock()
	defer fs.Mutex.Unlock()

	result := make(map[string]string)
	for shortURL, urlInfo := range fs.data {
		result[shortURL] = urlInfo.OriginalURL
	}
	return result
}

// GetUserURLs получает URL пользователя
func (fs *FileRepository) GetUserURLs(ctx context.Context, userID string) ([]models.URLPair, error) {
	fs.Mutex.Lock()
	defer fs.Mutex.Unlock()

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

// DeleteUserURLsBatch удаляет URL пользователя
func (fs *FileRepository) DeleteUserURLsBatch(ctx context.Context, shortURLs []string, userID string) error {
	fs.Mutex.Lock()

	for _, shortURL := range shortURLs {
		if urlInfo, exists := fs.data[shortURL]; exists && urlInfo.UserID == userID {
			urlInfo.IsDeleted = true
			fs.data[shortURL] = urlInfo
		}
	}

	fs.Mutex.Unlock()

	return fs.Close()
}

// Ping проверяет доступность хранилища
func (fs *FileRepository) Ping(ctx context.Context) error {
	fs.Mutex.Lock()
	defer fs.Mutex.Unlock()

	file, err := os.OpenFile(fs.Path, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("failed to ping file storage: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("error closing file: %v\n", err)
		}
	}()

	return nil
}

// Close закрывает хранилище
func (fs *FileRepository) Close() error {
	fs.Mutex.Lock()
	defer fs.Mutex.Unlock()

	file, err := os.OpenFile(fs.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("error closing file: %v\n", err)
		}
	}()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(fs.data); err != nil {
		return fmt.Errorf("failed to encode data to file: %v", err)
	}

	return nil
}
