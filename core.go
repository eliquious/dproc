package dproc

import (
	"context"
	"sync"
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
	Values    map[string]interface{}
}

// State manages the processor state
type State int

// States
const (
	StateRunning State = iota
	StateWaiting
	StateKilled
)

// Processor processes messages.
type Processor interface {
	Name() string
	// State() State
	Process(Message)
}

// ProcessorList is a list type for Processors
type ProcessorList []Processor

// Dispatch dispatches a message to a list of processes
func (p ProcessorList) Dispatch(m Message) {
	for i := 0; i < len(p); i++ {
		p[i].Process(m)
	}
}

// Engine manages the pipeline
type Engine interface {
	Start()
	Stop()
}

// Service allows for global processors.
type Service interface {
	Name() string
	Process(Message)
}

// ServiceList is a value type for a list of services
type ServiceList []Service

// SendTo sends a message to a particular service
func (s ServiceList) SendTo(name string, m Message) {
	for i := 0; i < len(s); i++ {
		if s[i].Name() == name {
			s[i].Process(m)
		}
	}
}

type contextKey string

var serviceKey = contextKey("svc")
var nameKey = contextKey("name")
var waitGroupKey = contextKey("waitgroup")

// SendTo allows for sending messages to services
func SendTo(ctx context.Context, svc string, msg Message) {
	if v := ctx.Value(serviceKey); v != nil {
		if svcs, ok := v.(ServiceList); ok {
			svcs.SendTo(svc, msg)
		}
		return
	}
}

// WithService adds a service to a context.Context.
func WithService(ctx context.Context, svc Service) context.Context {
	if v := ctx.Value(serviceKey); v != nil {
		if svcs, ok := v.(ServiceList); ok {
			svcs = append(svcs, svc)
			return context.WithValue(ctx, serviceKey, svcs)
		}
		return ctx
	}
	return context.WithValue(ctx, serviceKey, ServiceList{svc})
}

// WithName adds a name to a context.Context.
func WithName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, nameKey, name)
}

// Name gets a name from a context.Context.
func Name(ctx context.Context) string {
	if v := ctx.Value(nameKey); v != nil {
		if name, ok := v.(string); ok {
			return name
		}
		return ""
	}
	return ""
}

// WithWaitGroup adds a sync.WaitGroup to a context.Context.
func WithWaitGroup(ctx context.Context, wg *sync.WaitGroup) context.Context {
	return context.WithValue(ctx, waitGroupKey, wg)
}

// Done decrements a sync.WaitGroup from a context.Context.
func Done(ctx context.Context) {
	if v := ctx.Value(waitGroupKey); v != nil {
		if wg, ok := v.(*sync.WaitGroup); ok {
			wg.Done()
		}
	}
}

// Add increments a sync.WaitGroup from a context.Context.
func Add(ctx context.Context) {
	if v := ctx.Value(waitGroupKey); v != nil {
		if wg, ok := v.(*sync.WaitGroup); ok {
			wg.Add(1)
		}
	}
}

// NewEngine creates a new engine
func NewEngine(ctx context.Context, cancel context.CancelFunc, ps ProcessorList) Engine {
	return &engine{ctx, cancel, ps}
}

type engine struct {
	ctx      context.Context
	cancel   context.CancelFunc
	children ProcessorList
}

func (e *engine) Start() {
	e.children.Dispatch(Message{Timestamp: time.Now(), Type: MessageTypeStart, Forward: true})
}

func (e *engine) Stop() {
	e.children.Dispatch(Message{Timestamp: time.Now(), Type: MessageTypeStop, Forward: true})
	e.cancel()
}
