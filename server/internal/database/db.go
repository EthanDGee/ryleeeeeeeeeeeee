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

// DB returns the process-wide database handle, opening it on first use. The
// handle is safe for concurrent use by multiple goroutines and must not be
// closed by callers; use Close for that.
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

		// The database is a local file with a single writer, so a one
		// connection pool keeps writes serialized in Go rather than letting
		// them collide and wait on the driver's busy timeout. If read
		// throughput ever matters, add a second read-only pool instead of
		// raising this.
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

// Close releases the shared handle. Intended for process shutdown only.
func Close() error {
	if handle == nil {
		return nil
	}

	db := handle
	handle = nil
	return db.Close()
}
