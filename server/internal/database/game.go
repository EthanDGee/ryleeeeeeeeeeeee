package database

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"server/rest/internal/config"
	"server/rest/internal/models"
	"server/rest/internal/utils"

	"github.com/corentings/chess/v2"
)

const tagTimestampLayout = "2006.01.02 15:04:05"

func ParseTagInt(game *chess.Game, key string) int {
	return models.TagInt(game, key)
}

func ParseTagTimestamp(game *chess.Game) *time.Time {
	date := game.GetTagPair("UTCDate")
	timeOfDay := game.GetTagPair("UTCTime")
	timestamp, err := time.Parse(tagTimestampLayout, date+" "+timeOfDay)
	if err != nil {
		log.Printf("failed to parse game timestamp (UTCDate=%q UTCTime=%q): %v", date, timeOfDay, err)
		return nil
	}
	return &timestamp
}

func InsertGame(game chess.Game, minPlyID int, metadata models.Metadata) error {
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
			blackRatingDiff, eco, termination,  totalPlys
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
		ParseTagInt(&game, "WhiteElo"), ParseTagInt(&game, "BlackElo"),
		ParseTagInt(&game, "WhiteRatingDiff"), ParseTagInt(&game, "BlackRatingDiff"),
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
	db, err := sql.Open("turso", config.LOCAL_DATABASE_PATH)
	if err != nil {
		return 0, err
	}
	defer utils.Close(db, "database")

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

	db, err := sql.Open("turso", config.LOCAL_DATABASE_PATH)
	if err != nil {
		return err
	}
	defer utils.Close(db, "database")

	minPlyID, err := TotalPlysPlayed()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	const columnsPerRow = 15
	valuePlaceholders := make([]string, len(games))
	args := make([]any, 0, len(games)*columnsPerRow)

	for i, game := range games {
		valuePlaceholders[i] = "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

		plyCount := len(game.Moves())
		if minPlyID == 0 {
			plyCount -= 1
		}
		maxPlyID := minPlyID + plyCount

		args = append(
			args,
			metadata.Id, game.String(), string(game.Outcome()),
			ParseTagInt(&game, "WhiteElo"), ParseTagInt(&game, "BlackElo"),
			ParseTagInt(&game, "WhiteRatingDiff"), ParseTagInt(&game, "BlackRatingDiff"),
		)
		minPlyID = maxPlyID + 1
	}

	query := `INSERT INTO game (
		fileId, PGN, result, whiteElo, blackElo, whiteRatingDiff,
		blackRatingDiff, eco, termination,
		totalPlys, minPlyId, maxPlyId
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
