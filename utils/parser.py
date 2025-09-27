import chess
import chess.pgn
import sys


def convert_to_csv(file_path: str):
    print(file_path)
    pgn = open(file_path)

    # create the header for the csv file
    header = "white_elo,black_elo,blacks_move,"

    for letter in "ABCDEFGH":
        for num in range(1, 9):
            header += f"{letter}{num},"

    csv_name = file_path.split(".")[0] + ".csv"
    with open(csv_name, mode="w") as csv:
        csv.write(f"{header}\n")

        current_game = chess.pgn.read_game(pgn)

        while current_game is not None:
            # get game info
            elo = (current_game.headers["WhiteElo"], current_game.headers["BlackElo"])
            board = current_game.board()
            is_black = False

            # loop through moves and save board states to csv
            for move in current_game.mainline_moves():
                board.push(move)
                row = convert_board_to_row(board.fen(), elo, is_black)

                csv.write(row)
                is_black = not is_black  # toggle every move

            current_game = chess.pgn.read_game(pgn)


def convert_board_to_row(fen: str, elo: tuple, is_black: bool) -> str:
    piece_map = {
        "P": 1,
        "N": 2,
        "B": 3,
        "R": 4,
        "Q": 5,
        "K": 6,
        "p": 7,
        "n": 8,
        "b": 9,
        "r": 10,
        "q": 11,
        "k": 12,
    }

    board = []
    for char in fen:  # skip line splits
        if char == "/":
            continue
        # digits indicate sequences of empty spaces so we simply make the next 'n' spaces 0
        elif char.isdigit():
            board.extend([0] * int(char))
        else:
            board.append(piece_map.get(char, 0))

    return (
        f"{','.join(elo)},{int(is_black)},{','.join([str(piece) for piece in board])}\n"
    )


if __name__ == "__main__":
    print(sys.argv)

    # get file_path
    path = sys.argv[1]
    convert_to_csv(path)
