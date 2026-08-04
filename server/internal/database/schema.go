package database

import (
	"database/sql"
	"log"
	"os"

	"server/rest/internal/config"
	"server/rest/internal/utils"
)

func InitializeDatabase() error {
	if err := os.MkdirAll(config.CACHE_FOLDER, 0o755); err != nil {
		return err
	}

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
              processed INTEGER NOT NULL DEFAULT 0,
              downloaded BOOLEAN NOT NULL DEFAULT FALSE
      )`,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS game (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		fileId INTEGER NOT NULL,
		PGN TEXT NOT NULL,

		totalPlys INTEGER,
		minPlyId INTEGER,
		maxPlyId INTEGER,

		FOREIGN KEY(fileId) REFERENCES metadata(id)
		)`)
	if err != nil {
		return err
	}

	log.Println("database initialized")
	return nil
}
