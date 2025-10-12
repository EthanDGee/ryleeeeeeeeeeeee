import sqlite3
from typing import List

from src.train.data.database.config import DB_FILE, is_database_initialized
from src.train.data.models.file_metadata import FileMetadata

_TABLE_NAME: str = "files_metadata"


def initialize_files_metadata_table():
    with sqlite3.connect(DB_FILE) as conn:
        cursor = conn.cursor()
        cursor.execute(f"""
        CREATE TABLE IF NOT EXISTS {_TABLE_NAME} (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            url TEXT UNIQUE,
            filename TEXT,
            games INTEGER,
            size_gb REAL
        )
        """)
        conn.commit()


def files_metadata_exist() -> bool:
    """Return True if the 'files_metadata' table has any rows."""
    if not is_database_initialized():
        return False

    conn = sqlite3.connect(DB_FILE)
    cursor = conn.cursor()
    try:
        cursor.execute("SELECT 1 FROM files_metadata LIMIT 1;")
        has_rows = cursor.fetchone() is not None
    finally:
        conn.close()
    return has_rows


def save_files_metadata(files: List[FileMetadata]) -> None:
    """
    Save a list of FileMetadata objects into the database.
    Skips any files already present (based on URL).
    """
    with sqlite3.connect(DB_FILE) as conn:
        cursor = conn.cursor()

        for file in files:
            cursor.execute(f"SELECT 1 FROM {_TABLE_NAME} WHERE url = ?", (file.url,))
            if cursor.fetchone():
                continue  # skip existing

            cursor.execute(
                f"""
                INSERT INTO {_TABLE_NAME} (url, filename, games, size_gb)
                VALUES (?, ?, ?, ?)
                """,
                (file.url, file.filename, file.games, file.size_gb),
            )

        conn.commit()


def fetch_all_files_metadata() -> List[FileMetadata]:
    """
    Retrieve all entries from the files_metadata table as FileMetadata objects.
    """
    with sqlite3.connect(DB_FILE) as conn:
        cursor = conn.cursor()
        cursor.execute(f"SELECT url, filename, games, size_gb FROM {_TABLE_NAME}")
        rows = cursor.fetchall()

    return _rows_to_file_metadata(rows)


def fetch_files_metadata_under_size(max_gb: float) -> List[FileMetadata]:
    """
    Retrieve all entries from the files_metadata table with size less than max_gb.
    """
    with sqlite3.connect(DB_FILE) as conn:
        cursor = conn.cursor()
        cursor.execute(
            f"SELECT url, filename, games, size_gb FROM {_TABLE_NAME} WHERE size_gb < ?",
            (max_gb,),
        )
        rows = cursor.fetchall()

    return _rows_to_file_metadata(rows)


def _rows_to_file_metadata(rows: list[tuple]) -> List[FileMetadata]:
    """
    Helper function to convert database rows into FileMetadata objects.
    """
    return [
        FileMetadata(url=row[0], filename=row[1], games=row[2], size_gb=row[3])
        for row in rows
    ]
