package main

import (
	"log"

	"server/rest/internal/config"
	"server/rest/internal/database"
	"server/rest/internal/download"
	"server/rest/internal/parser"

	_ "turso.tech/database/tursogo"
)

func main() {
	if err := database.InitializeDatabase(); err != nil {
		log.Fatal(err)
	}

	download.FetchFilesMetadata()

	log.Printf("ensuring at least %d games are downloaded", config.GAME_COUNT)
	download.EnsureNGamesDownloaded(config.GAME_COUNT)

	log.Printf("ensuring at least %d games processed", config.GAME_COUNT)
	parser.EnsureNGamesProcessed(config.GAME_COUNT)

	log.Println("done")
}
