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
		Price:  fmt.Sprintf("%.14f", summRate/weight),
	}
}
