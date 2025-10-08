import sqlite3
from typing import List


class Database:
    def __init__(self, db_name: str):
        self.db_name = db_name

        # HEADER INFO
        self.INTEGER_LABELS: List[str] = ["white_elo", "black_elo"]
        self.BOOLEAN_LABELS: List[str] = ["is_black"]
        # Columns exactly as in the provided selection
        self.CLASS_LABELS: List[str] = [
            # board squares
            "A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8",
            "B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8",
            "C1", "C2", "C3", "C4", "C5", "C6", "C7", "C8",
            "D1", "D2", "D3", "D4", "D5", "D6", "D7", "D8",
            "E1", "E2", "E3", "E4", "E5", "E6", "E7", "E8",
            "F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8",
            "G1", "G2", "G3", "G4", "G5", "G6", "G7", "G8",
            "H1", "H2", "H3", "H4", "H5", "H6", "H7", "H8",
        ]
        self.TARGET_COLUMN: str = "selected_move"

        self._ensure_table_exist()

    def _ensure_table_exist(self) -> None:
        def _append_columns(headers: List[str], data_type: str) -> str:
            new_columns = ""
            for header in headers:
                new_columns += f"{header} {data_type}, "

            return new_columns

        """Create tables if they don't exist"""
        with sqlite3.connect(self.db_name) as conn:
            cur = conn.cursor()

            query = "CREATE TABLE IF NOT EXISTS state ("

            # add the various labels
            query += _append_columns(self.INTEGER_LABELS, "INTEGER")
            query += _append_columns(self.BOOLEAN_LABELS, "INTEGER")
            query += _append_columns(self.CLASS_LABELS, "INTEGER")  # classes are stored as indexes

            # add the target class
            query += f"{self.TARGET_COLUMN} TEXT)"
            cur.execute(query)
            conn.commit()

    def _save_state(self, state_data: List) -> None:

        # try and convert all labels to ints if any fail ignore given state and print warning
        labels = state_data[:-1]
        for i in range(len(labels)):
            try:
                labels[i] = int(labels[i])
            except ValueError:
                print(f"Warning: Could not convert label {i} - {labels[i]} to int. Ignoring state.")

        # format for SQL insert statement
        processed_data = ""
        for label in labels:
            processed_data += f"{label}, "

        # add target value

        processed_data += f"'{state_data[-1]}'"

        with sqlite3.connect(self.db_name) as conn:
            cur = conn.cursor()
            cur.execute(f"INSERT INTO state VALUES ({processed_data})")
            conn.commit()

    def import_csv(self, csv_path: str) -> None:
        """Import a CSV file into the database"""
        with open(csv_path, "r") as file:
            file.readline()  # skip the header line
            line = file.readline()
            while line:
                data = line.split(",")
                self._save_state(data)
                line = file.readline()
