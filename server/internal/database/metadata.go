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

	rows, err := db.Query(`SELECT id, url, filename, games, processed, downloaded FROM metadata`)
	if err != nil {
		return nil, err
	}
	defer utils.Close(rows, "rows")

	var results []models.Metadata
	for rows.Next() {
		var metadata models.Metadata
		if err := rows.Scan(&metadata.Id, &metadata.Url, &metadata.Filename, &metadata.Games, &metadata.Processed, &metadata.Downloaded); err != nil {
			return nil, err
		}
		results = append(results, metadata)
	}
	return results, rows.Err()
}

func MetadataExists(filename string) (bool, error) {
	db, err := sql.Open("turso", config.LOCAL_DATABASE_PATH)
	if err != nil {
		return false, err
	}
	defer utils.Close(db, "database")

	var id int
	err = db.QueryRow(`SELECT id FROM metadata WHERE filename = ?`, filename).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
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

func GetSmallestUnprocessedFile() (models.Metadata, error) {
	db, err := sql.Open("turso", config.LOCAL_DATABASE_PATH)
	if err != nil {
		return models.Metadata{}, err
	}
	defer utils.Close(db, "database")

	row := db.QueryRow(`SELECT id, url, filename, games, processed, downloaded FROM metadata WHERE processed < games ORDER BY games - processed ASC LIMIT 1`)

	var metadata models.Metadata
	if err := row.Scan(&metadata.Id, &metadata.Url, &metadata.Filename, &metadata.Games, &metadata.Processed, &metadata.Downloaded); err != nil {
		return models.Metadata{}, err
	}

	return metadata, nil
}

func MarkDownloaded(id int) (bool, error) {
	db, err := sql.Open("turso", config.LOCAL_DATABASE_PATH)
	if err != nil {
		return false, err
	}
	defer utils.Close(db, "database")

	_, err = db.Exec(`UPDATE metadata SET downloaded = true WHERE id = ?`, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func IncrementProcessed(id int) error {
	db, err := sql.Open("turso", config.LOCAL_DATABASE_PATH)
	if err != nil {
		return err
	}
	defer utils.Close(db, "database")

	_, err = db.Exec(`UPDATE metadata SET processed = processed + 1 WHERE id = ?`, id)
	return err
}

func IncrementProcessedBy(id int, count int) error {
	db, err := sql.Open("turso", config.LOCAL_DATABASE_PATH)
	if err != nil {
		return err
	}
	defer utils.Close(db, "database")

	_, err = db.Exec(`UPDATE metadata SET processed = processed + ? WHERE id = ?`, count, id)
	return err
}

func CountDownloaded() (int, error) {
	db, err := sql.Open("turso", config.LOCAL_DATABASE_PATH)
	if err != nil {
		return 0, err
	}
	defer utils.Close(db, "database")

	var count sql.NullInt64
	if err := db.QueryRow(`SELECT SUM(games) FROM metadata WHERE downloaded = true`).Scan(&count); err != nil {
		return 0, err
	}

	return int(count.Int64), nil
}

func CountProcessed() (int, error) {
	db, err := sql.Open("turso", config.LOCAL_DATABASE_PATH)
	if err != nil {
		return 0, err
	}
	defer utils.Close(db, "database")

	var count sql.NullInt64
	if err := db.QueryRow(`SELECT SUM(processed) FROM metadata WHERE downloaded = true`).Scan(&count); err != nil {
		return 0, err
	}

	return int(count.Int64), nil
}
