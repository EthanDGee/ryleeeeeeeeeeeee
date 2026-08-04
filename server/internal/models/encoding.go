package models

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/corentings/chess/v2"
)

const (
	SquareCount = 64
	// MoveSpaceSize is the size of the from-square/to-square move encoding.
	MoveSpaceSize = SquareCount * SquareCount
)

// promotionPieces is the order used by the promotion head.
var promotionPieces = [4]chess.PieceType{chess.Queen, chess.Rook, chess.Bishop, chess.Knight}

var pieceValues = map[chess.PieceType]int{
	chess.Pawn:   1,
	chess.Knight: 3,
	chess.Bishop: 3,
	chess.Rook:   5,
	chess.Queen:  9,
}

// TagInt reads a PGN tag pair as an int, returning 0 when it is missing or malformed.
func TagInt(game *chess.Game, key string) int {
	value, err := strconv.Atoi(game.GetTagPair(key))
	if err != nil {
		return 0
	}
	return value
}

// MoveAtPly returns the main line move played at ply along with the positions
// before and after it. Ply is zero based, so ply 0 is white's first move.
func MoveAtPly(game *chess.Game, ply int) (*chess.MoveHistory, error) {
	if game == nil {
		return nil, errors.New("models: nil game")
	}

	history := game.MoveHistory()
	if ply < 0 || ply >= len(history) {
		return nil, fmt.Errorf("models: ply %d out of range, game has %d moves", ply, len(history))
	}

	entry := history[ply]
	if entry == nil || entry.Move == nil || entry.PrePosition == nil || entry.PostPosition == nil {
		return nil, fmt.Errorf("models: incomplete move history at ply %d", ply)
	}

	return entry, nil
}

func Material(board *chess.Board) (int, int) {
	if board == nil {
		return 0, 0
	}

	white, black := 0, 0
	for _, piece := range board.SquareMap() {
		value := pieceValues[piece.Type()]
		if piece.Color() == chess.White {
			white += value
		} else {
			black += value
		}
	}

	return white, black
}

func MoveIndex(move *chess.Move) int {
	return int(move.S1())*SquareCount + int(move.S2())
}

// EncodeMove one hot encodes a move over the from-square/to-square move space.
func EncodeMove(move *chess.Move) []bool {
	encoded := make([]bool, MoveSpaceSize)
	if move == nil {
		return encoded
	}

	encoded[MoveIndex(move)] = true
	return encoded
}

func EncodeLegalMoves(position *chess.Position) []bool {
	encoded := make([]bool, MoveSpaceSize)
	if position == nil {
		return encoded
	}

	for _, move := range position.ValidMoves() {
		encoded[MoveIndex(&move)] = true
	}

	return encoded
}

// EncodePromotion one hot encodes the promoted piece. All false when the move
// is not a promotion.
func EncodePromotion(piece chess.PieceType) []bool {
	encoded := make([]bool, len(promotionPieces))
	for i, promotionPiece := range promotionPieces {
		if piece == promotionPiece {
			encoded[i] = true
			break
		}
	}

	return encoded
}
