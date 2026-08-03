package models

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
