package parser

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"server/rest/internal/config"
	"server/rest/internal/database"
	"server/rest/internal/models"
	"server/rest/internal/utils"

	"github.com/corentings/chess/v2"
	"github.com/klauspost/compress/zstd"
)

func DecompressPGN(metadata models.Metadata) {
	compressedPath := filepath.Join(config.FILES_FOLDER, metadata.Filename)
	decompressedPath := strings.TrimSuffix(compressedPath, ".zst")

	in, err := os.Open(compressedPath)
	if err != nil {
		log.Fatalf("failed to open %s: %v", compressedPath, err)
	}
	defer utils.Close(in, "compressed file")

	decoder, err := zstd.NewReader(in)
	if err != nil {
		log.Fatalf("failed to create zstd reader for %s: %v", compressedPath, err)
	}
	defer decoder.Close()

	out, err := os.Create(decompressedPath)
	if err != nil {
		log.Fatalf("failed to create %s: %v", decompressedPath, err)
	}
	defer utils.Close(out, "decompressed file")

	if _, err := io.Copy(out, decoder); err != nil {
		log.Fatalf("failed to decompress %s: %v", compressedPath, err)
	}
}

func ProcessFile(metadata models.Metadata) error {
	path := strings.TrimSuffix(filepath.Join(config.FILES_FOLDER, metadata.Filename), ".zst")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		DecompressPGN(metadata)
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer utils.Close(file, "file")

	scanner := chess.NewScanner(file)

	preprocessedGames := metadata.Processed
	log.Printf("processing %s: %d games total, %d already processed", metadata.Filename, metadata.Games, preprocessedGames)

	gameIndex := 0
	for scanner.HasNext() {
		if gameIndex < preprocessedGames {
			gameIndex++
			continue
		}

		game, err := scanner.ParseNext()
		if err != nil {
			log.Fatalf("Failed to process game: %d from %s", gameIndex, metadata.Filename)
		}

		err = database.InsertGame(*game, metadata)
		if err != nil {
			log.Printf("failed to insert game %d from %s: %v", gameIndex, metadata.Filename, err)
			return err
		}
		log.Printf("inserted game %d from %s", gameIndex, metadata.Filename)

		gameIndex++
	}

	log.Printf("finished processing %s: %d games processed", metadata.Filename, gameIndex)
	return nil
}
