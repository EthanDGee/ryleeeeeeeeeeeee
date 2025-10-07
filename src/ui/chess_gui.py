import datetime
import os
import random
import tkinter as tk
import urllib.request
from PIL import Image, ImageTk
import chess
import chess.pgn
import sys

from src.player.human_player import HumanPlayer

# ================= Configuration =================

# --- Default window size ---
DEFAULT_WINDOW_WIDTH = 1300
DEFAULT_WINDOW_HEIGHT = 1000

ENABLE_HIGHLIGHT_LEGAL_MOVES = True
ENABLE_HIGHLIGHT_LAST_MOVE = True
ENABLE_HIGHLIGHT_SELECTED = True
ENABLE_HIGHLIGHT_ILLEGAL_MOVE = True
ENABLE_HIGHLIGHT_CAPTURE_SQUARE = True

COLOR_BOARD_LIGHT = "#F0D9B5"
COLOR_BOARD_DARK = "#B58863"
COLOR_HIGHLIGHT_LEGAL = "#A9D18E"
COLOR_HIGHLIGHT_ILLEGAL = "#FF6666"
COLOR_HIGHLIGHT_LAST_MOVE = "#FBFFB8"
COLOR_HIGHLIGHT_SELECTED = "#FFFF00"
COLOR_HIGHLIGHT_CAPTURE_SQUARE = "#FFD700"

PIECE_CODES = {
    'P': 'wp', 'R': 'wr', 'N': 'wn', 'B': 'wb', 'Q': 'wq', 'K': 'wk',
    'p': 'bp', 'r': 'br', 'n': 'bn', 'b': 'bb', 'q': 'bq', 'k': 'bk',
}

PIECE_VALUES = {
    chess.PAWN: 1,
    chess.KNIGHT: 3,
    chess.BISHOP: 3,
    chess.ROOK: 5,
    chess.QUEEN: 9,
    chess.KING: 0,  # King is invaluable, ignore for score
}

IMAGE_DIR = os.path.join(os.path.dirname(__file__), "images")
BASE_URL = "https://images.chesscomfiles.com/chess-themes/pieces/neo/150/"

# Default PGN save directory
PGN_SAVE_DIR = os.path.join(os.path.dirname(__file__), "games")

# ================= GUI CLASS =================

class ChessGui:
    def __init__(self, root, white_player, black_player, save_dir=PGN_SAVE_DIR):
        self.root = root
        self.root.title("Dynamic Chess GUI")

        # --- Game state ---
        self.board = chess.Board()
        self.white_player = white_player
        self.black_player = black_player
        self.current_player = self.white_player if self.board.turn else self.black_player
        self.save_dir = save_dir

        self.selected_square = None
        self.legal_moves = []
        self.last_move = None
        self.illegal_dest = None
        self.capture_square = None

        self.tile_size = 64
        self.piece_images_raw = {}
        self.piece_images_scaled = {}

        # ================= Layout =================
        main_frame = tk.Frame(root)
        main_frame.pack(fill="both", expand=True)

        # Left = Chess Board
        self.canvas = tk.Canvas(main_frame, bg='black')
        self.canvas.pack(side="left", fill="both", expand=True)

        # Right = Sidebar (stats + moves)
        self.sidebar = tk.Frame(main_frame, width=220, bg="#EEE")
        self.sidebar.pack(side="right", fill="y")

        # Turn indicator
        self.turn_label = tk.Label(self.sidebar, text="", font=("Arial", 14, "bold"), bg="#EEE")
        self.turn_label.pack(pady=(10, 10))

        # ---------- Player Stats with current game score ----------
        # Black player frame
        black_frame = tk.Frame(self.sidebar, bg="#EEE")
        black_frame.pack(pady=(0, 10))
        tk.Label(black_frame, text="Black:", font=("Arial", 12, "bold"), bg="#EEE").pack(side="left")
        self.black_name_label = tk.Label(black_frame, text=self.black_player.name, font=("Arial", 12, "bold"),
                                         bg="#EEE")
        self.black_name_label.pack(side="left", padx=(5, 10))
        self.black_score_label = tk.Label(black_frame, text="Score: 0", font=("Arial", 10), bg="#EEE")
        self.black_score_label.pack(side="left")

        self.black_stats = tk.Label(self.sidebar, text="Wins: 0  Losses: 0", font=("Arial", 10), bg="#EEE")
        self.black_stats.pack(pady=(0, 10))

        # White player frame
        white_frame = tk.Frame(self.sidebar, bg="#EEE")
        white_frame.pack(pady=(0, 10))
        tk.Label(white_frame, text="White:", font=("Arial", 12, "bold"), bg="#EEE").pack(side="left")
        self.white_name_label = tk.Label(white_frame, text=self.white_player.name, font=("Arial", 12, "bold"),
                                         bg="#EEE")
        self.white_name_label.pack(side="left", padx=(5, 10))
        self.white_score_label = tk.Label(white_frame, text="Score: 0", font=("Arial", 10), bg="#EEE")
        self.white_score_label.pack(side="left")

        self.white_stats = tk.Label(self.sidebar, text="Wins: 0  Losses: 0", font=("Arial", 10), bg="#EEE")
        self.white_stats.pack(pady=(0, 20))

        # ---------- Move list ----------
        tk.Label(self.sidebar, text="Moves", font=("Arial", 12, "bold"), bg="#EEE").pack()
        self.move_listbox = tk.Listbox(self.sidebar, font=("Consolas", 10), width=25, height=20)
        self.move_listbox.pack(padx=10, pady=5, fill="both", expand=True)

        # Save button
        self.save_button = tk.Button(self.sidebar, text="Save Game", command=self.save_game)
        self.save_button.pack(pady=10, fill="x")

        # New Game button
        self.new_game_button = tk.Button(self.sidebar, text="New Game", command=self.reset_game, bg="#55AAFF",
                                         fg="white")
        self.new_game_button.pack(pady=(0, 5), fill="x")

        # Quit button
        self.quit_button = tk.Button(self.sidebar, text="Quit", command=self.root.destroy, bg="#FF5555", fg="white")
        self.quit_button.pack(pady=(5, 10), fill="x")

        # Bind events
        self.canvas.bind("<Configure>", self.on_resize)
        self.canvas.bind("<Button-1>", self.on_click)

        # Game over label
        self.game_over_label = tk.Label(self.sidebar, text="", font=("Arial", 12, "bold"), bg="#EEE", fg="red")
        self.game_over_label.pack(pady=5)

        # Play Again button (hidden initially)
        self.play_again_button = tk.Button(self.sidebar, text="Play Again", command=self.reset_game)
        self.play_again_button.pack(pady=10, fill="x")
        self.play_again_button.pack_forget()

        # Load piece images
        os.makedirs(IMAGE_DIR, exist_ok=True)
        self.download_and_load_images()

        # Initialize turn and start game loop
        self.update_turn_label()
        self.root.after(100, self.game_loop)

    # ---------------- PGN SAVE ----------------
    def save_game(self):
        """Save the current game in PGN format."""
        os.makedirs(self.save_dir, exist_ok=True)

        game = chess.pgn.Game()
        game.headers["Event"] = "Friendly Match"
        game.headers["White"] = self.white_player.name
        game.headers["Black"] = self.black_player.name
        game.headers["Date"] = datetime.datetime.now().strftime("%Y.%m.%d")
        game.headers["Result"] = self.board.result() if self.board.is_game_over() else "*"

        node = game
        for move in self.board.move_stack:
            node = node.add_variation(move)

        filename = os.path.join(
            self.save_dir,
            f"{self.white_player.name}_vs_{self.black_player.name}_{datetime.datetime.now().strftime('%Y-%m-%d_%H-%M-%S')}.pgn"
        )

        with open(filename, "w", encoding="utf-8") as f:
            print(game, file=f)

        print(f"Game saved to {filename}")

    def reset_game(self):
        self.board.reset()
        self.last_move = None
        self.selected_square = None
        self.legal_moves = []
        self.illegal_dest = None
        self.capture_square = None
        self.current_player = self.white_player if self.board.turn else self.black_player

        self.move_listbox.delete(0, tk.END)
        self.game_over_label.config(text="")
        self.play_again_button.pack_forget()
        self.update_turn_label()
        self.draw_board()

    def quit_game(self):
        """Stop all bots, cancel loops, and fully close the GUI."""
        for player in [self.white_player, self.black_player]:
            if hasattr(player, "stop"):
                try:
                    player.stop()
                except Exception:
                    pass
            if hasattr(player, "close"):
                try:
                    player.close()
                except Exception:
                    pass

        # Cancel the game loop after schedule
        if hasattr(self, "after_id"):
            self.root.after_cancel(self.after_id)

        # Destroy Tkinter window
        self.root.destroy()

        # Ensure Python process exits
        sys.exit(0)

    def update_scores(self):
        white_score = 0
        black_score = 0
        for square in chess.SQUARES:
            piece = self.board.piece_at(square)
            if piece:
                value = PIECE_VALUES[piece.piece_type]
                if piece.color == chess.WHITE:
                    white_score += value
                else:
                    black_score += value

        # Subtract from 39 to get material difference from starting position if you want relative?
        # But simpler: just show current material per player
        self.white_score_label.config(text=f"Score: {white_score}")
        self.black_score_label.config(text=f"Score: {black_score}")

    # ---------------- Image Loading ----------------
    def download_and_load_images(self):
        for symbol, code in PIECE_CODES.items():
            filename = f"{code}.png"
            path = os.path.join(IMAGE_DIR, filename)
            url = BASE_URL + filename
            if not os.path.exists(path):
                try:
                    print(f"Downloading {filename}...")
                    urllib.request.urlretrieve(url, path)
                except Exception as e:
                    print(f"Download failed for {filename}: {e}")
                    continue
            try:
                img = Image.open(path)
                self.piece_images_raw[symbol] = img.convert("RGBA")
            except Exception as e:
                print(f"Failed to load image {filename}: {e}")

    def scale_images(self, tile_size):
        resample = Image.Resampling.LANCZOS if hasattr(Image, "Resampling") else Image.LANCZOS
        self.piece_images_scaled = {
            sym: ImageTk.PhotoImage(img.resize((tile_size, tile_size), resample))
            for sym, img in self.piece_images_raw.items()
        }

    # ---------------- Drawing ----------------
    def square_to_xy(self, square):
        col = chess.square_file(square)
        row = 7 - chess.square_rank(square)
        return col * self.tile_size, row * self.tile_size

    def draw_board(self):
        self.canvas.delete("all")
        width = self.canvas.winfo_width()
        height = self.canvas.winfo_height()
        if width <= 0 or height <= 0:
            return

        self.tile_size = min(width, height) // 8
        self.scale_images(self.tile_size)

        # Draw board tiles
        for square in chess.SQUARES:
            x, y = self.square_to_xy(square)
            col = chess.square_file(square)
            row = chess.square_rank(square)
            color = COLOR_BOARD_LIGHT if (col + row) % 2 == 0 else COLOR_BOARD_DARK
            self.canvas.create_rectangle(x, y, x + self.tile_size, y + self.tile_size, fill=color, outline="")

        # Highlights
        if ENABLE_HIGHLIGHT_LAST_MOVE and self.last_move:
            for sq in (self.last_move.from_square, self.last_move.to_square):
                x, y = self.square_to_xy(sq)
                self.canvas.create_rectangle(x, y, x + self.tile_size, y + self.tile_size,
                                             fill=COLOR_HIGHLIGHT_LAST_MOVE, outline="")

        if ENABLE_HIGHLIGHT_CAPTURE_SQUARE and self.capture_square is not None:
            x, y = self.square_to_xy(self.capture_square)
            self.canvas.create_rectangle(x, y, x + self.tile_size, y + self.tile_size,
                                         fill=COLOR_HIGHLIGHT_CAPTURE_SQUARE, outline="")

        if ENABLE_HIGHLIGHT_ILLEGAL_MOVE and self.illegal_dest is not None:
            x, y = self.square_to_xy(self.illegal_dest)
            self.canvas.create_rectangle(x, y, x + self.tile_size, y + self.tile_size,
                                         fill=COLOR_HIGHLIGHT_ILLEGAL, outline="")
            self.illegal_dest = None

        if self.selected_square is not None:
            if ENABLE_HIGHLIGHT_SELECTED:
                x0, y0 = self.square_to_xy(self.selected_square)
                self.canvas.create_rectangle(x0, y0, x0 + self.tile_size, y0 + self.tile_size,
                                             fill=COLOR_HIGHLIGHT_SELECTED, outline="")
            if ENABLE_HIGHLIGHT_LEGAL_MOVES:
                for move in self.legal_moves:
                    x, y = self.square_to_xy(move.to_square)
                    self.canvas.create_rectangle(x, y, x + self.tile_size, y + self.tile_size,
                                                 fill=COLOR_HIGHLIGHT_LEGAL, outline="")

        # Draw pieces
        for square in chess.SQUARES:
            piece = self.board.piece_at(square)
            if piece:
                sym = piece.symbol()
                if sym in self.piece_images_scaled:
                    x, y = self.square_to_xy(square)
                    self.canvas.create_image(x, y, anchor="nw", image=self.piece_images_scaled[sym])

        # Game over message
        if self.board.is_game_over():
            result = self.board.result()

            if result == "1-0":
                winner = self.white_player
                loser = self.black_player
            elif result == "0-1":
                winner = self.black_player
                loser = self.white_player
            else:
                winner = None
                loser = None

            if winner:
                winner.record_win()
                loser.record_loss()
                self.game_over_label.config(text=f"Game Over: {winner.name} won!")
            else:
                self.game_over_label.config(text="Game Over: Draw!")

            # Update stats labels
            self.white_stats.config(text=f"Wins: {self.white_player.wins}  Losses: {self.white_player.losses}")
            self.black_stats.config(text=f"Wins: {self.black_player.wins}  Losses: {self.black_player.losses}")

            # Show Play Again button
            self.play_again_button.pack()

            self.save_game()

    # ---------------- Events ----------------
    def on_resize(self, event):
        self.draw_board()

    def on_click(self, event):
        if not isinstance(self.current_player, HumanPlayer):
            return

        col = event.x // self.tile_size
        row = 7 - (event.y // self.tile_size)
        square = chess.square(col, row)

        if self.selected_square is None:
            piece = self.board.piece_at(square)
            if piece and piece.color == self.board.turn:
                self.selected_square = square
                self.legal_moves = [m for m in self.board.legal_moves if m.from_square == square]
        else:
            move = chess.Move(self.selected_square, square)
            if move not in self.board.legal_moves:
                move = chess.Move(self.selected_square, square, promotion=chess.QUEEN)
            if move in self.board.legal_moves:
                self.current_player.pending_move = move
                self.selected_square = None
                self.legal_moves = []
            else:
                self.illegal_dest = square
                self.selected_square = None
                self.legal_moves = []

        self.draw_board()

    # ---------------- Game Loop ----------------
    def update_turn_label(self):
        if self.board.is_game_over():
            self.turn_label.config(text="")
        else:
            color_text = "White" if self.board.turn else "Black"
            player_name = self.white_player.name if self.board.turn else self.black_player.name
            self.turn_label.config(text=f"{color_text}'s turn — {player_name}")

    def update_move_list(self, san_move=None):
        if san_move:
            moves = list(self.board.move_stack)
            move_num = len(moves)
            if move_num % 2 == 1:
                self.move_listbox.insert(tk.END, f"{(move_num + 1) // 2}. {san_move}")
            else:
                last_line = self.move_listbox.get(tk.END)
                self.move_listbox.delete(tk.END)
                self.move_listbox.insert(tk.END, f"{last_line} {san_move}")
            self.move_listbox.see(tk.END)

    def game_loop(self):
        if self.board.is_game_over():
            self.draw_board()
            return

        move = self.current_player.get_move(self.board)
        if move:
            move_san = self.board.san(move)
            self.board.push(move)
            self.last_move = move
            self.capture_square = move.to_square if self.board.is_capture(move) else None
            self.current_player = self.white_player if self.board.turn else self.black_player
            self.update_move_list(move_san)
            self.update_scores()
            self.draw_board()
            self.update_turn_label()

        self.root.after_id = self.root.after(100, self.game_loop)

# ================= Run Example =================

if __name__ == "__main__":
    from src.player.lc0_bot_player import Lc0BotPlayer
    from src.player.random_bot_player import RandomBotPlayer
    from src.player.stockfish_bot_player import StockfishPlayer

    root = tk.Tk()
    root.geometry(f"{DEFAULT_WINDOW_WIDTH}x{DEFAULT_WINDOW_HEIGHT}")

    if random.choice([True, False]):
        white = Lc0BotPlayer("Leela", color=True, time_limit=1.0)
        black = StockfishPlayer("Stockfish", color=False, skill_level=10, time_limit=0.5)
        print("You are playing as White.")
    else:
        white = StockfishPlayer("Stockfish", color=True, skill_level=10, time_limit=0.5)
        black = Lc0BotPlayer("Leela", color=False, time_limit=1.0)
        print("You are playing as Black.")

    custom_dir = "/home/nate/projects/cs6640_project/data/pgn"
    gui = ChessGui(root, white, black, save_dir=custom_dir)

    root.mainloop()
