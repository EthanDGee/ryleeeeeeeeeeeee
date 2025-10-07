import chess
from abc import ABC, abstractmethod

class Player(ABC):
    def __init__(self, name: str, color: bool):
        self.name = name
        self.color = color  # True = White, False = Black
        self.wins = 0       # Number of wins
        self.losses = 0     # Number of losses

    def record_win(self):
        self.wins += 1

    def record_loss(self):
        self.losses += 1

    @abstractmethod
    def get_move(self, board: chess.Board) -> chess.Move:
        """Return a move to play on the given board."""
        pass
