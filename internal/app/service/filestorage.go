package service

import (
	"encoding/json"
	"os"
	"sync"
)

type URLInfo struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type FileStorage struct {
	Path string
	sync.Mutex
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{
		Path: path,
	}
}

func (fs *FileStorage) Save(urlInfo URLInfo) error {
	fs.Lock()
	defer fs.Unlock()

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

func (fs *FileStorage) Load() (map[string]URLInfo, error) {
	fs.Lock()
	defer fs.Unlock()

	file, err := os.Open(fs.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	urlStore := make(map[string]URLInfo)

	decoder := json.NewDecoder(file)
	for decoder.More() {
		var urlInfo URLInfo
		if err := decoder.Decode(&urlInfo); err != nil {
			continue
		}
		urlStore[urlInfo.ShortURL] = urlInfo
	}

	return urlStore, nil
}
