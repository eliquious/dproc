package main

import (
	"context"
	"fmt"
	"github.com/eliquious/dproc"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	// ctx, _ := context.WithCancel(context.Background())
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

// TypeRandom is the message type for random numbers
const TypeRandom = dproc.MessageType("Random")

// RandomMessage encapulates the random number for transmission.
type RandomMessage struct {
	Value float64
}

// RandomGenerator simply generates random numbers
type RandomGenerator struct {
	Duration time.Duration
}

// Handle sends prime messages to all the child processes.
func (p *RandomGenerator) Handle(ctx context.Context, proc dproc.Process, msg dproc.Message) {
	switch msg.Type {
	default:
		fmt.Println("Unknown message type: ", msg.Type)
	case dproc.MessageTypeStop:
		log.Printf("[%s] - Exiting...", proc.Name())
	case dproc.MessageTypeStart:
		log.Printf("[%s] - Starting...", proc.Name())

		timer := time.NewTimer(p.Duration)
		msg2 := dproc.Message{Forward: false, Type: TypeRandom, Value: RandomMessage{0.0}}
		for {
			select {
			default:
				msg2.Value = RandomMessage{rand.Float64()}
				proc.Children().Dispatch(msg2)
			case <-timer.C:
				proc.SetState(dproc.StateKilled)
				log.Printf("[%s] - Exiting...", proc.Name())
				return
			}
		}
	}
}

// RandomLogger simply logs random numbers per second
type RandomLogger struct {
	Ticker *time.Ticker

	count int
}

// Handle sends prime messages to all the child processes.
func (p *RandomLogger) Handle(ctx context.Context, proc dproc.Process, msg dproc.Message) {
	switch msg.Type {
	default:
		fmt.Println("Unknown message type: ", msg.Type)
	case dproc.MessageTypeStart:
		log.Printf("[%s] - Starting...", proc.Name())
	case dproc.MessageTypeStop:
		p.Ticker.Stop()
		log.Printf("[%s] - Exiting...", proc.Name())
	case TypeRandom:
		// log.Printf("[%s] - %.5f", proc.Name(), msg.(RandomMessage).Value)

		select {
		case <-p.Ticker.C:
			log.Printf("[%s] - %d/s", proc.Name(), p.count)
			p.count = 0
		default:
			p.count++
		}
	}
}
