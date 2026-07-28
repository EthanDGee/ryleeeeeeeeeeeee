# Development Guide

Detailed setup and workflow for contributors.

## Quick Start

```bash
python -m venv .venv
source .venv/bin/activate  # Windows: .venv\Scripts\activate
pip install -e .
prek install
pytest packages/*/tests/ -v  # Verify setup
```

## Project Structure

```text
human-chessbot/
├── packages/
│   ├── play/       # Chess game application
│   ├── convert/    # PGN conversion utilities
│   └── train/      # ML training and dataset ETL
├── docs/           # Documentation
└── pyproject.toml  # Project configuration
```

## Workflow

### Local Development

```bash
git checkout -b feature/my-feature  # Create branch
# Make changes
git add .
git commit -m "feat: add feature X"  # Prek hooks run automatically
# If hooks fail, review changes and recommit
git push origin feature/my-feature
```

**Commit Message Format:**

This project uses conventional commits. Prefix your message with:

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation
- `style:` - Formatting (no code change)
- `refactor:` - Code refactoring
- `test:` - Tests
- `chore:` - Build, dependencies, etc.

Example: `feat(play): add move validation`

### Continuous Integration

All commits are checked by GitHub Actions:

- Code linting and formatting
- Type checking with MyPy
- Unit tests with coverage
- Security scanning
- Markdown validation

See `.github/workflows/` for configuration and
`.github/workflows/README.md` for details.

## Tools

### Prek - Pre-commit Framework

This project uses **prek**, a Rust-based pre-commit framework that's faster
than the standard pre-commit tool.

**Configuration:** `prek.toml`

**Installation:**

Prek is automatically installed when you run `pip install -e .`. To set up
Git hooks:

```bash
prek install
```

**Hooks Configured:**

- **Trailing whitespace** - Remove trailing spaces
- **YAML/JSON/TOML validation** - Check syntax
- **Large file detection** - Prevent committing large files
- **Merge conflict detection** - Catch unresolved conflicts
- **Debug statements** - Remove debug code before commit
- **Conventional commits** - Enforce commit format
- **Spell checking** - Detect typos in code and messages
- **Markdown linting** - Validate markdown syntax
- **Code formatting** - Prettier for SQL/TOML/YAML
- **Ruff linting** - Python linting with auto-fix
- **Type checking** - MyPy for type safety

**Manual Usage:**

```bash
prek run --all-files              # Run all hooks
prek run --files path/to/file.py  # Specific files
prek list-hooks                   # Show configured hooks
prek autoupdate                   # Update hook versions
prek uninstall                    # Remove Git hooks
```

**Stages:**

Hooks run at different Git stages:

- `pre-commit` - Before committing (default)
- `commit-msg` - Validate commit message
- `pre-push` - Before pushing to remote

Bypass hooks (not recommended):

```bash
git commit --no-verify            # Skip pre-commit hooks
git push --no-verify              # Skip pre-push hooks
```

### Testing

```bash
pytest packages/*/tests/ -v                  # All tests
pytest packages/play/tests/ -v               # Specific package
pytest --cov=packages --cov-report=html      # With coverage
pytest -k "test_name"                        # By pattern
pytest tests/path/test.py::test_function     # Specific test
```

**Coverage targets**: 80%+ overall, 90%+ for critical paths

### Linting & Formatting

All tools run automatically via prek hooks, but can be run manually:

```bash
ruff check packages/        # Check for linting issues
ruff check --fix packages/  # Auto-fix issues
ruff format packages/       # Format code
mypy packages/              # Type checking
```

**Note:** These are included in prek hooks, so commits will fail if
violations exist. Use `--no-verify` only as a last resort.

## Code Style

- **Line length**: 80 characters (markdown), 100 characters (Python)
- **Python**: 3.11+
- **Type hints**: Required for functions
- **Docstrings**: Required for public APIs
- **Formatting**: Ruff (replaces Black and isort)
- **Linting**: Ruff
- **Type checking**: MyPy

### Import Order

```python
# 1. Standard library
import logging
from typing import Optional

# 2. Third-party
import chess
from pydantic import BaseModel

# 3. Local
from packages.play.src.game.game import Game
```

### Docstrings (Google style)

```python
def apply_move(self, move: chess.Move) -> str:
    """Apply a move to the board.

    Args:
        move: Chess move to apply

    Returns:
        Move in SAN notation

    Raises:
        ValueError: If move is illegal
    """
```

## Common Tasks

### Add Dependency

Edit `pyproject.toml` dependencies section, then:

```bash
pip install -e .
```

### Add Module

1. Create module with type hints and docstrings
2. Write tests in `tests/` directory
3. Update package README

### Add Package

```bash
mkdir -p packages/mypackage/{src/mypackage,tests,docs}
touch packages/mypackage/README.md
```

## Troubleshooting

| Issue | Solution |
| ------- | ---------- |
| Import errors | `pip install -e .` |
| Test failures | `rm -rf .pytest_cache && pip install -e .` |
| Hooks not running | `prek install` to reinstall |
| Formatting conflicts | `prek run --all-files` |
| Type errors | `mypy packages/` to identify |
| Commit message fails | Use conventional format |
| Hook version mismatch | `prek autoupdate` to update |
| GitHub Actions fail | Check Actions tab logs |

### Common Prek Issues

**Hooks not installed:**

```bash
prek install
```

**Specific hook failing:**

```bash
prek run --hook-id <hook-name>  # Run single hook
prek list-hooks                 # See hook names
```

**Force commit despite failures:**

```bash
git commit --no-verify  # Not recommended!
```

### Testing Locally

Before pushing, run all checks locally:

```bash
prek run --all-files              # All pre-commit hooks
pytest packages/*/tests/ -v       # Tests
ruff check packages/              # Linting
mypy packages/                    # Type checking
```

## Resources

- [Prek](https://prek.dev/) - Fast pre-commit framework
- [Ruff](https://docs.astral.sh/ruff/) - Python linter & formatter
- [MyPy](http://mypy-lang.org/) - Static type checker
- [Pytest](https://docs.pytest.org/) - Testing framework
- [Conventional Commits](https://www.conventionalcommits.org/) - Commit format
- [GitHub Actions](https://github.com/features/actions) - CI/CD
- [Python Chess](https://python-chess.readthedocs.io/) - Chess library
