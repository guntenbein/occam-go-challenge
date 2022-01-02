package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReconnectedRateSource_Rate(t *testing.T) {
	subscriber := NewPredefinedPriceSubscriber(time.Millisecond*10, "100")
	source := NewReconnectedRateSource(subscriber, BTCUSDTicker, time.Millisecond*50)
	err := source.Open()
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 10)
	rate := source.Rate()
	assert.NotNil(t, rate)

	subscriber.Close()
	time.Sleep(time.Millisecond * 50)
	rate = source.Rate()
	assert.Nil(t, rate)

	subscriber.AllowConnecting()
	time.Sleep(time.Millisecond * 60)
	rate = source.Rate()
	assert.NotNil(t, rate)
}

func TestReconnectedRateSource_Rate_StreamAlwaysClose(t *testing.T) {
	subscriber := NewErroredPriceSubscriber()
	source := NewReconnectedRateSource(subscriber, BTCUSDTicker, time.Millisecond*50)
	err := source.Open()
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 10)
	rate := source.Rate()
	assert.Nil(t, rate)

	subscriber.Close()
	time.Sleep(time.Millisecond * 50)
	rate = source.Rate()
	assert.Nil(t, rate)
}
