package main

import (
	"fmt"
	"server/rest/internal/download"

	_ "turso.tech/database/tursogo"
)

func main() {
	files := download.FetchFilesMetadata()
	fmt.Println(files)
}
