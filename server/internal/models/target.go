package models

import (
	"github.com/corentings/chess/v2"
)

type Target struct {
	ChosenMove []bool

	// Data for the auxiliary head to improve chess understanding
	Promotion  []bool // Optional Piece selected for promotion
	EnPassant  bool
	Capture    bool
	LegalMoves []bool // A list of legal moves
	Check      bool
	CheckMate  bool
}

// NewTarget builds the training labels for the move played at ply.
func NewTarget(game *chess.Game, ply int) (Target, error) {
	entry, err := MoveAtPly(game, ply)
	if err != nil {
		return Target{}, err
	}

	move := entry.Move

	return Target{
		ChosenMove: EncodeMove(move),
		Promotion:  EncodePromotion(move.Promo()),
		EnPassant:  move.HasTag(chess.EnPassant),
		Capture:    move.HasTag(chess.Capture),
		LegalMoves: EncodeLegalMoves(entry.PrePosition),
		Check:      move.HasTag(chess.Check),
		CheckMate:  entry.PostPosition.Status() == chess.Checkmate,
	}, nil
}
