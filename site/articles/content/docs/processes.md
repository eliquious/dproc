---
title: "Processes"

# Set type to 'docs' if you want to render page outside of configured section or if you render section other than 'docs'
type: docs

# Set page weight to re-arrange items in file-tree menu (if BookMenuBundle not set)
weight: 2

# (Optional) Set to mark page as flat section in file-tree menu (if BookMenuBundle not set)
bookFlatSection: true

# (Optional) Set to hide table of contents, overrides global value
bookShowToC: true
---

## **Processes**

Processes wrap a handler with a life-cycle and dispatch messages to the handler. Each processes runs in a separate go-routine to maximize concurrency. Each process has a unique ID, an internal `State`, an inbox and a `Handler`.

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

It should be noted that `SetState` is rarely used. However, it is useful in certain circumstances if the `Handler` needs to kill the `Process`.

`ProcessList` is a wrapper type for a slice of Processes. It allows for the quick dispatch of messages to a list of Processes. It is the primary method for sending messages to other processes and calls the `process.Send` method in `Dispatch`.

```go
// ProcessList is a list type for Process
type ProcessList []Process

// Dispatch dispatches a message to a list of processes
func (p ProcessList) Dispatch(m Message) {
	for _, proc := range p {
		proc.Send(m)
	}
}
```

### Life-cycle

Each process manages its own state which can be one of three possibilities:

```go
// State manages the processor state
type State string

// States
const (
	StateRunning State = "RUNNING"
	StateWaiting       = "WAITING"
	StateKilled        = "KILLED"
)
```

All processes start in the `WAITING` state until `engine.Start` is called and will run until in the `KILLED` state. The inner loop of the default process looks like this.

```go
switch msg.Type {
case MessageTypeStart:
	p.state = StateRunning
	p.Children().Dispatch(msg)
	p.handler.Handle(p.ctx, p, msg)
case MessageTypeStop:
	p.state = StateKilled
	p.handler.Handle(p.ctx, p, msg)
	p.Children().Dispatch(msg)
default:

	// Process message if running
	if p.state == StateRunning {
		p.handler.Handle(p.ctx, p, msg)
	}

	// Forward message if required
	if msg.Forward {
		p.Children().Dispatch(msg)
	}
}
```

`START` and `STOP` messages update the internal Process state, then handle and dispatch the messages. Otherwise, the handler processes the message if `RUNNING` and then the message is dispatched to child processes if required.

This is not the full implementation, however, it is the most important part. Additional logic exists for handling state updates outside the normal lifecycle as well as `context.Done` events.

### Creating processes

There is a helper method for creating a new default process. It is unsurprisingly called `NewDefaultProcess`. Using the simple handler that was created on the `Handler` page, creating the process looks like this:

```go
p := dproc.NewDefaultProcess(ctx, "Simple", &SimpleHandler{}, dproc.ProcessList{})
```

We give the process a `context.Context`, a name, a handler and an empty list of child processes. In a later example, we'll show how to build a pipeline.
