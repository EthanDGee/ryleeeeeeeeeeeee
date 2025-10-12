from dataclasses import dataclass

@dataclass
class FileMetadata:
    url: str
    filename: str
    games: int
    size_gb: float
