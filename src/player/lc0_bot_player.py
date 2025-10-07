import shutil
import chess
import chess.engine
from src.player.player import Player

class Lc0BotPlayer(Player):
    def __init__(self, name="Lc0", color=True, time_limit=1.0):
        super().__init__(name, color)
        self.pending_move = None
        self.time_limit = time_limit

        self.engine_path = shutil.which("lc0")
        if not self.engine_path:
            raise FileNotFoundError("lc0 executable not found in PATH")

        # Simply start Lc0; make sure your network is in ~/.local/share/lc0/networks/
        self.engine = chess.engine.SimpleEngine.popen_uci([self.engine_path])

    def get_move(self, board: chess.Board):
        if board.turn != self.color:
            return None
        try:
            result = self.engine.play(board, chess.engine.Limit(time=self.time_limit))
            return result.move
        except Exception as e:
            print("Lc0 error:", e)
            return None

    def close(self):
        try:
            self.engine.quit()
        except Exception:
            pass

    def __del__(self):
        self.close()
