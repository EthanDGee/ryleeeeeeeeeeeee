import requests
import zstandard as zstd
from src.train.data.database.files_metadata import fetch_files_metadata_under_size
from src.train.data.database.games import fetch_file_ids_in_db
from src.train.data.models.game import Game
from typing import List


def fetch_new_games(max_files: int = 5, max_size_gb: float = 1) -> List[Game]:
    """
    Fetch new games from Lichess PGN files that have not been downloaded yet.
    """
    candidate_files = fetch_files_metadata_under_size(max_gb=max_size_gb)
    existing_file_ids = fetch_file_ids_in_db()
    files_to_download = [f for f in candidate_files if f.id not in existing_file_ids][:max_files]

    all_games: List[Game] = []

    for file_meta in files_to_download:
        print(f"Downloading and decompressing PGN file: {file_meta.filename}...")
        response = requests.get(file_meta.url, stream=True)
        if response.status_code == 200:
            decompressor = zstd.ZstdDecompressor()
            with decompressor.stream_reader(response.raw) as reader:
                decompressed_text = reader.read().decode("utf-8")

            individual_pgns = _split_pgn_text_into_games(decompressed_text)
            for pgn in individual_pgns:
                all_games.append(Game(file_id=file_meta.id, pgn=pgn))
        else:
            print(f"Failed to download {file_meta.filename} (status {response.status_code})")

    return all_games


def _split_pgn_text_into_games(pgn_text: str) -> List[str]:
    """
    Split a single PGN file into individual games.
    Each game starts with '[Event '.
    """
    raw_games = pgn_text.strip().split("\n\n[Event ")
    pgn_list = []
    for i, raw in enumerate(raw_games):
        if i > 0:
            raw = "[Event " + raw
        pgn_list.append(raw.strip())
    return pgn_list
