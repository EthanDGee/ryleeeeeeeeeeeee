package database

import (
	"log"
)

func InitializeDatabase() error {
	db, err := DB()
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`
      CREATE TABLE IF NOT EXISTS metadata (
              id         INTEGER PRIMARY KEY AUTOINCREMENT,
              url        TEXT NOT NULL,
              filename   TEXT NOT NULL,
              games      INTEGER NOT NULL,
              processed  INTEGER NOT NULL DEFAULT 0,
				      corrupted  INTEGER NOT NULL DEFAULT 0,
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

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_game_min_ply ON game(minPlyId)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_game_max_ply ON game(maxPlyId)`)
	if err != nil {
		return err
	}

	log.Println("database initialized")
	return nil
}
