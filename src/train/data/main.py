from typing import List

from src.train.data.database.database import initialize_database
from src.train.data.database.files_metadata import save_files_metadata, files_metadata_exist
from src.train.data.database.games import save_games
from src.train.data.models.game import Game
from src.train.data.requesters.file_metadata import fetch_files_metadata, FileMetadata
from src.train.data.requesters.games import fetch_new_games


def main() -> None:
    initialize_database()

    # --- Fetch metadata if needed ---
    if not files_metadata_exist():
        print("Fetching metadata from Lichess...")
        files: List[FileMetadata] = fetch_files_metadata()
        print(f"Fetched {len(files)} files.\n")

        print("Saving metadata to database...")
        save_files_metadata(files)
        print("All metadata saved to database.\n")
    else:
        print("Files metadata already exists in the database.\n")

    games: List[Game] = fetch_new_games(max_files=10, max_size_gb=1)

    save_games(games)

    print("\nDone!")

if __name__ == "__main__":
    main()
