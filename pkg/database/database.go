package database

import (
	"database/sql"
	"fmt"

	"github.com/askomar/apod/config"
	_ "github.com/lib/pq"
)

func NewDatabase(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable", cfg.Driver, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
