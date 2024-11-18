package repository

import (
	"encoding/json"
	"github.com/google/uuid"
	"os"
	"sync"
)

type URLInfo struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type FileRepository struct {
	Path string
	sync.Mutex
}

func NewFileRepository(path string) *FileRepository {
	return &FileRepository{
		Path: path,
	}
}

func (fs *FileRepository) Save(key, value string) error {
	fs.Lock()
	defer fs.Unlock()

	urlInfo := URLInfo{
		UUID:        uuid.New().String(),
		ShortURL:    key,
		OriginalURL: value,
	}

	file, err := os.OpenFile(fs.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(urlInfo)
	if err != nil {
		return err
	}

	_, err = file.Write(append(data, '\n'))
	if err != nil {
		return err
	}

	return nil
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
