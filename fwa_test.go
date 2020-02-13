package fwa

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gobuffalo/buffalo/worker"
)

var q, _ = New()

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
			cancel()
			log.Fatal(ctx.Err())
		}
	}()

	go func() {
		fmt.Println("start q")
		err := q.Start(ctx)
		if err != nil {
			cancel()
			log.Fatal(err)
		}
	}()

	code := m.Run()

	err := q.Stop()
	if err != nil {
		log.Fatal(err)
	}
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
