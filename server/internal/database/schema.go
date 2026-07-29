package database

import (
	"database/sql"

	"server/rest/internal/config"
	"server/rest/internal/utils"
)

func InitializeDatabase() error {
	db, err := sql.Open("turso", config.LOCAL_DATABASE_PATH)
	if err != nil {
		return err
	}
	defer utils.Close(db, "database")

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
