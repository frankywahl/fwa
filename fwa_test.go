package fwa

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gobuffalo/buffalo/worker"
)

func must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}

func startAdapter(t *testing.T, ctx context.Context, options ...Option) (worker.Worker, func() error) {
	q := must(New(options...))
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() error {
		defer wg.Done()
		if err := q.Start(ctx); err != nil {
			return err
		}
		return nil
	}()
	close := func() error {
		cancel()
		wg.Wait()
		return nil
	}
	return q, close
}

func Test_Perform(t *testing.T) {
	q, close := startAdapter(t, context.Background())
	defer close()
	hit := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	q.Register("perform", func(worker.Args) error {
		hit = true
		wg.Done()
		return nil
	})
	q.Perform(worker.Job{
		Handler: "perform",
	})
	wg.Wait()
	if !hit {
		t.Errorf("should be true")
	}
}

func Test_PerformAt(t *testing.T) {
	q, close := startAdapter(t, context.Background())
	defer close()
	hit := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	q.Register("perform_at", func(args worker.Args) error {
		hit = true
		wg.Done()
		return nil
	})
	q.PerformAt(worker.Job{
		Handler: "perform_at",
	}, time.Now().Add(5*time.Nanosecond))
	wg.Wait()
	if !hit {
		t.Errorf("should be true")
	}
}

func Test_PerformIn(t *testing.T) {
	q, close := startAdapter(t, context.Background())
	defer close()
	hit := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	q.Register("perform_in", func(worker.Args) error {
		hit = true
		wg.Done()
		return nil
	})
	q.PerformIn(worker.Job{
		Handler: "perform_in",
	}, 5*time.Nanosecond)
	wg.Wait()
	if !hit {
		t.Errorf("should be true")
	}
}
