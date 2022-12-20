package fwa

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gobuffalo/buffalo/worker"
)

var q = must(New())

func must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}

func runTests(m *testing.M) int {
	stop, err := setupFaktory()
	if err != nil {
		panic(err)
	}
	defer stop()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	go func() {
		err := q.Start(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	code := m.Run()

	if err := q.Stop(); err != nil {
		log.Fatal(err)
	}
	return code
}
func TestMain(m *testing.M) {
	code := runTests(m)
	os.Exit(code)
}

func Test_Perform(t *testing.T) {
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
