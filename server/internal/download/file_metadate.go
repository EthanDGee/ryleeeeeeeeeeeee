package download

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func FetchFilesMetadata() []FileMetadata {
	// gameCounts := GameCounts()

	return nil
}

func GameCounts() map[string]uint {
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

	gameCounts := make(map[string]uint)

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

		gameCounts[split[0]] = uint(numGames)
	}
	return gameCounts
}
