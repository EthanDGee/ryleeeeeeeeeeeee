import json
import os
import sys

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
    # If user provides a path, use it. Otherwise fall back to default.
    if len(sys.argv) == 2:
        config_path = sys.argv[1]
    else:
        print(f"No config path provided, using default: {DEFAULT_CONFIG_PATH}")
        config_path = DEFAULT_CONFIG_PATH

    config = load_config(config_path)

    neural_network = NeuralNetwork()
    trainer = Trainer(config, neural_network)

    trainer.random_search(config["num_iterations"])


if __name__ == "__main__":
    main()
