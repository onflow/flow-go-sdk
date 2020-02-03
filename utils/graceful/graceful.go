package graceful

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// A Routine is a long-running task that can be started and stopped.
type Routine interface {
	// Start starts the routine and blocks until it completes.
	// This function will return an error if the routine stops unexpectedly.
	Start() error
	// Stop stops the running routine.
	// Calling this function will cause Start to return without an error.
	Stop()
}

// A Group synchronizes the startup and shutdown behaviour of multiple routines.
//
// All registered routines start running simultaneously. If any routine stops early,
// all other routines in the group are gracefully shut down.
//
// Similarly, a call to Stop will gracefully shut down all routines in the group.
//
// A group is thread-safe in the sense that it run multiple routines and track
// their status without data races. However, a single group is not safe to use across
// multiple host threads. Start and Stop should only be called at most once.
type Group struct {
	routines []Routine
	done     map[int]bool
	err      error
	mut      sync.Mutex
}

// NewGroup creates a new routine group.
func NewGroup() *Group {
	return &Group{
		routines: make([]Routine, 0),
		done:     make(map[int]bool),
	}
}

// Add adds a routine to the group.
func (g *Group) Add(r Routine) {
	g.mut.Lock()
	defer g.mut.Unlock()
	g.routines = append(g.routines, r)
}

// Start starts all routines in the group.
//
// This function blocks until all routines in the group have finished running.
//
// If one or more routines stop early and produce an error, this function
// returns with the first error that was captured.
//
// This function also captures any SIGINT, SIGQUIT or SIGTERM signals received
// from the host environment. After receiving any of these signals, all routines
// in the group are shut down gracefully.
func (g *Group) Start() error {
	wg := &sync.WaitGroup{}

	g.mut.Lock()

	for i, r := range g.routines {
		wg.Add(1)
		go func(i int, r Routine) {
			err := r.Start()

			g.markDone(i, err)

			if err != nil {
				g.Stop()
			}

			wg.Done()
		}(i, r)
	}

	g.mut.Unlock()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		<-c
		g.Stop()
	}()

	wg.Wait()

	return g.shutdownErr()
}

// Stop gracefully shuts down all routines in the group.
//
// In any scenario, a group guarantees that each routine's Stop function
// is called at most once.
func (g *Group) Stop() {
	for i, r := range g.routines {
		if g.isDone(i) {
			continue
		}

		r.Stop()
		g.markDone(i, nil)
	}
}

func (g *Group) markDone(i int, err error) {
	g.mut.Lock()
	defer g.mut.Unlock()

	if g.err == nil {
		g.err = err
	}

	g.done[i] = true
}

func (g *Group) isDone(i int) bool {
	g.mut.Lock()
	defer g.mut.Unlock()
	return g.done[i]
}

func (g *Group) shutdownErr() error {
	g.mut.Lock()
	defer g.mut.Unlock()
	return g.err
}
