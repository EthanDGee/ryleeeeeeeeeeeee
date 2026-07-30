package database

import (
	"database/sql"

	"server/rest/internal/config"
	"server/rest/internal/models"
	"server/rest/internal/utils"
)

func UpsertMetadata(metadata models.Metadata) error {
	db, err := sql.Open("turso", config.LOCAL_DATABASE_PATH)
	if err != nil {
		return err
	}
	defer utils.Close(db, "database")

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
	defer utils.Close(db, "database")

	rows, err := db.Query(`SELECT id, url, filename, games, processed FROM metadata`)
	if err != nil {
		return nil, err
	}
	defer utils.Close(rows, "rows")

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

func GetSmallestUndownloadedFile() (models.Metadata, error) {
	db, err := sql.Open("turso", config.LOCAL_DATABASE_PATH)
	if err != nil {
		return models.Metadata{}, err
	}
	defer utils.Close(db, "database")

	row := db.QueryRow(`SELECT id, url, filename, games, processed, downloaded FROM metadata WHERE downloaded = false ORDER BY games ASC LIMIT 1`)

	var metadata models.Metadata
	if err := row.Scan(&metadata.Id, &metadata.Url, &metadata.Filename, &metadata.Games, &metadata.Processed, &metadata.Downloaded); err != nil {
		return models.Metadata{}, err
	}

	return metadata, nil
}
