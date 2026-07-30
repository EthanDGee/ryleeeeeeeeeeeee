package utils

import (
	"io"
	"log"
)

func Close(c io.Closer, what string) {
	if err := c.Close(); err != nil {
		log.Printf("failed to close %s: %v", what, err)
	}
}
