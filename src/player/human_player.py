from src.player.player import Player
import chess

class HumanPlayer(Player):
    def __init__(self, name: str, color: bool):
        super().__init__(name, color)
        self.pending_move = None  # GUI sets this when user makes a move

    def get_move(self, board: chess.Board) -> chess.Move:
        """Wait for the GUI to set pending_move."""
        if self.pending_move:
            move = self.pending_move
            self.pending_move = None
            return move
        return None
