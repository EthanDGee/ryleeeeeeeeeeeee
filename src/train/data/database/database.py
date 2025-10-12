import sqlite3
from typing import List, Callable

from src.train.data.database.config import DB_FILE
from src.train.data.database.files_metadata import create_files_metadata_table
from src.train.data.database.games import create_games_table

# List of functions that create tables in the database
TABLE_CREATORS: List[Callable[[], None]] = [
    create_files_metadata_table,
    create_games_table
]

def initialize_database() -> None:
    """
    Creates the SQLite database and the tables if they don't exist.
    """
    # Connecting ensures the file exists
    with sqlite3.connect(DB_FILE):
        pass

    # Initialize tables
    for table_creator in TABLE_CREATORS:
        table_creator()
