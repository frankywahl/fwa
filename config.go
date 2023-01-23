package fwa

import (
	faktory_worker "github.com/contribsys/faktory_worker_go"
)

type configOptions struct {
	queues      map[string]int
	middlewares []faktory_worker.MiddlewareFunc
	concurrency int
	logger      Logger
}

func newConfig() *configOptions {
	return &configOptions{
		queues:      map[string]int{"default": 1},
		middlewares: []faktory_worker.MiddlewareFunc{},
		concurrency: 20,
		logger:      &noopLogger{},
	}
}

// Option is a list configuration for the workers.
type Option interface {
	apply(c *configOptions) error
}

type withQueueOption struct {
	queues map[string]int
}

func (c *withQueueOption) apply(config *configOptions) error {
	config.queues = c.queues
	return nil
}

// WithQueues sets the queues read from. Map of queue name to queue priority.
func WithQueues(queues map[string]int) Option {
	return &withQueueOption{
		queues: queues,
	}
}

type middlewareOption struct {
	middlewares []faktory_worker.MiddlewareFunc
}

func (c *middlewareOption) apply(config *configOptions) error {
	config.middlewares = c.middlewares
	return nil
}

// WithMiddleware allows us to register faktory middleware with the functions
//
// This is to be used with faktory middleware.
func WithMiddleware(m ...faktory_worker.MiddlewareFunc) Option {
	return &middlewareOption{
		middlewares: m,
	}
}

type concurrencyOption struct {
	concurrency int
}

func (c *concurrencyOption) apply(config *configOptions) error {
	config.concurrency = c.concurrency
	return nil
}

// SetConcurrency is an option that will set the
// number of workers associated with a manager.
func SetConcurrency(value int) Option {
	return &concurrencyOption{
		concurrency: value,
	}
}

type loggerOption struct {
	logger Logger
}

func (c *loggerOption) apply(config *configOptions) error {
	config.logger = c.logger
	return nil
}

// WithLogger defines a logger that the adapter will use.
func WithLogger(l Logger) Option {
	return &loggerOption{
		logger: l,
	}
}
