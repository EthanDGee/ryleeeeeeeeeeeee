package database

import (
	"database/sql"
	"server/rest/internal/config"
	"server/rest/internal/models"
)

func UpsertMetadata(metadata models.Metadata) error {
	db, err := sql.Open("turso", config.LOCAL_DATABASE_PATH)
	if err != nil {
		return err
	}
	defer db.Close()

	var existingID int
	err = db.QueryRow(`SELECT id FROM metadata WHERE filename = ?`, metadata.Filename).Scan(&existingID)

	switch {
	case err == sql.ErrNoRows:
		_, err = db.Exec(
			`INSERT INTO metadata (url, filename, games, processed) VALUES (?, ?, ?, ?)`,
			metadata.Url, metadata.Filename, metadata.Games, metadata.Processed,
		)
		return err
	case err != nil:
		return err
	default:
		_, err = db.Exec(
			`UPDATE metadata SET url = ?, games = ?, processed = ? WHERE filename = ?`,
			metadata.Url, metadata.Games, metadata.Processed, metadata.Filename,
		)
		return err
	}
}

func GetAllMetadata() ([]models.Metadata, error) {
	db, err := sql.Open("turso", config.LOCAL_DATABASE_PATH)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, url, filename, games, processed FROM metadata`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Metadata
	for rows.Next() {
		var metadata models.Metadata
		if err := rows.Scan(&metadata.Id, &metadata.Url, &metadata.Filename, &metadata.Games, &metadata.Processed); err != nil {
			return nil, err
		}
		results = append(results, metadata)
	}
	return results, rows.Err()
}
