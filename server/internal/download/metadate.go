package download

import (
	"io"
	"log"
	"net/http"
	"server/rest/internal/config"
	"server/rest/internal/database"
	"server/rest/internal/models"
	"server/rest/internal/utils"
	"strconv"
	"strings"
)

func FetchFilesMetadata() {
	gameCounts := GameCounts()

	id := 0
	for filename, count := range gameCounts {
		fileUrl := config.LICHESS_BASE_URL + filename
		resp, err := http.Head(fileUrl)
		if err != nil {
			log.Fatal(err)
		}
		utils.Close(resp.Body, "response body")

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
}

func GameCounts() map[string]int {
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
