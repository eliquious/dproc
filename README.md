<p align="center">
  <img src="https://eliquious.github.io/dproc/images/dProc_small.svg" style="margin-bottom:24px">
</p>
<p align="center">
    <a href="https://godoc.org/github.com/eliquious/dproc">
      <img src="https://img.shields.io/badge/godoc-reference-blue.svg" alt="godoc">
    </a>
    <a href="https://eliquious.github.io/dproc">
      <img src="https://img.shields.io/badge/support-articles-red.svg" alt="Docs">
    </a>
</p>

<br/>

`dProc` is a small, generic data processing library written in Go. It is modeled after [Actors](https://en.wikipedia.org/wiki/Actor_model) and a bit of [Flow-based programming](https://en.wikipedia.org/wiki/Flow-based_programming). It is a small library (< 250 lines of code) with very few interfaces.

## Usage

`dProc` constructs data pipelines for execution with user-defined data handlers. Each handler processes or generates messages and sends them to other processes until the pipeline is complete. Messages are simple structs which encapsulate the `MessageType` and value. The pipeline is composes of concurrent processes and each has a handler for adding behavior.

### Handlers

Handlers are composed of single functions which process each message based on its type. Here's the interface.

```go
type Handler interface {
    Handle(ctx context.Context, proc Process, msg Message)
}
```

### Processes

A `Process` wraps a `Handler` with a life-cycle which allows the `Engine` to manage the entire pipeline from start to finish. Here's the interface.

```go
// Process processes incoming messages and sends messages to other Processes.
type Process interface {

	// Name is used primarily for logging purposes
	Name() string

	// SetState can be used to update the Process' internal state outside of the normal lifecycle.
	SetState(State)

	// Start is called by the Engine when the pipeline starts and initializes the go-routine and lifecycle.
	Start(*sync.WaitGroup)

	// Send adds another message to the inbox for the Process.
	Send(Message)

	// Children returns a list of child processes.
	Children() ProcessList
}
```

#### Creating a new process

Creating a new Process is simple and there's a built-in helper function. Unsurprisingly, it is called `NewDefaultProcess`. 

```go
p := dproc.NewDefaultProcess(ctx, "Process Name", &Handler{}, dproc.ProcessList{})
```

Each `Process` is given a `context.Context`, a name, a `Handler` and a `ProcessList` which is a list of child processes.

### Engine

The engine manages the data pipeline for all the processes and is created using another helper function.

```go
// Engine executes the pipeline and runs until no processes are running.
type Engine interface {
	Start(*sync.WaitGroup)
	Stop()
}

// NewEngine creates an Engine with the given context.Context and a list of processes.
func NewEngine(ctx context.Context, ps ProcessList) Engine
```

A `sync.WaitGroup` is used to manage the number of Processes and is used to wait until everything is finished.

### Documentation

There more indepth articles about each interface [here](https://eliquious.github.io/dproc).
