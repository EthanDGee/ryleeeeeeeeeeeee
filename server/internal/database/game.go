package database

import (
	"database/sql"
	"strings"

	"server/rest/internal/models"

	"github.com/corentings/chess/v2"
)

func ParseTagInt(game *chess.Game, key string) int {
	return models.TagInt(game, key)
}

// NextPlyID scans the whole game table, so call it once when resuming rather
// than once per batch.
func NextPlyID() (int, error) {
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

// InsertGames gives each game a contiguous block of ply IDs starting at
// nextPlyID and returns the first unused ply ID.
func InsertGames(games []chess.Game, metadata models.Metadata, nextPlyID int) (int, error) {
	if len(games) == 0 {
		return nextPlyID, nil
	}

	db, err := DB()
	if err != nil {
		return nextPlyID, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nextPlyID, err
	}
	defer func() {
		_ = tx.Rollback() // no-op once Commit succeeded
	}()

	const columnsPerRow = 5
	valuePlaceholders := make([]string, len(games))
	args := make([]any, 0, len(games)*columnsPerRow)

	plyID := nextPlyID
	for i, game := range games {
		valuePlaceholders[i] = "(?, ?, ?, ?, ?)"

		// n plies owns IDs [plyID, plyID+n-1], so ply p is ID plyID+p, matching
		// the 0-based indexing models.MoveAtPly expects. A game with no plies
		// gets an empty range that no lookup can match.
		plyCount := len(game.MoveHistory())

		args = append(
			args,
			metadata.Id, game.String(), plyCount, plyID, plyID+plyCount-1,
		)
		plyID += plyCount
	}

	query := `INSERT INTO game (
		fileId, PGN, totalPlys, minPlyId, maxPlyId
	) VALUES ` + strings.Join(valuePlaceholders, ", ")

	if _, err := tx.Exec(query, args...); err != nil {
		return nextPlyID, err
	}

	// Same transaction as the rows, so a crash cannot leave processed
	// disagreeing with what was stored.
	if _, err := tx.Exec(
		`UPDATE metadata SET processed = processed + ? WHERE id = ?`,
		len(games), metadata.Id,
	); err != nil {
		return nextPlyID, err
	}

	if err := tx.Commit(); err != nil {
		return nextPlyID, err
	}

	return plyID, nil
}
