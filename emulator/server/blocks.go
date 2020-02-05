package server

import (
	"time"
)

type BlocksTicker struct {
	backend *Backend
	ticker  *time.Ticker
	done    chan bool
}

func NewBlocksTicker(
	backend *Backend,
	blockTime time.Duration,
) *BlocksTicker {
	return &BlocksTicker{
		backend: backend,
		ticker:  time.NewTicker(blockTime),
		done:    make(chan bool, 1),
	}
}

func (t *BlocksTicker) Start() error {
	for {
		select {
		case <-t.ticker.C:
			t.backend.commitBlock()
		case <-t.done:
			return nil
		}
	}
}

func (t *BlocksTicker) Stop() {
	t.done <- true
}
