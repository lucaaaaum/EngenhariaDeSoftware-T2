package db

import (
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	DSN string
}

func NewPostgresDatabase(config Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", config.DSN)
	if err != nil {
		return nil, errors.Join(errors.New("Failed to connect to PostgreSQL database"), err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	return db, nil
}
