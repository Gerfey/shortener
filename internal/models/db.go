package models

import "database/sql"

// DB интерфейс для работы с базой данных
type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Close() error
}
