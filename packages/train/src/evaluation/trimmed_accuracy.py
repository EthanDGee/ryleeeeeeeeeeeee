import os
import sqlite3

import torch

from packages.train.src.constants import DB_FILE
from packages.train.src.dataset.loaders.game_snapshots import GameSnapshotsDataset
from packages.train.src.dataset.repositories.processed_snapshots import count_processed_snapshots
from packages.train.src.models.neural_network import NeuralNetwork


def trimmed_accuracy(model_path: str):
    """Calculate model accuracy metrics on qualified game snapshots."""

    test_set = GameSnapshotsDataset(
        int(count_processed_snapshots() * 0.8), int(count_processed_snapshots() * 0.2)
    )
    model = NeuralNetwork()
    model.load_state_dict(
        torch.load(model_path, weights_only=True, map_location=torch.device("cpu"))
    )
    test_set_start = int(count_processed_snapshots() * 0.8)
    test_set_end = test_set_start + int(count_processed_snapshots() * 0.2)
    print(test_set_start, test_set_end)

    correct_moves_top1 = 0
    correct_moves_top5 = 0
    total_snapshots = 0

    with sqlite3.connect(DB_FILE) as conn:
        cur = conn.cursor()

        # get snapshots from test set that meet the maia papers requirements
        cur.execute(
            """
            SELECT gs.id FROM game_snapshots gs
            JOIN game_statistics gst ON gs.raw_game_id = gst.raw_game_id
            WHERE gs.id >= ? AND gs.id < ?
            AND gst.event = 'Rated Blitz game'
            AND gst.white_elo > 1100 AND gst.white_elo < 1900
            AND gst.black_elo > 1100 AND gst.black_elo < 1900
            AND ABS(gst.black_elo - gst.white_elo) < 200
            AND gs.move_number > 10
        """,
            (test_set_start, test_set_end),
        )

        snapshots = cur.fetchall()
        total_snapshots = len(snapshots)
        tensors = [test_set[i[0] - test_set_start] for i in snapshots]

        if not tensors:
            print("No qualifying snapshots found in test set.")
        else:
            print(f"Evaluating {len(tensors)} snapshots...")

            # Process tensors in batch
            model.eval()
            with torch.no_grad():
                boards = torch.stack([t[0][0] for t in tensors])
                metadata = torch.stack([t[0][1] for t in tensors])
                chosen_moves = torch.tensor([t[1][0] for t in tensors], dtype=torch.long)

                predicted_moves, _ = model(metadata, boards)

                # Calculate top-1 accuracy
                _, predicted_moves_top1 = torch.max(predicted_moves.data, 1)
                correct_moves_top1 += (predicted_moves_top1 == chosen_moves).sum().item()

                # Calculate top-5 accuracy
                _, predicted_moves_top5 = torch.topk(predicted_moves.data, k=5, dim=1)
                correct_moves_top5 += (
                    (predicted_moves_top5 == chosen_moves.unsqueeze(1)).any(dim=1).sum().item()
                )

    top_1_accuracy = correct_moves_top1 / total_snapshots * 100
    top_5_accuracy = correct_moves_top5 / total_snapshots * 100
    print(f"Top-1 accuracy: {top_1_accuracy:.2f}%")
    print(f"Top-5 accuracy: {top_5_accuracy:.2f}%")
    print(f"Total snapshots: {total_snapshots}")
    return top_1_accuracy, top_5_accuracy, total_snapshots


if __name__ == "__main__":
    model_path = os.path.expanduser("~/rylee.pth")
    trimmed_accuracy(model_path)
