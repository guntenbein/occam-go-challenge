package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimpleMediumRateCalculator_RatePrice(t *testing.T) {
	t.Run("normal work", func(t *testing.T) {
		mediumPrice := SimpleMediumRateCalculator{}.RatePrice([]TickerPrice{
			{
				Ticker: BTCUSDTicker,
				Time:   time.Now(),
				Price:  "100.1",
			},
			{
				Ticker: BTCUSDTicker,
				Time:   time.Now(),
				Price:  "100.8",
			},
			{
				Ticker: BTCUSDTicker,
				Time:   time.Now(),
				Price:  "99.8",
			},
		}, time.Minute, BTCUSDTicker)
		assert.Equal(t, BTCUSDTicker, mediumPrice.Ticker)
		assert.Equal(t, "100.23333333333333", mediumPrice.Price)
	})
	t.Run("normal work, some prices are outdated", func(t *testing.T) {
		mediumPrice := SimpleMediumRateCalculator{}.RatePrice([]TickerPrice{
			{
				Ticker: BTCUSDTicker,
				Time:   time.Now(),
				Price:  "100.1",
			},
			{
				Ticker: BTCUSDTicker,
				Time:   time.Now().Add(-time.Hour),
				Price:  "100.8",
			},
			{
				Ticker: BTCUSDTicker,
				Time:   time.Now(),
				Price:  "99.8",
			},
		}, time.Minute, BTCUSDTicker)
		assert.Equal(t, BTCUSDTicker, mediumPrice.Ticker)
		assert.Equal(t, "99.94999999999999", mediumPrice.Price)
	})
}
