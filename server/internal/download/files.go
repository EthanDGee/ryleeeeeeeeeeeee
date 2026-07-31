package download

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"server/rest/internal/config"
	"server/rest/internal/database"
	"server/rest/internal/models"
	"server/rest/internal/utils"
)

func FetchGameFile(metadata models.Metadata) {
	if err := os.MkdirAll(config.FILES_FOLDER, 0o755); err != nil {
		log.Fatal(err)
	}

	log.Printf("downloading %s", metadata.Filename)

	resp, err := http.Get(metadata.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer utils.Close(resp.Body, "Game File Response Body")

	if resp.StatusCode != http.StatusOK {
		log.Fatal(fmt.Errorf("unexpected status downloading %s: %s", metadata.Url, resp.Status))
	}

	path := filepath.Join(config.FILES_FOLDER, metadata.Filename)

	out, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer utils.Close(out, "Downloaded file")

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	success, err := database.MarkDownloaded(metadata.Id)
	if err != nil {
		log.Fatal(err)
	} else if !success {
		log.Fatalf("Failed to mark %s as downloaded", metadata.Filename)
	}

	log.Printf("downloaded %s", metadata.Filename)
}

func EnsureNGamesDownloaded(n int) {
	downloadedCount, err := database.CountDownloaded()
	if err != nil {
		log.Fatal(err)
	}

	for downloadedCount < n {
		nextFile, err := database.GetSmallestUndownloadedFile()
		if err != nil {
			log.Fatal(err)
		}

		FetchGameFile(nextFile)
		downloadedCount, err = database.CountDownloaded()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("downloaded games: %d/%d (%.1f%%)", downloadedCount, n, float64(downloadedCount)/float64(n)*100)
	}

	log.Printf("finished ensuring %d games downloaded\n", n)
}
