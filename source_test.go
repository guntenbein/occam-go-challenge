package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReconnectedRateSource_Rate(t *testing.T) {
	subscriber := NewPredefinedPriceSubscriber(time.Millisecond*10, "100")
	source := NewReconnectedRateSource(subscriber, BTCUSDTicker, time.Millisecond*50)
	err := source.Open()
	assert.NoError(t, err)
	waitFor(t, source, ValueReady)
	rate := source.Rate()
	assert.NotNil(t, rate)

	subscriber.Close()
	waitFor(t, source, Disconnected)
	rate = source.Rate()
	assert.Nil(t, rate)

	subscriber.AllowConnecting()
	waitFor(t, source, ValueReady)
	rate = source.Rate()
	assert.NotNil(t, rate)
}

func waitFor(t *testing.T, source *ReconnectedRateSource, status int32) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	for source.Status() != status {
		select {
		case <-ctx.Done():
			t.Fatal("the source cannot be opened")
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
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
	time.Sleep(time.Millisecond * 100)
	rate = source.Rate()
	assert.Nil(t, rate)
}
