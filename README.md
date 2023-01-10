# Faktory/Buffalo Worker Adapter

This package implements the `github.com/gobuffalo/buffalo/worker.Worker`.

It allows the usage of [Faktory](https://contribsys.com/faktory/) as a background worker in the [Buffalo](https://gobuffalo.io/) framework.

## Setup

```go
import(
	"github.com/frankywahl/fwa"
	// ...
)

// ...

worker, err := fwa.New(
	fwa.WithQueues(map[string]int{
		"default":	   1,
		"other_queue": 1,
	}),
)
// handle the error

buffalo.New(buffalo.Options{
	// ...
	Worker: worker,
	// ...
})
```

The rest of the setup for running buffalo and faktory can be found in [buffalo's documentation](https://gobuffalo.io/documentation/guides/workers/)

## Tested with

| App             | Version                                                            |
| ---             | ---                                                                |
| Buffalo-CLI     | [0.18.2](https://github.com/gobuffalo/buffalo/releases/tag/v1.0.1) |
| Buffalo-Package | [1.0.1]()                                                          |
| Faktory         | 1.6.2                                                              |
| FaktoryWorkerGo | 1.6.0                                                              |
