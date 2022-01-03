package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartIndexRateTicker(t *testing.T) {
	t.Run("normal status", func(t *testing.T) {
		ticker := NewGenericTimeTicker(time.Millisecond * 100)
		sources := []RateSource{
			NewReconnectedRateSource(NewPredefinedPriceSubscriber(time.Millisecond*10, "100"), BTCUSDTicker, time.Second*10),
			NewReconnectedRateSource(NewPredefinedPriceSubscriber(time.Millisecond*10, "150"), BTCUSDTicker, time.Second*10),
		}
		for _, src := range sources {
			if err := src.Open(); err != nil {
				t.Fatalf("cannot open a stream source for a ticker: %s; error: %s", src, err)
			}
		}
		defer func() {
			for _, src := range sources {
				src.Close()
			}
		}()
		rateCalculator := SimpleMediumRateCalculator{}
		out, stop := StartIndexRateTicker(BTCUSDTicker, sources, ticker, rateCalculator)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	waitForResult:
		for {
			select {
			case <-ctx.Done():
				t.Fatal("the source cannot be opened")
			case val, ok := <-out:
				if ok {
					assert.Equal(t, BTCUSDTicker, val.Ticker)
					assert.Equal(t, "125.00000000000000", val.Price)
					break waitForResult
				} else {
					t.Fatal("the source was closed without sending a value")
				}
			}
		}
		cancel()
		stop()
	})
	t.Run("no sources are connected", func(t *testing.T) {
		ticker := NewGenericTimeTicker(time.Millisecond * 100)
		sources := []RateSource{
			NewReconnectedRateSource(&ErroredPriceSubscriber{}, BTCUSDTicker, time.Second*10),
		}
		for _, src := range sources {
			if err := src.Open(); err != nil {
				t.Fatalf("cannot open a stream source for a ticker: %s; error: %s", src, err)
			}
		}
		defer func() {
			for _, src := range sources {
				src.Close()
			}
		}()
		rateCalculator := SimpleMediumRateCalculator{}
		out, stop := StartIndexRateTicker(BTCUSDTicker, sources, ticker, rateCalculator)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	waitForResult:
		for {
			select {
			case <-ctx.Done():
				t.Fatal("the source cannot be opened")
			case val, ok := <-out:
				if ok {
					assert.Equal(t, BTCUSDTicker, val.Ticker)
					assert.Equal(t, "undefined", val.Price)
					break waitForResult
				} else {
					t.Fatal("the source was closed without sending a value")
				}
			}
		}
		cancel()
		stop()
	})
}
