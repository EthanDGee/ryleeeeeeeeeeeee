package database

import (
	"database/sql"
	"os"
	"sync"

	"server/rest/internal/config"

	_ "turso.tech/database/tursogo"
)

var (
	handle   *sql.DB
	openOnce sync.Once
	openErr  error
)

func DB() (*sql.DB, error) {
	openOnce.Do(func() {
		if err := os.MkdirAll(config.CACHE_FOLDER, 0o755); err != nil {
			openErr = err
			return
		}

		db, err := sql.Open("turso", config.LOCAL_DATABASE_PATH)
		if err != nil {
			openErr = err
			return
		}

		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
		db.SetConnMaxLifetime(0)

		if err := db.Ping(); err != nil {
			_ = db.Close()
			openErr = err
			return
		}

		handle = db
	})

	return handle, openErr
}

func Close() error {
	if handle == nil {
		return nil
	}

	db := handle
	handle = nil
	return db.Close()
}
