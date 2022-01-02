package main

import (
	"log"
	"sync"
	"time"
)

type ReconnectedRateSource struct {
	subscriber        PriceStreamSubscriber
	ticker            Ticker
	rate              TickerPrice
	functioning       bool
	reconnectInterval time.Duration
	rateChan          chan TickerPrice
	lock              sync.RWMutex
	stopc             chan struct{}
	runWG             sync.WaitGroup
}

func NewReconnectedRateSource(subscriber PriceStreamSubscriber, ticker Ticker,
	reconnectInterval time.Duration) *ReconnectedRateSource {
	return &ReconnectedRateSource{
		subscriber:        subscriber,
		ticker:            ticker,
		stopc:             make(chan struct{}),
		reconnectInterval: reconnectInterval,
	}

}

func (r *ReconnectedRateSource) Open() error {
	r.runWG.Add(2)
	go r.keepConnected()
	go r.refreshRate()
	return nil
}

func (r *ReconnectedRateSource) Rate() *TickerPrice {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if r.functioning {
		out := r.rate
		return &out
	} else {
		return nil
	}
}

func (r *ReconnectedRateSource) Close() {
	close(r.stopc)
	r.runWG.Wait()
}

func (r *ReconnectedRateSource) keepConnected() {
	for {
		select {
		case <-r.stopc:
			r.runWG.Done()
			return
		default:
			r.lock.RLock()
			functioning := r.functioning
			r.lock.RUnlock()
			if !functioning {
				if !r.connect() {
					// todo: make exponential backoff retry here
					retry := time.NewTicker(r.reconnectInterval)
					for range retry.C {
						// here we interrupt the reconnection retry
						select {
						case <-r.stopc:
							r.runWG.Done()
							return
						default:
							if r.connect() {
								retry.Stop()
								break
							}
						}
					}
				}
			} else {
				time.Sleep(r.reconnectInterval)
			}
		}
	}
}

func (r *ReconnectedRateSource) connect() bool {
	rateChan, err := r.subscriber.SubscribePriceStream(r.ticker)
	if err == nil {
		r.lock.Lock()
		r.rateChan = rateChan
		r.functioning = true
		r.lock.Unlock()
		return true
	} else {
		log.Printf("error when reconnecting to the subscriber: %s", err)
	}
	return false
}

func (r *ReconnectedRateSource) refreshRate() {
	for {
		select {
		case <-r.stopc:
			r.runWG.Done()
			return
		default:
			r.lock.RLock()
			functioning := r.functioning
			r.lock.RUnlock()
			if functioning {
				value, ok := <-r.rateChan
				r.lock.Lock()
				if ok {
					r.functioning = true
					r.rate = value
				} else {
					r.functioning = false
				}
				r.lock.Unlock()
			} else {
				time.Sleep(r.reconnectInterval)
			}
		}

	}
}
