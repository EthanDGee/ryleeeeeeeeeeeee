package models

import (
	"errors"

	"github.com/corentings/chess/v2"
)

type BoardGrid struct {
	WhitePawn   [][]bool
	BlackPawn   [][]bool
	WhiteKnight [][]bool
	BlackKnight [][]bool
	WhiteRook   [][]bool
	BlackRook   [][]bool
	WhiteBishop [][]bool
	BlackBishop [][]bool
	WhiteQueen  [][]bool
	BlackQueen  [][]bool
	WhiteKing   [][]bool
	BlackKing   [][]bool
}

func emptyPieceGrid() [][]bool {
	boardSize := 8
	grid := make([][]bool, boardSize)
	for i := range grid {
		grid[i] = make([]bool, boardSize)
	}
	return grid
}

func NewBoardGrid(position *chess.Position) (BoardGrid, error) {
	if position == nil {
		return BoardGrid{}, errors.New("models: game has no position")
	}

	grid := BoardGrid{
		WhitePawn:   emptyPieceGrid(),
		BlackPawn:   emptyPieceGrid(),
		WhiteKnight: emptyPieceGrid(),
		BlackKnight: emptyPieceGrid(),
		WhiteRook:   emptyPieceGrid(),
		BlackRook:   emptyPieceGrid(),
		WhiteBishop: emptyPieceGrid(),
		BlackBishop: emptyPieceGrid(),
		WhiteQueen:  emptyPieceGrid(),
		BlackQueen:  emptyPieceGrid(),
		WhiteKing:   emptyPieceGrid(),
		BlackKing:   emptyPieceGrid(),
	}

	planes := map[chess.Piece][][]bool{
		chess.WhitePawn:   grid.WhitePawn,
		chess.BlackPawn:   grid.BlackPawn,
		chess.WhiteKnight: grid.WhiteKnight,
		chess.BlackKnight: grid.BlackKnight,
		chess.WhiteRook:   grid.WhiteRook,
		chess.BlackRook:   grid.BlackRook,
		chess.WhiteBishop: grid.WhiteBishop,
		chess.BlackBishop: grid.BlackBishop,
		chess.WhiteQueen:  grid.WhiteQueen,
		chess.BlackQueen:  grid.BlackQueen,
		chess.WhiteKing:   grid.WhiteKing,
		chess.BlackKing:   grid.BlackKing,
	}

	for square, piece := range position.Board().SquareMap() {
		plane, ok := planes[piece]
		if !ok {
			continue
		}

		plane[int(square)/8][int(square)%8] = true
	}

	return grid, nil
}
