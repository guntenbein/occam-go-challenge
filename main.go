package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	rateCalculator := SimpleMediumRateCalculator{}
	sources := []RateSource{
		NewReconnectedRateSource(NewPredefinedPriceSubscriber(time.Second*10, "100"), BTCUSDTicker, time.Second*10),
		NewReconnectedRateSource(NewPredefinedPriceSubscriber(time.Second*20, "150"), BTCUSDTicker, time.Second*10),
	}
	timer := NewMinuteTimeTicker()
	out, stop := StartIndexRateTicker(BTCUSDTicker, sources, timer, rateCalculator)

	stopWG := sync.WaitGroup{}
	stopWG.Add(1)

	go func() {
		for price := range out {
			fmt.Println(price)
		}
		stopWG.Done()
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	stop()
	stopWG.Wait()
}
