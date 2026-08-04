package models

import (
	"github.com/corentings/chess/v2"
)

type Snapshot struct {
	WhiteElo        int
	BlackElo        int
	WhiteRatingDiff int
	BlackRatingDiff int
	WhiteMaterial   int
	BlackMaterial   int
	WhitesPly       bool
	PlyCount        int
	Board           BoardGrid
}

// NewSnapshot builds the model input describing the position right before the
// move played at ply.
func NewSnapshot(game *chess.Game, ply int) (Snapshot, error) {
	entry, err := MoveAtPly(game, ply)
	if err != nil {
		return Snapshot{}, err
	}

	position := entry.PrePosition
	board, err := NewBoardGrid(position)
	if err != nil {
		return Snapshot{}, err
	}

	whiteMaterial, blackMaterial := Material(position.Board())

	return Snapshot{
		WhiteElo:        TagInt(game, "WhiteElo"),
		BlackElo:        TagInt(game, "BlackElo"),
		WhiteRatingDiff: TagInt(game, "WhiteRatingDiff"),
		BlackRatingDiff: TagInt(game, "BlackRatingDiff"),
		WhiteMaterial:   whiteMaterial,
		BlackMaterial:   blackMaterial,
		WhitesPly:       position.Turn() == chess.White,
		PlyCount:        ply,
		Board:           board,
	}, nil
}
