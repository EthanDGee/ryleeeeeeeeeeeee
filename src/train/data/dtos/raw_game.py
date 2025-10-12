from dataclasses import dataclass
from typing import Optional

@dataclass
class RawGame:
    file_id: Optional[int] = None  # Foreign key to file_metadata
    pgn: str = ""
