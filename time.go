package main

import (
	"sync"
	"time"
)

type MinuteTimeTicker struct {
	minuteTimec chan time.Time
	heartbeat   *time.Ticker
	stopc       chan struct{}
	runWG       sync.WaitGroup
}

func NewMinuteTimeTicker() *MinuteTimeTicker {
	return &MinuteTimeTicker{stopc: make(chan struct{})}
}

func (m *MinuteTimeTicker) Duration() time.Duration {
	return time.Minute
}

func (m *MinuteTimeTicker) Tick() chan time.Time {
	if m.minuteTimec == nil {
		m.minuteTimec = make(chan time.Time)
		m.heartbeat = time.NewTicker(100 * time.Millisecond)
		m.runWG.Add(1)
		go func() {
			defer m.runWG.Done()
			lastMinute := time.Time{}
			for {
				select {
				case t := <-m.heartbeat.C:
					if t.Second() == 0 && t.Sub(lastMinute) > time.Second {
						lastMinute = t
						m.minuteTimec <- t
					}
				case <-m.stopc:
					return
				}
			}
		}()

	}
	return m.minuteTimec
}

func (m *MinuteTimeTicker) Stop() {
	if m.minuteTimec != nil {
		m.heartbeat.Stop()
		close(m.stopc)
		m.runWG.Wait()
		close(m.minuteTimec)
	}
}
