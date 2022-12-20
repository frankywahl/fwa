package fwa

import (
	"context"
	"fmt"
	"time"

	faktory "github.com/contribsys/faktory/client"
	faktory_worker "github.com/contribsys/faktory_worker_go"
	"github.com/gobuffalo/buffalo/worker"
)

// Adapter adapts faktory to use with buffalo.
type Adapter struct {
	Logger Logger
	mgr    *faktory_worker.Manager
	ctx    context.Context
	pool   *faktory.Pool
}

// Option is a list configuration for the
// workers
type Option func(a *Adapter) error

// Queues to read from. map of queue name to queue priority
func WithQueues(queues map[string]int) Option {
	return func(a *Adapter) error {
		for _, v := range queues {
			if v <= 0 {
				return fmt.Errorf("queue priority with %d is invalid", v)
			}
		}
		a.mgr.ProcessWeightedPriorityQueues(queues)
		return nil
	}
}

// SetConcurrency is an option that will set the
// number of workers associated with a mgr
func SetConcurrency(value int) Option {
	return func(a *Adapter) error {
		a.mgr.Concurrency = value
		return nil
	}
}

// SetPool Defines a pool to get from
// the client from
func SetPool(p *faktory.Pool) Option {
	return func(a *Adapter) error {
		a.pool = p
		return nil
	}

}

func WithLogger(l Logger) Option {
	return func(a *Adapter) error {
		a.Logger = l
		return nil
	}
}

// New constructs a new adapter.
func New(opts ...Option) (*Adapter, error) {
	pool, err := faktory.NewPool(20)
	if err != nil {
		return nil, fmt.Errorf("could not create pool: %w", err)
	}
	mgr := faktory_worker.NewManager()
	mgr.ProcessWeightedPriorityQueues(map[string]int{"default": 1})
	adapter := &Adapter{
		Logger: &noopLogger{},
		mgr:    mgr,
		ctx:    context.Background(),
		pool:   pool,
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
	q.mgr.RunWithContext(ctx)
	return nil
}

// Stop stops the adapter event loop.
func (q *Adapter) Stop() error {
	q.Logger.Info("Stopping faktory Worker")
	q.mgr.Terminate(false)
	return nil
}

// Register binds a new job, with a name and a handler.
func (q *Adapter) Register(name string, h worker.Handler) error {
	f := func(ctx context.Context, args ...interface{}) error {
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
	return q.pool.With(
		func(c *faktory.Client) error {
			job := faktory.NewJob(h.Handler, h.Args)
			job.Queue = h.Queue
			job.Args = []interface{}{h.Args}
			job.At = t.Format(time.RFC3339)
			return c.Push(job)
		},
	)
}
