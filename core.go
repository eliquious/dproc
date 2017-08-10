package dproc

import (
	"context"
	"strings"
	"time"
)

// MessageType is the message type
type MessageType string

// Common message types
const (
	MessageTypeStart MessageType = "START"
	MessageTypeStop  MessageType = "STOP"
)

// Message is the value type exchanged between data processors
type Message struct {
	Timestamp time.Time
	Type      MessageType
	Forward   bool
	Values    map[string]string
}

// State manages the processor state
type State int

// Processor processes messages.
type Processor interface {
	Name() string
	State() State
	SetContext(context.Context)
	Process(Message)
}

// Dispatcher dispatches messages to child processors.
type Dispatcher interface {
	AddProcessor(Processor)
	Dispatch(Message)
}

// Engine manages the pipeline
type Engine interface {
	Dispatcher
	Start()
	Stop()
}

// Service allows for global processors.
type Service interface {
	Name() string
	Process(Message)
}

type contextKey string

var serviceKey = contextKey("svc")

// SendTo allows for sending messages to services
func SendTo(ctx context.Context, svc string, msg Message) {
	if v := ctx.Value(serviceKey); v != nil {
		if svcs, ok := v.([]Service); ok {
			for i := 0; i < len(svcs); i++ {
				if strings.Equal(svcs[i].Name(), svc) {
					svcs[i].Process(msg)
				}
			}
		}
		return
	}
}

// WithService adds a service to a context.Context.
func WithService(ctx context.Context, svc Service) context.Context {
	if v := ctx.Value(serviceKey); v != nil {
		if svcs, ok := v.([]Service); ok {
			svcs = append(svcs, svc)
			return context.WithValue(ctx, serviceKey, svcs)
		}
		return ctx
	}
	return context.WithValue(ctx, serviceKey, []Service{svc})
}

// NewEngine creates a new engine
func NewEngine(ctx context.Context) Engine {
	return &engine{&dispatcher{make([]Processor, 0)}}
}

type engine struct {
	dispatcher Dispatcher
}

func (e *engine) Start() {
	e.dispatcher.Dispatch(Message{Timestamp: time.Now(), MessageType: MessageTypeStart, Forward: true})
}

func (e *engine) Stop() {
	e.dispatcher.Dispatch(Message{Timestamp: time.Now(), MessageType: MessageTypeStop, Forward: true})
}

type dispatcher struct {
	children []Processor
}

func (d *dispatcher) Dispatch(msg Message) {
	for i := 0; i < len(d.children); i++ {
		d.children[i].Process(msg)
	}
}

func (d *dispatcher) AddProcessor(proc Processor) {
	d.children = append(d.children, proc)
}
