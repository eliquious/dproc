---
title: "Messages"

# Set type to 'docs' if you want to render page outside of configured section or if you render section other than 'docs'
type: docs

# Set page weight to re-arrange items in file-tree menu (if BookMenuBundle not set)
weight: 2

# (Optional) Set to mark page as flat section in file-tree menu (if BookMenuBundle not set)
bookFlatSection: true

# (Optional) Set to hide table of contents, overrides global value
bookShowToC: true
---

## **Messages**

Messages are the primary mechanism for communicating with other processes. The struct looks like this:

```go
type Message struct {
    Timestamp time.Time
    Type      MessageType
    Forward   bool
    Value     interface{}
}
```

**Message fields**

- **Timestamp**: When the message was created (e.g. `time.Now().UTC()`)
- **Type**: MessageType is just a string type for quick type-switching when processing the message
- **Forward**: If `true`, the message will be sent to all child nodes recursively. Otherwise, it will only be sent to the process' direct children and no further.
- **Value**: Contains information fo be passed on. It will need to be type cast when received by the child processes.

Forwarding is neccessary for process life-cycle management. For instance, starting and stopping all processes. Most messages will not be forwarded. The `Timestamp` field is not strictly necessary.

*Note*: Setting the timestamp for every message will degrade performance.

### MessageTypes

All messages must have a type in order to be processed. MessageTypes are simply const strings which identify the message's purpose.

```go
// MessageType helps enumerate message types
type MessageType string
```

Custom message types can be defined like so:

```go
// TypeRandom is the message type for random numbers
const TypeRandom = dproc.MessageType("Random")
```

### Built-in messages

There are only two built-in messages types. They are also forwarded to all children when using the default process. This informs all processes in the pipeline about the life-cycle state. For instance, starting and stopping the pipeline.

```go
// Common message types
const (
	MessageTypeStart MessageType = "START"
	MessageTypeStop  MessageType = "STOP"
)
```

Handlers can also listen for these messages if additional initialization or clean-up is required, like opening/closing a file.
