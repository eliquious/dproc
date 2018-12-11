---
title: "Engine"

# Set type to 'docs' if you want to render page outside of configured section or if you render section other than 'docs'
type: docs

# Set page weight to re-arrange items in file-tree menu (if BookMenuBundle not set)
weight: 2

# (Optional) Set to mark page as flat section in file-tree menu (if BookMenuBundle not set)
bookFlatSection: true

# (Optional) Set to hide table of contents, overrides global value
bookShowToC: true
---

## **Engine**

The `dproc.Engine` is similar to a pipeline. It allows for many processes to be created and organized before the entire pipeline is executed. The `Engine` can be managed with a `context.Context`, `context.CancelFunc` and a `sync.WaitGroup` allowing for the pipeline to be killed.

Here's the interface:

```go
// Engine manages the pipeline
type Engine interface {
	Start(*sync.WaitGroup)
	Stop()
}
```

### Starting the engine

Starting the `Engine` is simple. It requires calling `engine.Start(sync.WaitGroup)`. However, there a helper method for creating the `Engine` which is surprisingly called `NewEngine`.

Its signature is:

```go
func NewEngine(ctx context.Context, ps ProcessList) Engine
```

If the context is cancelled, it will be caught and will stop everything. However, the built-in Engine has a built-in `context.CancelFunc` that is called with `engine.Stop`. 

### Example: Creating and starting an Engine

This is the `main` from the included example. It is somewhat contrived. One process sends random numbers and the child process prints a count every second of how many messages were received. Normally between 2-3M messages are sent per second through this pipeline.

```
func main() {
	log.SetOutput(os.Stdout)
	ctx := context.Background()

	var wg sync.WaitGroup
	engine := dproc.NewEngine(ctx, dproc.ProcessList{
		dproc.NewDefaultProcess(ctx, "Random Numbers", &RandomGenerator{time.Second * 5}, dproc.ProcessList{
			dproc.NewDefaultProcess(ctx, "Random Logger", &RandomLogger{Ticker: time.NewTicker(time.Second)}, dproc.ProcessList{}),
		}),
	})
	start := time.Now()
	engine.Start(&wg)

	wg.Wait()
	fmt.Println("Elapsed: ", time.Since(start))
	fmt.Println("\nExiting...")
}
```
