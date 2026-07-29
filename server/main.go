package main

import (
	"fmt"
	"server/rest/internal/download"
)

func main() {
	files := download.FetchFilesMetadata()
	fmt.Println(files)
}
