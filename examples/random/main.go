package main

import (
	"context"
	"fmt"
	"github.com/eliquious/dproc"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx := context.Background()
	engine := dproc.NewEngine(ctx, dproc.ProcessorList{
		NewRandomProc(ctx, "RandomProc", dproc.ProcessorList{
			NewRandomLogger(ctx, "Logger"),
		}),
	})
	engine.Start()
	// time.Sleep(time.Second * 5)

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	<-c
	fmt.Println("\nExiting...")
	engine.Stop()
}

// NewRandomProc creates a new dproc.Processor that emits random numbers
func NewRandomProc(ctx context.Context, name string, ps dproc.ProcessorList) dproc.Processor {
	proc := &randomProc{ctx, name, ps, dproc.StateWaiting, make(chan dproc.Message, 1)}
	go proc.start()
	return proc
}

// TypeRandom is the message type for random numbers
const TypeRandom = dproc.MessageType("Random")

type randomProc struct {
	ctx      context.Context
	name     string
	children dproc.ProcessorList
	state    dproc.State
	inbox    chan dproc.Message
}

func (r *randomProc) Name() string {
	return r.name
}

func (r *randomProc) Process(msg dproc.Message) {
	r.inbox <- msg
}

func (r *randomProc) start() {
	for {
		select {
		case <-time.After(time.Millisecond):
			if r.state == dproc.StateRunning {
				r.children.Dispatch(dproc.Message{
					Timestamp: time.Now(),
					Type:      TypeRandom,
					Forward:   false,
					Values: map[string]interface{}{
						"random": rand.Float64(),
					},
				})
			}
		case msg := <-r.inbox:
			// fmt.Printf("%s - [%s] %s %+v\n", time.Now().UTC().Format(time.RFC3339), r.name, msg.Type, msg.Values)

			switch msg.Type {
			case dproc.MessageTypeStart:
				r.state = dproc.StateRunning
			case dproc.MessageTypeStop:
				r.state = dproc.StateKilled
			}
			if msg.Forward {
				r.children.Dispatch(msg)
			}
		case <-r.ctx.Done():
			return
		}
	}
}

// NewRandomLogger creates a new dproc.Processor that logs random numbers
func NewRandomLogger(ctx context.Context, name string) dproc.Processor {
	proc := &randomLoggerProc{ctx, name, dproc.StateWaiting, make(chan dproc.Message, 1)}
	go proc.start()
	return proc
}

type randomLoggerProc struct {
	ctx   context.Context
	name  string
	state dproc.State
	inbox chan dproc.Message
}

func (r *randomLoggerProc) Name() string {
	return r.name
}

func (r *randomLoggerProc) Process(msg dproc.Message) {
	r.inbox <- msg
}

func (r *randomLoggerProc) start() {
	// var state dproc.State
	// state = dproc.StateWaiting
	var sum float64
	var count int
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			fmt.Printf("%s\t%s\t%.8f\t%d\n", time.Now().UTC().Format(time.RFC3339), TypeRandom, sum/float64(count), count)
			sum = 0
			count = 0
		case msg := <-r.inbox:
			// fmt.Printf("%s - [%s] %s %+v\n", time.Now().UTC().Format(time.RFC3339), r.name, msg.Type, msg.Values)
			switch msg.Type {
			case dproc.MessageTypeStart:
				r.state = dproc.StateRunning
			case dproc.MessageTypeStop:
				r.state = dproc.StateKilled
				ticker.Stop()
			case TypeRandom:
				sum += float64(msg.Values["random"].(float64))
				count++
			}
		case <-r.ctx.Done():
			r.state = dproc.StateKilled
			return
		}
	}
}
