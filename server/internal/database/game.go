package database

import (
	"database/sql"
	"log"
	"strings"

	"server/rest/internal/models"
	"server/rest/internal/utils"

	"github.com/corentings/chess/v2"
)

func ParseTagInt(game *chess.Game, key string) int {
	return models.TagInt(game, key)
}

func InsertGame(game chess.Game, minPlyID int, metadata models.Metadata) error {
	db, err := DB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(
		`INSERT INTO game (
			fileId, PGN,  totalPlys, minPlyId, maxPlyId
		) VALUES (?, ?, ?, ?,  ?)`,
	)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	defer utils.Close(stmt, "statement")

	plyCount := len(game.Moves())
	if minPlyID == 0 {
		plyCount -= 1
	}

	_, err = stmt.Exec(
		metadata.Id, game.String(),
		plyCount, minPlyID, minPlyID+plyCount,
	)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return IncrementProcessed(metadata.Id)
}

func TotalPlysPlayed() (int, error) {
	db, err := DB()
	if err != nil {
		return 0, err
	}

	var count sql.NullInt64
	if err := db.QueryRow(`SELECT SUM(totalPlys) FROM game`).Scan(&count); err != nil {
		return 0, err
	}

	return int(count.Int64), nil
}

func InsertGames(games []chess.Game, metadata models.Metadata) error {
	if len(games) == 0 {
		return nil
	}

	db, err := DB()
	if err != nil {
		return err
	}

	minPlyID, err := TotalPlysPlayed()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	const columnsPerRow = 5
	valuePlaceholders := make([]string, len(games))
	args := make([]any, 0, len(games)*columnsPerRow)

	for i, game := range games {
		valuePlaceholders[i] = "(?, ?, ?, ?, ?)"

		plyCount := len(game.Moves())
		if minPlyID == 0 {
			plyCount -= 1
		}
		maxPlyID := minPlyID + plyCount

		args = append(
			args,
			metadata.Id, game.String(), plyCount, minPlyID, maxPlyID,
		)
		minPlyID = maxPlyID + 1
	}

	query := `INSERT INTO game (
		fileId, PGN, totalPlys, minPlyId, maxPlyId
	) VALUES ` + strings.Join(valuePlaceholders, ", ")

	_, err = tx.Exec(query, args...)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return IncrementProcessedBy(metadata.Id, len(games))
}
