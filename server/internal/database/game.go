package database

import (
	"database/sql"

	"server/rest/internal/config"
	"server/rest/internal/models"
	"server/rest/internal/utils"
)

func InsertGames(games []models.Game) error {
	db, err := sql.Open("turso", config.LOCAL_DATABASE_PATH)
	if err != nil {
		return err
	}
	defer utils.Close(db, "database")

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(
		`INSERT INTO game (
			fileId, PGN, processed, result,
			whiteElo, blackElo, whiteRatingDiff, blackRatingDiff,
			timeControl, eco, termination, timestamp, variant, totalMoves
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
	)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	defer utils.Close(stmt, "statement")

	for _, game := range games {
		_, err = stmt.Exec(
			game.FileID, game.PGN, game.Processed, game.Result,
			game.WhiteElo, game.BlackElo, game.WhiteRatingDiff, game.BlackRatingDiff,
			game.TimeControl, game.ECO, game.Termination, game.Timestamp, game.Variant, game.TotalMoves,
		)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	}

	return tx.Commit()
}
