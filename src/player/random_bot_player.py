import random
import chess
from src.player.player import Player

class RandomBotPlayer(Player):
    def get_move(self, board: chess.Board) -> chess.Move:
        moves = list(board.legal_moves)
        return random.choice(moves) if moves else None
