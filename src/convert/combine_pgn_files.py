import argparse
import os

def combine_pgn_files(pgn1_path, pgn2_path, output_path, delete_old=False):
    """
    Combines two PGN files into one and optionally deletes the originals.
    """
    # Check that both input files exist
    if not os.path.isfile(pgn1_path):
        raise FileNotFoundError(f"File not found: {pgn1_path}")
    if not os.path.isfile(pgn2_path):
        raise FileNotFoundError(f"File not found: {pgn2_path}")

    # Combine file contents
    with open(pgn1_path, "r", encoding="utf-8") as f1, open(pgn2_path, "r", encoding="utf-8") as f2:
        content1 = f1.read().strip()
        content2 = f2.read().strip()

    # Ensure a blank line between games if not already present
    combined = content1 + "\n\n" + content2 + "\n"

    # Write combined PGN
    os.makedirs(os.path.dirname(output_path), exist_ok=True)
    with open(output_path, "w", encoding="utf-8") as out:
        out.write(combined)

    print(f"Combined PGNs saved to: {output_path}")

    # Optionally delete old files
    if delete_old:
        os.remove(pgn1_path)
        os.remove(pgn2_path)
        print("Original files deleted.")

def main():
    parser = argparse.ArgumentParser(description="Combine two PGN files into one.")
    parser.add_argument("pgn1", help="Path to first PGN file")
    parser.add_argument("pgn2", help="Path to second PGN file")
    parser.add_argument("output", help="Path to output combined PGN file")
    parser.add_argument("--delete-old", action="store_true", help="Delete original PGN files after combining")

    args = parser.parse_args()

    combine_pgn_files(args.pgn1, args.pgn2, args.output, args.delete_old)

if __name__ == "__main__":
    main()
