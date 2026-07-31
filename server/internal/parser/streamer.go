package parser

import (
	"log"
	"os"

	"server/rest/internal/config"
	"server/rest/internal/database"
	"server/rest/internal/models"
	"server/rest/internal/utils"

	"github.com/corentings/chess/v2"
)

func ProcessFile(metadata models.Metadata) error {
	filepath := config.FILES_FOLDER + metadata.Filename

	file, err := os.Open(filepath)
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
			log.Fatalf("Failed to process game: %d in %s", gameIndex, metadata.Filename)
		}

		err = database.InsertGame(*game, metadata)
		if err != nil {
			log.Printf("failed to insert game %d in %s: %v", gameIndex, metadata.Filename, err)
			return err
		}
		log.Printf("inserted game %d in %s", gameIndex, metadata.Filename)

		gameIndex++
	}

	log.Printf("finished processing %s: %d games processed", metadata.Filename, gameIndex)
	return nil
}
