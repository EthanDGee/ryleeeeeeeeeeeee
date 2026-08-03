package models

type Snapshot struct {
	WhiteElo        int
	BlackElo        int
	WhiteRatingDiff int
	BlackRatingDiff int
	WhiteMaterial   int
	BlackMaterial   int
	WhitesMove      bool
	MoveCount       int
	TimeControl     []bool
	Board           BoardGrid
}
