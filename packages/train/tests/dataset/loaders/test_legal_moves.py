import sqlite3

import pytest
import torch

from packages.train.src.dataset.loaders.legal_moves import LegalMovesDataset


class TestLegalMovesDataset:
    """Tests for LegalMovesDataset class."""

    @pytest.fixture
    def temp_database(self, tmp_path):
        """Fixture to create a temporary database for testing."""
        db_path = tmp_path / "test_legal_moves.db"
        connection = sqlite3.connect(db_path)
        cursor = connection.cursor()

        # Create table
        cursor.execute(
            """
            CREATE TABLE legal_moves
            (
                move  TEXT,
                types TEXT
            )
            """
        )

        # Insert sample data
        cursor.executemany(
            """
            INSERT INTO legal_moves (move, types)
            VALUES (?, ?)
            """,
            [("e2e4", "pawn"), ("g1f3", "knight"), ("e7e8q", "pawn,queen"), ("g8=Q+", "pawn")],
        )

        connection.commit()
        connection.close()
        return db_path

    def test_load_data(self, temp_database):
        """Test that _load_data correctly fetches and processes data."""
        dataset = LegalMovesDataset(db_path=temp_database)
        assert len(dataset.data) == 4
        assert dataset.data[0] == {"move": "e2e4", "piece_types": ["pawn"]}

    def test_vocab_building(self, temp_database):
        """Test whether vocabulary creation produces expected results."""
        dataset = LegalMovesDataset(db_path=temp_database)
        vocab = dataset.vocab

        expected_vocab = {"e2e4": 0, "g1f3": 1, "e7e8q": 2, "g8=Q+": 3}
        assert vocab == expected_vocab

    def test_encode_move(self, temp_database):
        """Test move encoding matches expected tensor."""
        dataset = LegalMovesDataset(db_path=temp_database)

        encoded_move = dataset._encode_move("e2e4")
        assert torch.allclose(
            encoded_move, torch.tensor([4 / 7, 1 / 7, 4 / 7, 3 / 7], dtype=torch.float32)
        )

        encoded_promotion = dataset._encode_move("e7e8q")
        assert torch.allclose(
            encoded_promotion,
            torch.tensor([4 / 7, 6 / 7, 4 / 7, 7 / 7, 4 / 4], dtype=torch.float32),
        )

    def test_encode_piece_types(self, temp_database):
        """Test piece type encoding produces one-hot encoded tensor."""
        dataset = LegalMovesDataset(db_path=temp_database)
        encoded = dataset._encode_piece_types(["pawn", "queen"])

        expected = torch.tensor([1.0, 0.0, 0.0, 0.0, 1.0, 0.0], dtype=torch.float32)
        assert torch.allclose(encoded, expected)

    def test_get_move_from_index(self, temp_database):
        """Test fetching move string from index."""
        dataset = LegalMovesDataset(db_path=temp_database)

        move = dataset.get_move_from_index(0)
        assert move == "e2e4"

        move = dataset.get_move_from_index(1)
        assert move == "g1f3"

        move = dataset.get_move_from_index(999)
        assert move == ""

    def test_get_index_from_move(self, temp_database):
        """Test fetching vocabulary index from move."""
        dataset = LegalMovesDataset(db_path=temp_database)

        index = dataset.get_index_from_move("e2e4")
        assert index == 0

        index = dataset.get_index_from_move("g1f3")
        assert index == 1

        # make sure that pawn promotions are handled properly
        index = dataset.get_index_from_move("g8=Q+")
        assert index == 3

        index = dataset.get_index_from_move("nonexistent")
        assert index == -1

    def test_len_method(self, temp_database):
        """Test that __len__ method returns correct data length."""
        dataset = LegalMovesDataset(db_path=temp_database)
        assert len(dataset) == 4

    def test_getitem_method(self, temp_database):
        """Test that __getitem__ returns properly encoded data."""
        dataset = LegalMovesDataset(db_path=temp_database)

        sample = dataset[0]
        assert sample["move"] == "e2e4"
        assert sample["move_index"] == 0
        assert torch.allclose(
            sample["move_encoding"], torch.tensor([4 / 7, 1 / 7, 4 / 7, 3 / 7], dtype=torch.float32)
        )
        assert torch.allclose(
            sample["piece_encoding"],
            torch.tensor([1.0, 0.0, 0.0, 0.0, 0.0, 0.0], dtype=torch.float32),
        )

    def test_custom_vocab(self, temp_database):
        """Test dataset with a predefined vocabulary."""
        custom_vocab = {"e2e4": 100, "g1f3": 101}
        dataset = LegalMovesDataset(db_path=temp_database, vocab=custom_vocab)

        assert dataset.vocab == custom_vocab
        assert dataset.get_move_from_index(100) == "e2e4"
        assert dataset.get_index_from_move("g1f3") == 101
