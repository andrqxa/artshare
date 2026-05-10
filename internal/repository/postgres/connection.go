package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
)

type Connection struct {
	DB *sql.DB
}

func NewConnection(ctx context.Context, databaseURL string) (*Connection, error) {
	if databaseURL == "" {
		return nil, errors.New("database url is required")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, err
	}

	return &Connection{DB: db}, nil
}

func (c *Connection) Close() error {
	if c == nil || c.DB == nil {
		return nil
	}
	return c.DB.Close()
}
