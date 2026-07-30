package models

import "time"

type Game struct {
	ID        int
	FileID    int
	PGN       string
	Processed bool

	Result string

	// Player ratings
	WhiteElo        int
	BlackElo        int
	WhiteRatingDiff int
	BlackRatingDiff int

	// Time control
	TimeControl string

	// Opening information
	ECO string

	// Game termination
	Termination string

	// Timestamp information
	Timestamp time.Time

	// Variant
	Variant string

	// Calculated fields
	TotalMoves int
}
