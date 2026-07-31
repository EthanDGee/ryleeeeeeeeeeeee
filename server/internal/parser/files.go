package parser

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"server/rest/internal/config"
	"server/rest/internal/models"
	"server/rest/internal/utils"

	"github.com/klauspost/compress/zstd"
)

func DecompressPGN(metadata models.Metadata) {
	compressedPath := filepath.Join(config.FILES_FOLDER, metadata.Filename)
	decompressedPath := strings.TrimSuffix(compressedPath, ".zst")

	in, err := os.Open(compressedPath)
	if err != nil {
		log.Fatalf("failed to open %s: %v", compressedPath, err)
	}
	defer utils.Close(in, "compressed file")

	decoder, err := zstd.NewReader(in)
	if err != nil {
		log.Fatalf("failed to create zstd reader for %s: %v", compressedPath, err)
	}
	defer decoder.Close()

	out, err := os.Create(decompressedPath)
	if err != nil {
		log.Fatalf("failed to create %s: %v", decompressedPath, err)
	}
	defer utils.Close(out, "decompressed file")

	if _, err := io.Copy(out, decoder); err != nil {
		log.Fatalf("failed to decompress %s: %v", compressedPath, err)
	}
}
