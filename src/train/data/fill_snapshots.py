from src.train.data.database.database import initialize_database
from src.train.data.database.files_metadata import files_metadata_exist, save_file_metadata
from src.train.data.database.game_snapshots import save_snapshot, count_snapshots
from src.train.data.processers.game_snapshots import raw_game_to_snapshots
from src.train.data.requesters.file_metadata import fetch_files_metadata
from src.train.data.requesters.raw_games import fetch_new_raw_games


def fill_database_with_snapshots(
    snapshots_threshold: int = 100000,
    max_size_gb: float = 10,
    print_interval: int = 1000,
) -> None:
    """
    Fetch raw games from Lichess, convert to snapshots, and fill the database
    until the specified snapshot threshold is reached.

    Args:
        snapshots_threshold: Stop once this many snapshots are saved.
        max_size_gb: Only download raw games smaller than this size.
        print_interval: Print progress every N snapshots.
    """
    initialize_database()

    if not files_metadata_exist():
        print("Fetching metadata from Lichess...")
        for file in fetch_files_metadata():
            print(f"Saving metadata to database {file.filename}...")
            save_file_metadata(file)
        print("All metadata saved to database.\n")
    else:
        print("Files metadata already exists in the database.\n")

    print("Fetching and saving games...")

    snapshot_counter = count_snapshots()
    while snapshot_counter < snapshots_threshold:
        any_files_fetched = False

        for raw_game in fetch_new_raw_games(max_files=1, max_size_gb=max_size_gb):
            any_files_fetched = True
            if raw_game is None:
                continue

            for snapshot in raw_game_to_snapshots(raw_game):
                save_snapshot(snapshot)
                snapshot_counter += 1

                if print_interval and snapshot_counter % print_interval == 0:
                    print(f"Saved {snapshot_counter} snapshots so far...")

        if not any_files_fetched:
            print(f"No unprocessed files found under {max_size_gb} GB. Stopping.")
            break


if __name__ == "__main__":
    fill_database_with_snapshots(
        snapshots_threshold=100000,
        max_size_gb=10,
        print_interval=1000,
    )
