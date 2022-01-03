package main

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	Disconnected int32 = 0
	Connected    int32 = 1
	ValueReady   int32 = 2
)

type ReconnectedRateSource struct {
	subscriber        PriceStreamSubscriber
	ticker            Ticker
	rate              TickerPrice
	status            *int32
	reconnectInterval time.Duration
	lock              sync.RWMutex
	stopc             chan struct{}
	runWG             sync.WaitGroup
}

func NewReconnectedRateSource(subscriber PriceStreamSubscriber, ticker Ticker,
	reconnectInterval time.Duration) *ReconnectedRateSource {
	status := Disconnected
	return &ReconnectedRateSource{
		subscriber:        subscriber,
		ticker:            ticker,
		reconnectInterval: reconnectInterval,
		status:            &status,
	}

}

func (r *ReconnectedRateSource) Open() error {
	r.stopc = make(chan struct{})
	atomic.StoreInt32(r.status, Disconnected)
	go r.keepConnected()
	return nil
}

func (r *ReconnectedRateSource) Rate() *TickerPrice {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if atomic.LoadInt32(r.status) != ValueReady {
		return nil
	}
	out := r.rate
	return &out
}

func (r *ReconnectedRateSource) Close() {
	close(r.stopc)
	r.runWG.Wait()
}

func (r *ReconnectedRateSource) Status() int32 {
	return atomic.LoadInt32(r.status)
}

func (r *ReconnectedRateSource) keepConnected() {
	r.runWG.Add(1)
	defer r.runWG.Done()
	for {
		select {
		case <-r.stopc:
			return
		default:
			if atomic.LoadInt32(r.status) == Disconnected {
				if err := r.connect(); err == nil {
					continue
				}
				// todo: make exponential backoff retry here
				retry := time.NewTicker(r.reconnectInterval)
			retryLoop:
				for range retry.C {
					select {
					// here we interrupt the reconnection retry
					case <-r.stopc:
						return
					default:
						if err := r.connect(); err == nil {
							retry.Stop()
							break retryLoop
						}
					}
				}
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}
	}
}

func (r *ReconnectedRateSource) connect() (err error) {
	rateChan, err := r.subscriber.SubscribePriceStream(r.ticker)
	if err == nil {
		atomic.StoreInt32(r.status, Connected)
		go r.refreshRate(rateChan)

	} else {
		log.Printf("error %s when reconnecting subscriber %+v", err, r.subscriber)
	}
	return
}

func (r *ReconnectedRateSource) refreshRate(rateChan chan TickerPrice) {
	r.runWG.Add(1)
	defer r.runWG.Done()
	for {
		select {
		case <-r.stopc:
			r.subscriber.Close()
			return
		case value, ok := <-rateChan:
			if ok {
				r.lock.Lock()
				atomic.StoreInt32(r.status, ValueReady)
				r.rate = value
				r.lock.Unlock()
			} else {
				atomic.StoreInt32(r.status, Disconnected)
				return
			}

		}
	}
}
