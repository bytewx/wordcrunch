package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"wordcrunch/internal/concurrency"
	"wordcrunch/internal/models"
	"wordcrunch/internal/utils"
)

const chunkSize = 1024 * 1024

var ErrFileDoesNotExist = errors.New("fatal error: file does not exist")

func main() {
	filePtr := flag.String("file", "words.txt", "help: file with words")
	workersPtr := flag.Int("workers", 5, "help: number of workers")

	flag.Parse()

	if *workersPtr < 1 {
		log.Fatalf("fatal error: -workers must be >= 1, got %d", *workersPtr)
	}

	fmt.Printf("info: running word counter on file: %s, count of workers: %d\n", *filePtr, *workersPtr)

	file, err := os.Open(*filePtr)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatal(ErrFileDoesNotExist)
		}

		log.Fatalf("fatal error: failed to open file: %s: %v", *filePtr, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("warning: failed to close file: %s: %v", *filePtr, err)
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("fatal error: failed to get file statistics: %s: %v", *filePtr, err)
	}
	fileSize := fileInfo.Size()

	var totalWordCount atomic.Int64
	var chunks []models.Chunk
	var prevEnd int64

	for off := int64(0); off < fileSize; {
		end := min(off+chunkSize, fileSize)
		if end < fileSize {
			end = alignToRuneStart(file, end, fileSize)
			if end <= off {
				end = min(off+chunkSize, fileSize) // defence from endless loops from "junk" bytes
			}
		}
		chunks = append(chunks, models.Chunk{Offset: off, Length: end - off})
		off = end
		_ = prevEnd
	}

	err = concurrency.RunWorkerPool(context.Background(), chunks, *workersPtr, func(ctx context.Context, c models.Chunk) error {
		n, err := utils.ProcessChunk(file, c)
		if err != nil {
			return fmt.Errorf("failed to process chunk at offset %d: %w", c.Offset, err)
		}
		totalWordCount.Add(int64(n))
		return nil
	})
	if err != nil {
		log.Fatalf("fatal error: %v", err)
	}

	fmt.Printf("info: total word count: %d\n", totalWordCount.Load())
}

func alignToRuneStart(file *os.File, offset int64, fileSize int64) int64 {
	if offset >= fileSize {
		return offset
	}
	buf := make([]byte, 1)
	for offset > 0 {
		if _, err := file.ReadAt(buf, offset); err != nil {
			break
		}
		// continuation byte starts with 10xxxxxx
		if buf[0]&0xC0 != 0x80 {
			break
		}
		offset--
	}
	return offset
}
