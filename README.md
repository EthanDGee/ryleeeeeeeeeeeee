# Rylee - The Human-like Chessbot

Rylee is a chess application suite with multiple engines, data processing tools, and ML training utilities.

## Features

- Interactive chess game (GUI and CLI)
- Human focused chess engine (Rylee)
- Multiple chess engines (Stockfish, LCZero, Random bot, and Rylee Bot)
- PGN file processing and conversion to CSV
- Lichess data pipeline for ML training
- Time controls and automatic game recording

## Quick Start

```bash
# Clone and setup
git clone <repository-url>
cd rylee
poetry install --no-root

# Setup the configs
cp packages/convert/.env.example packages/convert/.env
cp packages/play/.env.example packages/play/.env
cp packages/train/.env.example packages/train/.env
# See config.json if you want to tweak the training configs

# Build Legal Moves Database
python -m packages.train.src.dataset.fillers.fill_legal_moves

# Run chess application
poetry run python -m packages.train.src.dataset.main
poetry run python -m packages.play.src.main
poetry run poetry -m packages.train.src.train.main
```

## Project Structure

```
rylee/
├── packages/
│   ├── convert/       # PGN conversion utilities
│   ├── play/          # Chess game application
│   └── train/         # ML training and dataset pipeline
├── docs/              # Development documentation
└── pyproject.toml
```

## Packages

### [Play](packages/play/README.md)

Interactive chess game with Stockfish, LCZero, and random bot support.

### [Convert](packages/convert/README.md)

PGN file combination and conversion to CSV for ML training.

### [Train](packages/train/README.md)

ML training pipeline and Lichess dataset ETL.

