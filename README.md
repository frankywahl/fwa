# Faktory/work Adapter for Buffalo

This package implements the `github.com/gobuffalo/buffalo/worker.Worker`

## Setup

```go

// ...

workerAdapter, _ := fwa.New()

buffalo.New(buffalo.Options{
  // ...
  Worker: workerAdatper,
  // ...
})
```
