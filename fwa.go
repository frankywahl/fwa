package fwa

import (
	"context"
	"fmt"
	"time"

	faktory "github.com/contribsys/faktory/client"
	faktory_worker "github.com/contribsys/faktory_worker_go"
	"github.com/gobuffalo/buffalo/worker"
	"github.com/sirupsen/logrus"
)

// Adapter adapts faktory to use with buffalo.
type Adapter struct {
	Logger Logger
	mgr    *faktory_worker.Manager
	ctx    context.Context
}

// Option is a list configuration for the
// workers
type Option func(a *Adapter) error

// SetConcurrency is an option that will set the
// number of workers associated with a mgr
func SetConcurrency(value int) Option {
	return func(a *Adapter) error {
		a.mgr.Concurrency = value
		return nil
	}
}

// New constructs a new adapter.
func New(opts ...Option) (*Adapter, error) {
	adapter := &Adapter{
		Logger: logrus.New(),
		mgr:    faktory_worker.NewManager(),
		ctx:    context.Background(),
	}
	for _, opt := range opts {
		if err := opt(adapter); err != nil {
			return nil, fmt.Errorf("could not apply options: %w", err)
		}
	}
	return adapter, nil
}

// Start starts the adapter event loop.
func (q *Adapter) Start(ctx context.Context) error {
	q.Logger.Info("Starting the workers")
	q.ctx = ctx
	go func() {
		select {
		case <-ctx.Done():
			q.Stop()
		}
	}()
	q.mgr.Run()
	return nil
}

// Stop stops the adapter event loop.
func (q *Adapter) Stop() error {
	q.Logger.Info("Stopping faktory Worker")
	q.mgr.Terminate()
	return nil
}

// Register binds a new job, with a name and a handler.
func (q *Adapter) Register(name string, h worker.Handler) error {
	f := func(ctx faktory_worker.Context, args ...interface{}) error {
		if len(args) != 1 {
			return fmt.Errorf("error with arguments passing")
		}
		wargs := map[string]interface{}{}
		switch v := args[0].(type) {
		case nil:
		case map[string]interface{}:
			wargs = v
		default:
		}
		return h(wargs)
	}
	q.mgr.Register(name, f)
	return nil
}

// Perform sends a new job to the queue
func (q *Adapter) Perform(h worker.Job) error {
	return q.PerformAt(h, time.Now())
}

// PerformIn sends a new job to the queue, with a given delay.
func (q Adapter) PerformIn(job worker.Job, t time.Duration) error {
	return q.PerformAt(job, time.Now().Add(t))
}

// PerformAt sends a new job to the queue, with a given start time.
func (q *Adapter) PerformAt(h worker.Job, t time.Time) error {
	client, err := faktory.Open()
	if err != nil {
		return fmt.Errorf("could not open connection: %w", err)
	}
	defer client.Close()

	job := faktory.NewJob(h.Handler, h.Args)
	job.Queue = h.Queue
	job.Args = []interface{}{h.Args}
	job.At = t.Format(time.RFC3339)
	return client.Push(job)
}
