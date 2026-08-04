package database

import (
	"database/sql"
	"strings"

	"server/rest/internal/config"
	"server/rest/internal/models"
	"server/rest/internal/utils"

	"github.com/corentings/chess/v2"
)

func GetPly(id int) (models.Snapshot, models.Target, error) {
	db, err := sql.Open("turso", config.LOCAL_DATABASE_PATH)
	if err != nil {
		return models.Snapshot{}, models.Target{}, err
	}
	defer utils.Close(db, "database")

	var (
		pgn      string
		minPlyID int
	)
	err = db.QueryRow(
		`SELECT PGN, minPlyId FROM game WHERE minPlyId <= ? AND ? <= maxPlyId LIMIT 1`, id, id,
	).Scan(&pgn, &minPlyID)
	if err != nil {
		return models.Snapshot{}, models.Target{}, err
	}

	fromPGN, err := chess.PGN(strings.NewReader(pgn))
	if err != nil {
		return models.Snapshot{}, models.Target{}, err
	}
	game := chess.NewGame(fromPGN)

	ply := id - minPlyID

	snapshot, err := models.NewSnapshot(game, ply)
	if err != nil {
		return models.Snapshot{}, models.Target{}, err
	}

	target, err := models.NewTarget(game, ply)
	if err != nil {
		return models.Snapshot{}, models.Target{}, err
	}

	return snapshot, target, nil
}
