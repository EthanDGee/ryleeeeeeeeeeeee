package database

import (
	"database/sql"

	"server/rest/internal/config"
)

func initializeDatabase() {
	db, _ := sql.Open("turso", config.LOCAL_DATABASE_PATH)

	db.Exec(
		`
      CREATE TABLE IF NOT EXISTS metadata (
              id        INTEGER PRIMARY KEY AUTOINCREMENT,
              url       TEXT NOT NULL,
              filename  TEXT NOT NULL,
              games     INTEGER NOT NULL,
              processed BOOLEAN NOT NULL DEFAULT FALSE
      )`,
	)
}
