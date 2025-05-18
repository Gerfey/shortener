package models

import "errors"

// Стандартные ошибки сервиса
var (
	// ErrURLExists возвращается, когда URL уже существует в системе
	ErrURLExists = errors.New("url already exists")
	// ErrURLNotFound возвращается, когда URL не найден в системе
	ErrURLNotFound = errors.New("url not found")
)
