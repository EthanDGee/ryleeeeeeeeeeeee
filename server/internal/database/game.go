package database

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"server/rest/internal/config"
	"server/rest/internal/models"
	"server/rest/internal/utils"

	"github.com/corentings/chess/v2"
)

const tagTimestampLayout = "2006.01.02 15:04:05"

func parseTagInt(game *chess.Game, key string) int {
	value, err := strconv.Atoi(game.GetTagPair(key))
	if err != nil {
		return 0
	}
	return value
}

func parseTagTimestamp(game *chess.Game) *time.Time {
	date := game.GetTagPair("UTCDate")
	timeOfDay := game.GetTagPair("UTCTime")
	timestamp, err := time.Parse(tagTimestampLayout, date+" "+timeOfDay)
	if err != nil {
		log.Printf("failed to parse game timestamp (UTCDate=%q UTCTime=%q): %v", date, timeOfDay, err)
		return nil
	}
	return &timestamp
}

func InsertGame(game chess.Game, metadata models.Metadata) error {
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
			fileId, PGN, result, whiteElo, blackElo, whiteRatingDiff,
			blackRatingDiff, timeControl, eco, termination,
			timestamp, variant, totalMoves
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
	)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	defer utils.Close(stmt, "statement")

	_, err = stmt.Exec(
		metadata.Id, game.String(), string(game.Outcome()),
		parseTagInt(&game, "WhiteElo"), parseTagInt(&game, "BlackElo"),
		parseTagInt(&game, "WhiteRatingDiff"), parseTagInt(&game, "BlackRatingDiff"),
		game.GetTagPair("TimeControl"), game.GetTagPair("ECO"), game.GetTagPair("Termination"),
		parseTagTimestamp(&game), game.GetTagPair("Variant"), len(game.Moves()),
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
