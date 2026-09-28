# wordcrunch

A Go CLI utility that counts words in a file using concurrent processing.

The file is split into fixed-size chunks (1 MB), and each chunk is processed
by a separate worker in a goroutine pool. Word boundaries at chunk edges are
handled correctly (see `internal/utils.ProcessChunk`).

## Usage

```bash
go run ./cmd/wordcrunch -file words.txt -workers 8
```

Flags:

- `-file` — path to the file with words (default: `words.txt`)
- `-workers` — number of workers in the pool, must be >= 1 (default: `5`)

## Architecture

- `cmd/wordcrunch` — entry point: flag parsing, opening the file, splitting it
  into chunks (`internal/models.Chunk`), running the worker pool, and
  aggregating the result.
- `internal/concurrency` — generic worker pool, `RunWorkerPool[T any]`, built
  on `errgroup` with context cancellation: an error from any worker stops
  the remaining work and is propagated to the caller.
- `internal/utils` — counts words in a chunk (`ProcessChunk`), correctly
  handling the boundary between adjacent chunks.
- `internal/models` — shared data structures (`Chunk`).

## Tests

```bash
go test -race ./...
```

CI runs `go vet` and `go test -race` on every push and pull request
(see `.github/workflows/ci.yml`).

## License

MIT — see [LICENSE](LICENSE).