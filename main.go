package main

import (
	"fmt"
	"log"
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
	for _, src := range sources {
		if err := src.Open(); err != nil {
			log.Fatalf("cannot open a stream source for a ticker: %s; error: %s", src, err)
		}
	}
	defer func() {
		for _, src := range sources {
			src.Close()
		}
	}()
	timer := NewMinuteTimeTicker()
	out, stop := StartIndexRateTicker(BTCUSDTicker, sources, timer, rateCalculator)

	stopWG := sync.WaitGroup{}
	stopWG.Add(1)

	go func() {
		for price := range out {
			fmt.Println(fmt.Sprintf("%s : %s", price.Time.Format("2006-01-02 15:04:05"), price.Price))
		}
		stopWG.Done()
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	stop()
	stopWG.Wait()
}
