package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type SimpleMediumRateCalculator struct{}

func (s SimpleMediumRateCalculator) RatePrice(sources []TickerPrice, interval time.Duration, currencyRateType Ticker) TickerPrice {
	now := time.Now()
	var summRate, weight float64
	for _, src := range sources {
		// filter out the outdated rates
		if src.Time.Add(interval).After(now) && src.Ticker == currencyRateType {
			rate, err := strconv.ParseFloat(src.Price, 64)
			if err != nil {
				log.Printf("cannot parse the rate value %s", src.Price)
			} else {
				summRate += rate
				weight++
			}
		}
	}
	return TickerPrice{
		Ticker: currencyRateType,
		Time:   now,
		Price:  calculateMedium(summRate, weight),
	}
}

func calculateMedium(rate float64, weight float64) string {
	if weight == 0 {
		return "undefined"
	}
	return fmt.Sprintf("%.14f", rate/weight)
}
