import re
import requests
from urllib.parse import urljoin
from typing import List

from src.train.data.models.file_metadata import FileMetadata

_BASE_URL = "https://database.lichess.org/standard/"
_COUNTS_URL = urljoin(_BASE_URL, "counts.txt")


def fetch_files_metadata() -> List[FileMetadata]:
    """
    Fetch metadata about all standard Lichess files.

    Returns:
        List[FileMetadata]: A list of metadata objects for each .pgn.zst file.
    """
    # --- Fetch counts.txt to get game counts ---
    counts_resp = requests.get(_COUNTS_URL)
    counts_resp.raise_for_status()

    counts: dict[str, int] = {}
    for line in counts_resp.text.strip().splitlines():
        parts = line.split()
        if len(parts) == 2:
            filename, games = parts
            counts[filename] = int(games.replace(',', ''))

    # --- Fetch standard directory page ---
    resp = requests.get(_BASE_URL)
    resp.raise_for_status()
    html = resp.text

    # --- Extract .pgn.zst files ---
    file_names = re.findall(r'href="(lichess_db_standard_rated_[^"]+\.pgn\.zst)"', html)
    metadata_list: List[FileMetadata] = []

    for filename in file_names:
        file_url = urljoin(_BASE_URL, filename)

        # Get file size from HEAD request
        head_resp = requests.head(file_url)
        size_gb = int(head_resp.headers.get("Content-Length", 0)) / (1024 ** 3)

        metadata_list.append(
            FileMetadata(
                url=file_url,
                filename=filename,
                games=counts.get(filename, 0),
                size_gb=round(size_gb, 2),
            )
        )

    return metadata_list

if __name__ == "__main__":
    files_metadata = fetch_files_metadata()
    print(f"Found {len(files_metadata)} standard files.\n")
    for file in files_metadata[:10]:  # just print first 10 for preview
        print(file)
