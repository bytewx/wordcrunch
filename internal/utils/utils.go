package utils

import (
	"errors"
	"io"
	"os"
	"wordcrunch/internal/models"
)

func ProcessChunk(file *os.File, c models.Chunk) (int, error) {
	start := c.Offset
	size := c.Length
	prevIsSpace := true

	if c.Offset > 0 {
		start--
		size++
	}

	buf := make([]byte, size)

	n, err := file.ReadAt(buf, start)
	if err != nil && !errors.Is(err, io.EOF) {
		return 0, err
	}
	buf = buf[:n]

	if c.Offset > 0 {
		if len(buf) == 0 {
			return 0, nil
		}
		prevIsSpace = isSpace(buf[0])
		buf = buf[1:]
	}

	return countWordStarts(buf, prevIsSpace), nil
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\t' || b == '\r'
}

func countWordStarts(buf []byte, prevIsSpace bool) int {
	count := 0
	for _, b := range buf {
		space := isSpace(b)
		if !space && prevIsSpace {
			count++
		}
		prevIsSpace = space
	}
	return count
}
