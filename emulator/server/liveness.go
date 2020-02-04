package server

import (
	"net/http"
	"time"

	"github.com/dapperlabs/flow-go-sdk/utils/liveness"
)

type LivenessTicker struct {
	collector *liveness.CheckCollector
	ticker    *time.Ticker
	done      chan bool
}

func NewLivenessTicker(tolerance time.Duration) *LivenessTicker {
	return &LivenessTicker{
		collector: liveness.NewCheckCollector(tolerance),
		ticker:    time.NewTicker(tolerance / 2),
		done:      make(chan bool, 1),
	}
}

func (l *LivenessTicker) Start() error {
	check := l.collector.NewCheck()

	for {
		select {
		case <-l.ticker.C:
			check.CheckIn()
		case <-l.done:
			return nil
		}
	}
}

func (l *LivenessTicker) Stop() {
	l.done <- true
}

func (l *LivenessTicker) Handler() http.Handler {
	return l.collector
}
