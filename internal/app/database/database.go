package database

import (
	"fmt"
	"github.com/jackc/pgx"
)

type Database struct {
	config pgx.ConnConfig
}

func NewDatabase(s string) (*Database, error) {
	config, err := pgx.ParseDSN(s)

	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	return &Database{config: config}, nil
}

func (db *Database) Connect() (*pgx.Conn, error) {
	conn, err := pgx.Connect(db.config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return conn, nil
}
