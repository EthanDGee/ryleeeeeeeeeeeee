import json
import os
import sys

import torch

from packages.train.src.models.neural_network import NeuralNetwork
from packages.train.src.train.trainer import Trainer

DEFAULT_CONFIG_PATH = "./config.json"


def print_usage():
    print("Usage: python train.py [config_path]")
    print("If no config_path is provided, it will try to load './config.json'.")


def load_config(path: str):
    if not os.path.isfile(path):
        print(f"Error: Config file '{path}' does not exist.")
        sys.exit(1)

    try:
        with open(path) as f:
            return json.load(f)
    except Exception as e:
        print(f"Error loading JSON from '{path}': {e}")
        sys.exit(1)


def main():
    # If user provides a path, use it. Other wise fall back to default.
    if len(sys.argv) <= 1:
        print(f"No config path provided, using default: {DEFAULT_CONFIG_PATH}")
        config_path = DEFAULT_CONFIG_PATH
    else:
        config_path = sys.argv[1]

    config = load_config(config_path)

    starting_model = NeuralNetwork()

    # load a user specified model as the starting point for further training
    if len(sys.argv) > 2:
        print("Loading base model...")
        base_model_path = sys.argv[2]
        map_location = torch.device("cuda") if config["cuda_enabled"] else torch.device("cpu")

        starting_model.load_state_dict(
            torch.load(base_model_path, weights_only=True, map_location=map_location)
        )
        print(f"Successfully loaded {base_model_path}")

    trainer = Trainer(config, starting_model)

    trainer.random_search(config["num_iterations"])


if __name__ == "__main__":
    main()
