from unittest.mock import patch

import pytest

from packages.train.src.dataset.loaders.game_snapshots import GameSnapshotsDataset


class TestGameSnapshotsDataset:
    """Tests for the GameSnapshotsDataset class."""

    @patch("packages.train.src.dataset.loaders.game_snapshots.count_snapshots")
    def test_num_indexes_exceeds_database_count(self, mock_count_snapshots):
        """Test ValueError raised when num_indexes exceeds database snapshot count."""
        mock_count_snapshots.return_value = 100  # Mock snapshot count
        with pytest.raises(ValueError, match="num_indexes is larger than the size of the database"):
            GameSnapshotsDataset(start_index=50, num_indexes=60, db_path=":memory:")

    def test_fen_to_tensor(self):
        """Test conversion of FEN string to tensor."""
        dataset = GameSnapshotsDataset(0, 1, db_path=":memory:")
        fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"  # Starting position
        tensor = dataset._fen_to_tensor(fen)
        assert tensor.shape.numel() == (8 * 8 * 12)
        assert tensor.sum().item() == 32  # 32 pieces on the board

    def test_encode_result(self):
        """Test encoding of game results based on the player's perspective."""
        dataset = GameSnapshotsDataset(0, 1, db_path=":memory:")
        result_win = dataset._encode_result("1-0", "w")
        result_draw = dataset._encode_result("1/2-1/2", "b")
        result_loss = dataset._encode_result("0-1", "w")
        assert result_win.item() == 1.0
        assert result_draw.item() == 0.5
        assert result_loss.item() == 0.0

    def test_encode_turn(self):
        """Test encoding of player turn as one-hot vector."""
        dataset = GameSnapshotsDataset(0, 1, db_path=":memory:")
        turn_white = dataset._encode_turn("w")
        turn_black = dataset._encode_turn("b")
        assert list(turn_white.numpy()) == [1.0, 0.0]
        assert list(turn_black.numpy()) == [0.0, 1.0]

    def test_normalize_elo(self):
        """Test normalization of ELO ratings."""
        dataset = GameSnapshotsDataset(0, 1, db_path=":memory:")
        result = dataset._normalize_elo(2000, 1500)
        assert result.shape == (2,)
        assert pytest.approx(result[0].item(), 0.001) == (2000 - 1638.43153) / 185.80054702756055
        assert pytest.approx(result[1].item(), 0.001) == (1500 - 1638.43153) / 185.80054702756055

    def test_encode_move_valid(self):
        """Test encoding of valid chess move."""
        dataset = GameSnapshotsDataset(0, 1, db_path=":memory:")
        fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
        move_san = "e4"
        move, promo = dataset._encode_move(fen, move_san)
        assert promo == 0  # No promotion
        assert isinstance(move, int)

    def test_encode_move_invalid(self):
        """Test encoding of invalid chess move."""
        dataset = GameSnapshotsDataset(0, 1, db_path=":memory:")
        fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
        move_san = "invalid_move"
        move, promo = dataset._encode_move(fen, move_san)
        assert move == 0
        assert promo == 0

    def test_len_method(self):
        """Test the __len__ method for dataset size."""
        dataset = GameSnapshotsDataset(start_index=0, num_indexes=10, db_path=":memory:")
        assert len(dataset) == 10
