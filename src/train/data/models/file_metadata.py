from dataclasses import dataclass
from typing import Optional

@dataclass
class FileMetadata:
    url: str
    filename: str
    games: int
    size_gb: float
    id: Optional[int] = None  # DB primary key
    processed: bool = False   # New field
