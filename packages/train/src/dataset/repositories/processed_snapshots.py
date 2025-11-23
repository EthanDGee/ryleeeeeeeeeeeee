import sqlite3

from packages.train.src.constants import DB_FILE

_TABLE_NAME = "processed_snapshots"


def create_processed_snapshots_table():
    """Create the 'processed_snapshots' table if it does not exist.

    This table caches processed game_snapshots data to avoid
    processing the same row twice.
    """
    conn = sqlite3.connect(DB_FILE)
    c = conn.cursor()

    c.execute(
        f"""
    CREATE TABLE IF NOT EXISTS {_TABLE_NAME} (
        snapshot_id INTEGER PRIMARY KEY,
        board BLOB NOT NULL,
        metadata BLOB NOT NULL,
        chosen_move INTEGER NOT NULL,
        valid_moves BLOB NOT NULL,
        FOREIGN KEY(snapshot_id) REFERENCES game_snapshots(id)
    )
    """
    )
    conn.commit()
    conn.close()


def save_processed_snapshots(data: list[tuple[int, bytes, bytes, int, bytes]]):
    """Save multiple processed snapshots in a single transaction.

    Args:
        data: List of (snapshot_id, board_bytes, metadata_bytes, chosen_move, valid_moves_bytes)
    """
    if not data:
        return

    with sqlite3.connect(DB_FILE) as conn:
        c = conn.cursor()
        c.executemany(
            f"INSERT OR IGNORE INTO {_TABLE_NAME} (snapshot_id, board, metadata, chosen_move, valid_moves) VALUES (?, ?, ?, ?, ?)",
            data,
        )
        conn.commit()


def get_processed_snapshots(
    snapshot_ids: list[int],
) -> dict[int, tuple[bytes, bytes, int, bytes]]:
    """Get already processed snapshots from cache.

    Args:
        snapshot_ids: List of snapshot IDs to look up

    Returns:
        Dict mapping snapshot_id -> (board_bytes, metadata_bytes, chosen_move, valid_moves_bytes)
    """
    if not snapshot_ids:
        return {}

    with sqlite3.connect(DB_FILE) as conn:
        placeholders = ",".join("?" * len(snapshot_ids))
        c = conn.cursor()
        c.execute(
            f"SELECT snapshot_id, board, metadata, chosen_move, valid_moves FROM {_TABLE_NAME} WHERE snapshot_id IN ({placeholders})",
            snapshot_ids,
        )
        return {row[0]: (row[1], row[2], row[3], row[4]) for row in c.fetchall()}
