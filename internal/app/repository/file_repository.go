package repository

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"
	"sync"
)

type URLInfo struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type FileRepository struct {
	data map[string]string
	Path string
	sync.Mutex
}

func NewFileRepository(path string) *FileRepository {
	return &FileRepository{
		data: make(map[string]string),
		Path: path,
	}
}

func (fs *FileRepository) SaveBatch(urls map[string]string) error {
	fs.Lock()
	defer fs.Unlock()

	for shortURL, originalURL := range urls {
		fs.data[shortURL] = originalURL
	}

	return fs.saveToFile()
}

func (fs *FileRepository) Load() error {
	fs.Lock()
	defer fs.Unlock()

	file, err := os.Open(fs.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&fs.data)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to decode file: %w", err)
	}

	return nil
}

func (fs *FileRepository) Save(key, value string) (string, error) {
	fs.Lock()
	defer fs.Unlock()

	urlInfo := URLInfo{
		UUID:        uuid.New().String(),
		ShortURL:    key,
		OriginalURL: value,
	}

	file, err := os.OpenFile(fs.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return key, err
	}
	defer file.Close()

	data, err := json.Marshal(urlInfo)
	if err != nil {
		return key, err
	}

	_, err = file.Write(append(data, '\n'))
	if err != nil {
		return key, err
	}

	return key, nil
}

func (fs *FileRepository) All() map[string]string {
	fs.Lock()
	defer fs.Unlock()

	file, err := os.Open(fs.Path)
	if err != nil {
		return nil
	}
	defer file.Close()

	urlStore := make(map[string]string)

	decoder := json.NewDecoder(file)
	for decoder.More() {
		var urlInfo URLInfo
		if err := decoder.Decode(&urlInfo); err != nil {
			continue
		}
		urlStore[urlInfo.ShortURL] = urlInfo.OriginalURL
	}

	return urlStore
}

func (fs *FileRepository) Find(key string) (string, bool) {
	fs.Lock()
	defer fs.Unlock()

	file, err := os.Open(fs.Path)
	if err != nil {
		return "", false
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for decoder.More() {
		var urlInfo URLInfo
		if err := decoder.Decode(&urlInfo); err != nil {
			continue
		}

		if urlInfo.ShortURL == key {
			return urlInfo.OriginalURL, true
		}
	}

	return "", false
}

func (fs *FileRepository) FindShortURL(originalURL string) (string, error) {
	fs.Lock()
	defer fs.Unlock()

	for shortURL, storedOriginalURL := range fs.data {
		if storedOriginalURL == originalURL {
			return shortURL, nil
		}
	}

	return "", fmt.Errorf("original URL not found")
}

func (fs *FileRepository) saveToFile() error {
	fs.Lock()
	defer fs.Unlock()

	file, err := os.OpenFile(fs.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(fs.data)
}
