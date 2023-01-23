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
	logger     Logger
	mgr        *faktory_worker.Manager
	pool       *faktory.Pool
	cancelFunc context.CancelFunc
}

// New constructs a new adapter.
func New(opts ...Option) (*Adapter, error) {
	config := newConfig()
	for _, opt := range opts {
		if err := opt.apply(config); err != nil {
			return nil, fmt.Errorf("could not apply options: %w", err)
		}
	}
	pool, err := faktory.NewPool(config.concurrency)
	if err != nil {
		return nil, fmt.Errorf("could not create pool: %w", err)
	}
	mgr := faktory_worker.NewManager()
	mgr.Use(config.middlewares...)
	mgr.ProcessWeightedPriorityQueues(config.queues)
	adapter := &Adapter{
		logger: config.logger,
		mgr:    mgr,
		pool:   pool,
	}
	return adapter, nil
}

// Start starts the adapter event loop.
func (q *Adapter) Start(ctx context.Context) error {
	q.logger.Infof("Starting the faktory workers\n")
	ctx, cancel := context.WithCancel(ctx)
	q.cancelFunc = cancel

	if err := q.mgr.RunWithContext(ctx); err != nil {
		return err
	}
	return nil
}

// Stop stops the adapter event loop.
func (q *Adapter) Stop() error {
	q.logger.Infof("Stopping faktory Worker\n")
	q.cancelFunc()
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

// Perform sends a new job to the queue.
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
