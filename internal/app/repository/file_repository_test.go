package repository

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var testFileStorage *FileRepository

func init() {
	testFileStorage = NewFileRepository("test.json")
}

func TestFileStorageSave(t *testing.T) {
	var testCases = []struct {
		desc     string
		urlInfo  URLInfo
		expected error
	}{
		{
			desc: "Valid URL info",
			urlInfo: URLInfo{
				UUID:        "123",
				ShortURL:    "short",
				OriginalURL: "https://original.com",
			},
			expected: nil,
		},
		{
			desc:     "Empty URL info",
			urlInfo:  URLInfo{},
			expected: nil,
		},
		{
			desc: "Invalid Original URL",
			urlInfo: URLInfo{
				UUID:        "123",
				ShortURL:    "short",
				OriginalURL: "https//original.com",
			},
			expected: nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := testFileStorage.Save(tC.urlInfo.ShortURL, tC.urlInfo.OriginalURL)
			if !errors.Is(err, tC.expected) {
				assert.Error(t, err)
			}
		})

		cleanFileStorage()
	}
}

func TestFileStorageLoad(t *testing.T) {
	var testCases = []struct {
		desc        string
		urlInfo     URLInfo
		expectError bool
	}{
		{
			desc: "Valid URL Info",
			urlInfo: URLInfo{
				UUID:        "123",
				ShortURL:    "short",
				OriginalURL: "https://original.com",
			},
			expectError: false,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := testFileStorage.Save(tC.urlInfo.ShortURL, tC.urlInfo.OriginalURL)
			if err != nil {
				t.Fatal(err)
			}

			_ = testFileStorage.All()

			cleanFileStorage()
		})
	}
}

func cleanFileStorage() {
	err := os.Remove(testFileStorage.Path)
	if err != nil {
		panic(err)
	}
}
