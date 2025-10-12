from src.train.data.database.config import is_database_initialized
from src.train.data.database.database import initialize_database
from src.train.data.database.files_metadata import save_files_metadata, fetch_files_metadata_under_size, files_metadata_exist
from src.train.data.requesters.file_metadata import fetch_files_metadata, FileMetadata
from typing import List

def main() -> None:
    # --- Initialize database if needed ---
    if not is_database_initialized():
        print("Initializing database...")
        initialize_database()
        print("Database initialized.\n")
    else:
        print("Database already initialized.\n")

    # --- Fetch metadata if needed ---
    if not files_metadata_exist():
        print("Fetching metadata from Lichess...")
        files: List[FileMetadata] = fetch_files_metadata()
        print(f"Fetched {len(files)} files.\n")

        print("Saving metadata to database...")
        save_files_metadata(files)
        print("All metadata saved to database.\n")
    else:
        print("Files metadata already exists in the database.\n")

    # Example usage: get all files under 1 GB
    max_size_gb = 1.0
    small_files = fetch_files_metadata_under_size(max_size_gb)
    print(f"Files under {max_size_gb} GB: {len(small_files)}")
    for file in small_files[:10]:
        print(f"{file.filename} ({file.size_gb} GB)")

    print("\nDone!")

if __name__ == "__main__":
    main()
