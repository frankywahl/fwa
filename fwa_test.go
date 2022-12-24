package fwa

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	faktory "github.com/contribsys/faktory/client"
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
	now := time.Now().UnixNano()
	handlerName := fmt.Sprintf("perform_%d", now)
	if err := q.Register(handlerName, func(worker.Args) error {
		hit = true
		wg.Done()
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := q.Perform(worker.Job{
		Handler: handlerName,
	}); err != nil {
		t.Fatal(err)
	}
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
	now := time.Now().UnixNano()
	handlerName := fmt.Sprintf("performAt_%d", now)
	if err := q.Register(handlerName, func(args worker.Args) error {
		hit = true
		wg.Done()
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := q.PerformAt(worker.Job{
		Handler: handlerName,
	}, time.Now().Add(5*time.Nanosecond)); err != nil {
		t.Fatal(err)
	}
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
	now := time.Now().UnixNano()
	handlerName := fmt.Sprintf("performIn_%d", now)
	if err := q.Register(handlerName, func(worker.Args) error {
		hit = true
		wg.Done()
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := q.PerformIn(worker.Job{
		Handler: handlerName,
	}, 5*time.Nanosecond); err != nil {
		t.Fatal(err)
	}
	wg.Wait()
	if !hit {
		t.Errorf("should be true")
	}
}

func Test_Middleware(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	middleCount := 0
	middleFunc := func(ctx context.Context, job *faktory.Job, next func(ctx context.Context) error) error {
		middleCount++
		return next(ctx)
	}
	q, close := startAdapter(t, context.Background(), WithMiddleware(middleFunc))
	defer close()
	now := time.Now().UnixNano()
	handlerName := fmt.Sprintf("perform_%d", now)
	if err := q.Register(handlerName, func(worker.Args) error {
		defer wg.Done()
		// Does nothing
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := q.PerformIn(worker.Job{
		Handler: handlerName,
	}, 1*time.Second); err != nil {
		t.Fatal(err)
	}
	wg.Wait()
	if middleCount != 1 {
		t.Errorf("expected to hit the middleware 1, hit it %d times", middleCount)
	}
}
