---
title: "Handlers"

# Set type to 'docs' if you want to render page outside of configured section or if you render section other than 'docs'
type: docs

# Set page weight to re-arrange items in file-tree menu (if BookMenuBundle not set)
weight: 2

# (Optional) Set to mark page as flat section in file-tree menu (if BookMenuBundle not set)
bookFlatSection: true

# (Optional) Set to hide table of contents, overrides global value
bookShowToC: true
---

## **Handlers**

Handlers handle incoming messages as well as send messages to child processes. The interface only has one method.

```go
type Handler interface {
    Handle(ctx context.Context, proc Process, msg Message)
}
```

The `Engine` `context.Context` is passed in mainly for long running processes. However, `ctx.Done()` is monitored between function calls regardless. The `Process` is included as the `Handler` can modify the process `State` and is needed to dispatch messages.

### **Implementing Handlers**

Implementing handlers is a relatively straightforward task. Afterall, only one method needs to be implemented. However, in reality, a `Message` and a `MessageType` are needed as well for it to be useful.

```go
// TypeFloat64
const TypeFloat64 = dproc.MessageType("Float64")

// SimpleHandler handles dproc.MessageTypeStart, dproc.MessageTypeStop and TypeFloat64.
type SimpleHandler struct {
}

// Handle handles each individual messages. The message value for TypeFloat64 is a float64.
func (h *SimpleHandler) Handle(ctx context.Context, proc dproc.Process, msg dproc.Message) {
	switch msg.Type {
	default:

		// Log any unknown message types.
		log.Printf("[%s] - Unknown message type: %s", proc.Name(), msg.Type)

	case dproc.MessageTypeStart:

		// Log a starting message. Generally, initialization 
		// code executes here.
		log.Printf("[%s] - Starting...", proc.Name())

	case dproc.MessageTypeStop:

		// Log the last message for the handler. The default process will 
		// not call the handler for any message after this.
		log.Printf("[%s] - Exiting...", proc.Name())

	case TypeFloat64:

		// Type cast and log the message value
		value := msg.Value.(float64)
		log.Printf("[%s] - Message: %f", proc.Name(), value)

		// Dispatch sends the square of the value to any child processes.
		proc.Children().Dispatch(dproc.Message{
			Type:    TypeFloat64,
			Value: value * value,
		})
	}
}
```

This example is somewhat contrived as a `Process` will most likely not be used simply to square a value. However, it outlines how to handle message types and send messages. 
