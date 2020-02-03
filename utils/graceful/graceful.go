package graceful

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Routine interface {
	Start() error
	Stop()
}

type Group struct {
	routines []Routine
	done     map[int]bool
	err      error
	mut      sync.Mutex
}

func NewGroup() *Group {
	return &Group{
		routines: make([]Routine, 0),
		done:     make(map[int]bool),
	}
}

func (g *Group) Add(r Routine) {
	g.mut.Lock()
	defer g.mut.Unlock()
	g.routines = append(g.routines, r)
}

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
