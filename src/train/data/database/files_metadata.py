import sqlite3
from typing import List

from src.train.data.database.config import DB_FILE, is_database_initialized
from src.train.data.models.file_metadata import FileMetadata

_TABLE_NAME: str = "files_metadata"


def create_files_metadata_table():
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


def save_file_metadata(file: FileMetadata) -> None:
    """
    Save a single FileMetadata object into the database.
    Skips if a file with the same URL already exists.
    Updates the object's `id` after insertion.
    """
    with sqlite3.connect(DB_FILE) as conn:
        cursor = conn.cursor()
        cursor.execute(f"SELECT id FROM {_TABLE_NAME} WHERE url = ?", (file.url,))
        row = cursor.fetchone()
        if row:
            file.id = row[0]  # assign existing id
            return

        cursor.execute(
            f"""
            INSERT INTO {_TABLE_NAME} (url, filename, games, size_gb)
            VALUES (?, ?, ?, ?)
            """,
            (file.url, file.filename, file.games, file.size_gb),
        )
        file.id = cursor.lastrowid  # assign new id
        conn.commit()


def save_files_metadata(files: List[FileMetadata]) -> None:
    """
    Save a list of FileMetadata objects into the database.
    Calls `save_file_metadata` for each file.
    """
    for file in files:
        save_file_metadata(file)


def fetch_all_files_metadata() -> List[FileMetadata]:
    """
    Retrieve all entries from the files_metadata table as FileMetadata objects.
    """
    with sqlite3.connect(DB_FILE) as conn:
        cursor = conn.cursor()
        cursor.execute(f"SELECT id, url, filename, games, size_gb FROM {_TABLE_NAME}")
        rows = cursor.fetchall()

    return [
        _row_to_file_metadata(row)
        for row in rows
    ]


def fetch_files_metadata_under_size(max_gb: float) -> List[FileMetadata]:
    """
    Retrieve all entries from the files_metadata table with size less than max_gb.
    """
    with sqlite3.connect(DB_FILE) as conn:
        cursor = conn.cursor()
        cursor.execute(
            f"SELECT id, url, filename, games, size_gb FROM {_TABLE_NAME} WHERE size_gb < ?",
            (max_gb,),
        )
        rows = cursor.fetchall()

    return [
        _row_to_file_metadata(row)
        for row in rows
    ]



def _row_to_file_metadata(row: tuple) -> FileMetadata:
    return FileMetadata(
            id=row[0],
            url=row[1],
            filename=row[2],
            games=row[3],
            size_gb=row[4]
        )
