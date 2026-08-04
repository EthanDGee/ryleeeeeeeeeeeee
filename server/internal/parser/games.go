package parser

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"server/rest/internal/config"
	"server/rest/internal/database"
	"server/rest/internal/models"
	"server/rest/internal/utils"

	"github.com/corentings/chess/v2"
)

func ProcessFile(metadata models.Metadata, cap int) error {
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

	nextPlyID, err := database.NextPlyID()
	if err != nil {
		return err
	}

	gameIndex := 0
	batch := make([]chess.Game, 0, config.BATCH_SIZE)
	for scanner.HasNext() && gameIndex < cap {
		if gameIndex < preprocessedGames {
			gameIndex++
			_, err := scanner.ParseNext()
			if err != nil {
				log.Fatalf("Failed to process game: %d from %s", gameIndex, metadata.Filename)
			}
			continue
		}

		game, err := scanner.ParseNext()
		if err != nil {
			log.Fatalf("Failed to process game: %d from %s", gameIndex, metadata.Filename)
		}

		batch = append(batch, *game)
		gameIndex++

		if len(batch) >= config.BATCH_SIZE {
			nextPlyID, err = database.InsertGames(batch, metadata, nextPlyID)
			if err != nil {
				log.Printf("failed to insert batch ending at game %d from %s: %v", gameIndex, metadata.Filename, err)
				return err
			}
			if gameIndex%5000 == 0 {
				log.Printf("inserted game %d from %s", gameIndex, metadata.Filename)
			}
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		if _, err := database.InsertGames(batch, metadata, nextPlyID); err != nil {
			log.Printf("failed to insert final batch from %s: %v", metadata.Filename, err)
			return err
		}
	}

	log.Printf("finished processing %s: %d games processed", metadata.Filename, gameIndex)
	return nil
}

func EnsureNGamesProcessed(n int) {
	processedGames, err := database.CountProcessed()
	if err != nil {
		log.Fatal(err)
	}

	remainingGames := n - processedGames

	for remainingGames > 0 {
		nextFile, err := database.GetSmallestUnprocessedFile()
		if err != nil {
			log.Fatal(err)
		}

		err = ProcessFile(nextFile, nextFile.Processed+remainingGames)
		if err != nil {
			log.Fatal(err)
		}

		processedGames, err = database.CountProcessed()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("processed games: %d/%d (%.1f%%)", processedGames, n, float64(processedGames)/float64(n)*100)
		remainingGames = n - processedGames
	}

	log.Printf("finished ensuring %d games processed\n", n)
}
