package main

import (
	"errors"
	"sync"
	"time"
)

type Ticker string

const (
	BTCUSDTicker Ticker = "BTC_USD"
)

type TickerPrice struct {
	Ticker Ticker
	Time   time.Time
	Price  string // decimal value. example: "0", "10", "12.2", "13.2345122"
}

type PriceStreamSubscriber interface {
	SubscribePriceStream(Ticker) (chan TickerPrice, error)
	Close()
}

type PredefinedPriceSubscriber struct {
	period         time.Duration
	price          string
	canBeConnected bool
	output         chan TickerPrice
	subscribeLock  sync.Mutex
	periodTicker   *time.Ticker
	stopc          chan struct{}
	runWG          sync.WaitGroup
}

func NewPredefinedPriceSubscriber(period time.Duration, price string) *PredefinedPriceSubscriber {
	return &PredefinedPriceSubscriber{period: period, price: price, canBeConnected: true}
}

func (p *PredefinedPriceSubscriber) SubscribePriceStream(ticker Ticker) (chan TickerPrice, error) {
	p.subscribeLock.Lock()
	defer p.subscribeLock.Unlock()
	if !p.canBeConnected {
		return nil, errors.New("the subscriber cannot be connected")
	}
	if p.output == nil {
		p.runWG.Add(1)
		p.output = make(chan TickerPrice)
		p.stopc = make(chan struct{})
		p.periodTicker = time.NewTicker(p.period)
		go func() {
			p.sendPrice(ticker, time.Now())
			for {
				select {
				case t := <-p.periodTicker.C:
					p.sendPrice(ticker, t)
				case <-p.stopc:
					p.runWG.Done()
					return
				}
			}
		}()
	}
	return p.output, nil
}

func (p *PredefinedPriceSubscriber) sendPrice(ticker Ticker, t time.Time) {
	select {
	case p.output <- TickerPrice{
		Ticker: ticker,
		Time:   t,
		Price:  p.price,
	}:
	case <-p.stopc:
	}
}

func (p *PredefinedPriceSubscriber) Close() {
	p.subscribeLock.Lock()
	defer p.subscribeLock.Unlock()
	p.periodTicker.Stop()
	close(p.stopc)
	p.runWG.Wait()
	close(p.output)
	p.output = nil
	p.canBeConnected = false
}

func (p *PredefinedPriceSubscriber) AllowConnecting() {
	p.subscribeLock.Lock()
	defer p.subscribeLock.Unlock()
	p.canBeConnected = true
}

type ErroredPriceSubscriber struct {
}

func NewErroredPriceSubscriber() *ErroredPriceSubscriber {
	return &ErroredPriceSubscriber{}
}

func (p *ErroredPriceSubscriber) SubscribePriceStream(_ Ticker) (chan TickerPrice, error) {
	tickc := make(chan TickerPrice)
	close(tickc)
	return tickc, nil
}

func (p *ErroredPriceSubscriber) Close() {
	// noop
}
