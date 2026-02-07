package closer

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var globalCloser = New(os.Interrupt, syscall.SIGTERM) // global singleton

// Add registers cleanup functions to the global closer
func Add(f ...func(context.Context) error) {
	globalCloser.Add(f...)
}

// Wait blocks until all registered cleanup functions finish
func Wait() {
	globalCloser.Wait()
}

// Closer manages graceful shutdown functions
type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func(ctx context.Context) error
}

// New creates a new Closer, optionally listening for OS signals
func New(sig ...os.Signal) *Closer {
	c := &Closer{
		done: make(chan struct{}),
	}

	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}
	return c
}

// Add registers functions to be called on shutdown
func (c *Closer) Add(f ...func(context.Context) error) {
	log.Println("Add func to Closer")

	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

// Wait blocks until CloseAll completes
func (c *Closer) Wait() {
	<-c.done
}

// CloseAll executes all registered functions (concurrently) and signals completion
// func (c *Closer) CloseAll(shutdownTimeout time.Duration) {
func (c *Closer) CloseAll() {

	log.Println("Gracful shutdown - Started ... ")

	c.once.Do(func() {
		timeout := 10 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))

		// Run all cleanup functions concurrenctly
		for _, f := range funcs {
			go func(f func(ctx context.Context) error) {
				errs <- f(shutdownCtx)
			}(f)
		}

		// Collect errors
		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				log.Printf("error returned from Closer: %v\n", err)
			}
		}

	})

	log.Println("Greaceful shutdown - Success")
}
