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

// ProcessFile parses PGNs from the file specified by metadata and inserts the resulting
// game data into the database while notating progress and corruption issues.
//
// Continues until either the file is fully processed or the global cap of total.processed
// games (globally) has been reached.
func ProcessFile(metadata models.Metadata, cap int) error {
	path := strings.TrimSuffix(filepath.Join(config.FILES_FOLDER, metadata.Filename), ".zst")

	// decompress file to prepare for parsing
	if _, err := os.Stat(path); os.IsNotExist(err) {
		DecompressPGN(metadata)
	}

	// open decompressed file
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	scanner := chess.NewScanner(file)

	preprocessedGames := metadata.Processed
	log.Printf("processing %s: %d games total, %d already processed", metadata.Filename, metadata.Games, preprocessedGames)

	// get the starting point for ply index to use as base ply index
	nextPlyID, err := database.NextPlyID()
	if err != nil {
		return err
	}

	// process games while there are still games to process or until the
	// required amount of games have been processed
	gameIndex := 0
	batch := make([]chess.Game, 0, config.BATCH_SIZE)
	for scanner.HasNext() && gameIndex < cap {

		// skip games that have already been processed to avoid duplicate entries
		if gameIndex < preprocessedGames {
			gameIndex++
			_, err := scanner.ParseNext()

			if err != nil {
				log.Fatalf("Failed to process game: %d from %s", gameIndex, metadata.Filename)
			}

			continue
		}

		game, err := scanner.ParseNext()

		// increment corrupted if the PGN was unable to be parsed
		if err != nil {
			log.Printf("Failed to process game: %d from %s\n", gameIndex, metadata.Filename)

			err = database.IncrementCorrupted(metadata.Id)
			if err != nil {
				log.Fatal(err)
			}

			err = database.IncrementProcessed(metadata.Id)
			if err != nil {
				log.Fatal(err)
			}
			// since we failed to parse we don't increment the
			// index and move onto the next game.
			continue
		}

		// add game to batch and insert batch if batch size has been reached
		batch = append(batch, *game)
		gameIndex++

		if len(batch) >= config.BATCH_SIZE {
			nextPlyID, err = database.InsertGames(batch, metadata, nextPlyID)
			if err != nil {
				log.Printf("failed to insert batch ending at game %d from %s: %v", gameIndex, metadata.Filename, err)
				return err
			}
			log.Printf("inserted game %d from %s", gameIndex, metadata.Filename)
			batch = batch[:0]
		}
	}

	utils.Close(file, "file")

	// insert remaining batch entries
	if len(batch) > 0 {
		if _, err := database.InsertGames(batch, metadata, nextPlyID); err != nil {
			log.Printf("failed to insert final batch from %s: %v", metadata.Filename, err)
			return err
		}
	}

	// delete the file if fully processed
	processed, err := database.FileProcessed(metadata.Id)
	if err != nil {
		return err
	}

	if processed {
		err = os.Remove(path)
		if err != nil {
			log.Fatalf("Unable to delete %s: %s", path, err)
		}
	}

	log.Printf("finished processing %s: %d games processed", metadata.Filename, gameIndex)
	return nil
}

// EnsureNGamesProcessed checks the local database to ensure that at least n games have been
// processed. If not, parse local file data until the threshold is reached.
func EnsureNGamesProcessed(n int) {
	processedGames, err := database.CountProcessed()
	if err != nil {
		log.Fatal(err)
	}

	remainingGames := n - processedGames

	// parse local PGN files until required number of games has been hit
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
