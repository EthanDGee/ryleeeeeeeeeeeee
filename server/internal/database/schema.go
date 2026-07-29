package database

import (
	"database/sql"

	"server/rest/internal/config"
)

func initializeDatabase() error {
	db, err := sql.Open("turso", config.LOCAL_DATABASE_PATH)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		`
      CREATE TABLE IF NOT EXISTS metadata (
              id        INTEGER PRIMARY KEY AUTOINCREMENT,
              url       TEXT NOT NULL,
              filename  TEXT NOT NULL,
              games     INTEGER NOT NULL,
              processed BOOLEAN NOT NULL DEFAULT FALSE
      )`,
	)
	return err
}
