package main

import (
	"log"

	"server/rest/internal/config"
	"server/rest/internal/database"
	"server/rest/internal/download"

	_ "turso.tech/database/tursogo"
)

func main() {
	if err := database.InitializeDatabase(); err != nil {
		log.Fatal(err)
	}

	download.FetchFilesMetadata()

	download.EnsureNGamesDownloaded(config.GAME_COUNT)
}
