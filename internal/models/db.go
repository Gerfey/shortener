package models

import "database/sql"

type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Close() error
}
