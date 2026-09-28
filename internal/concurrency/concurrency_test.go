package concurrency

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
)

func TestRunWorkerPool(t *testing.T) {
	items := make([]int, 100)
	for i := range items {
		items[i] = i
	}

	for _, workers := range []int{0, 1, 3, 50} {
		var processed atomic.Int64
		err := RunWorkerPool(context.Background(), items, workers, func(ctx context.Context, item int) error {
			processed.Add(1)
			return nil
		})
		if err != nil {
			t.Fatalf("workers=%d: unexpected error: %v", workers, err)
		}
		if got := processed.Load(); got != int64(len(items)) {
			t.Errorf("workers=%d: processed %d items, want %d", workers, got, len(items))
		}
	}
}

func TestRunWorkerPoolEmptyItems(t *testing.T) {
	var processed atomic.Int64
	err := RunWorkerPool(context.Background(), []int{}, 5, func(ctx context.Context, item int) error {
		processed.Add(1)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := processed.Load(); got != 0 {
		t.Errorf("processed %d items, want 0", got)
	}
}

func TestRunWorkerPoolPropagatesError(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	wantErr := errors.New("boom")

	err := RunWorkerPool(context.Background(), items, 2, func(ctx context.Context, item int) error {
		if item == 3 {
			return wantErr
		}
		return nil
	})
	if !errors.Is(err, wantErr) {
		t.Errorf("got error %v, want %v", err, wantErr)
	}
}
