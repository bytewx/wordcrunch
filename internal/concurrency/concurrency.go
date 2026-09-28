package concurrency

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// RunWorkerPool processes items concurrently with the given number of workers.
// It stops early and returns the first error encountered; ctx cancellation
// also stops remaining work.
func RunWorkerPool[T any](ctx context.Context, items []T, workers int, fn func(ctx context.Context, item T) error) error {
	if workers < 1 {
		workers = 1
	}
	if workers > len(items) {
		workers = max(len(items), 1)
	}

	g, ctx := errgroup.WithContext(ctx)
	jobs := make(chan T)

	g.Go(func() error {
		defer close(jobs)
		for _, item := range items {
			select {
			case jobs <- item:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	})

	for i := 0; i < workers; i++ {
		g.Go(func() error {
			for job := range jobs {
				if err := fn(ctx, job); err != nil {
					return err
				}
			}
			return nil
		})
	}

	return g.Wait()
}
