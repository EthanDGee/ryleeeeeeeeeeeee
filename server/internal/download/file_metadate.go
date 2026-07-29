package download

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func FetchFilesMetadata() []FileMetadata {
	gameCounts := GameCounts()

	if gameCounts == nil {
		return nil
	}

	var metadata []FileMetadata

	id := 0
	for filename, count := range gameCounts {
		fileUrl := LICHESS_BASE_URL + filename
		resp, err := http.Head(fileUrl)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		resp.Body.Close()

		metadata = append(metadata, FileMetadata{
			url:      fileUrl,
			filename: filename,
			games:    count,
			id:       id,
		})
		id++
	}

	return metadata
}

func GameCounts() map[string]int {
	response, err := http.Get(GAME_COUNTS_URL)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	defer response.Body.Close()

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
