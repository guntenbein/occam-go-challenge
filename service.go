package main

import (
	"sync"
	"time"
)

type IndexRateTicker struct {
	stopc chan struct{}
	runWG sync.WaitGroup
}

type RateSource interface {
	Rate() *TickerPrice
	Open() error
	Close()
}

type TimeTicker interface {
	Duration() time.Duration
	Tick() <-chan time.Time
	Stop()
}

type RateCalculator interface {
	RatePrice(sources []TickerPrice, interval time.Duration, currencyRateType Ticker) TickerPrice
}

func StartIndexRateTicker(currencyRateType Ticker, sources []RateSource,
	tt TimeTicker, pricer RateCalculator) (tickerPrisec <-chan TickerPrice, stop func()) {
	irTicker := &IndexRateTicker{
		stopc: make(chan struct{}),
		runWG: sync.WaitGroup{},
	}
	prisec := make(chan TickerPrice)
	irTicker.runWG.Add(1)
	timec := tt.Tick()
	duration := tt.Duration()
	go func() {
		for {
			select {
			case <-timec:
				relatedRates := getRateSources(sources)
				mediumPrice := pricer.RatePrice(relatedRates, duration, currencyRateType)
				sendPrice(prisec, mediumPrice, irTicker.stopc)
			case <-irTicker.stopc:
				tt.Stop()
				close(prisec)
				irTicker.runWG.Done()
				return
			}
		}
	}()
	return prisec, irTicker.stop
}

func sendPrice(prisec chan TickerPrice, mediumPrice TickerPrice, stopc chan struct{}) {
	select {
	case prisec <- mediumPrice:
	case <-stopc:
	}
}

func (s *IndexRateTicker) stop() {
	close(s.stopc)
	s.runWG.Wait()
}

func getRateSources(sources []RateSource) []TickerPrice {
	out := make([]TickerPrice, 0, len(sources))
	for _, source := range sources {
		tp := source.Rate()
		// nil means that there is no value for the price
		if tp != nil {
			out = append(out, *tp)
		}
	}
	return out
}
