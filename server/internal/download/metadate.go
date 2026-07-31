package download

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"server/rest/internal/config"
	"server/rest/internal/database"
	"server/rest/internal/models"
	"server/rest/internal/utils"
)

func FetchFilesMetadata() {
	gameCounts := GameCounts()
	log.Printf("found %d game files to check for metadata", len(gameCounts))

	id := 0
	for filename, count := range gameCounts {
		exists, err := database.MetadataExists(filename)
		if err != nil {
			log.Printf("failed to look up metadata for %s: %v", filename, err)
			id++
			continue
		}
		if exists {
			id++
			continue
		}

		log.Printf("fetching metadata for %s", filename)

		fileUrl := config.LICHESS_BASE_URL + filename
		resp, err := http.Head(fileUrl)
		if err != nil {
			log.Fatal(err)
		}
		utils.Close(resp.Body, "response body")
		time.Sleep(200 * time.Millisecond)

		metadata := models.Metadata{
			Url:      fileUrl,
			Filename: filename,
			Games:    count,
			Id:       id,
		}
		if err := database.UpsertMetadata(metadata); err != nil {
			log.Printf("failed to upsert metadata for %s: %v", filename, err)
		}

		id++
	}

	log.Println("metadata fetch complete")
}

func GameCounts() map[string]int {
	log.Println("fetching game counts")
	response, err := http.Get(config.GAME_COUNTS_URL)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	defer utils.Close(response.Body, "response body")

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	gameCounts := make(map[string]int)

	for _, line := range strings.Split(string(body), "\n") {
		split := strings.Split(line, " ")
		if len(split) < 2 {
			continue
		}

		numGames, err := strconv.Atoi(split[1])
		if err != nil {
			log.Fatal(err)
			return nil
		}

		gameCounts[split[0]] = numGames
	}
	return gameCounts
}
