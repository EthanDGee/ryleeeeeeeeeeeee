import chess
import chess.engine
import os
import shutil

from src.player.player import Player


class StockfishPlayer(Player):
    def __init__(self, name="Stockfish", color=True, skill_level=10, time_limit=0.5):
        """
        color: True for white, False for black
        skill_level: 0–20 (Stockfish difficulty)
        time_limit: seconds per move (float)
        """
        super().__init__(name, color)
        self.pending_move = None
        self.time_limit = time_limit
        self.skill_level = skill_level

        self.engine_path = shutil.which("stockfish")
        if self.engine_path is None:
            raise FileNotFoundError(
                "Stockfish binary not found in PATH. Try installing it with `sudo apt install stockfish`.")

        # Initialize Stockfish
        self.engine = chess.engine.SimpleEngine.popen_uci(self.engine_path)
        self.engine.configure({"Skill Level": skill_level})

    def get_move(self, board: chess.Board):
        if board.turn != self.color:
            return None  # Not this bot's turn

        try:
            result = self.engine.play(board, chess.engine.Limit(time=self.time_limit))
            return result.move
        except Exception as e:
            print(f"Stockfish error: {e}")
            return None

    def close(self):
        try:
            self.engine.quit()
        except Exception:
            pass  # Engine may already be closed

    def __del__(self):
        self.close()
