package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMinuteTimeTicker_Tick(t *testing.T) {
	mtt := NewMinuteTimeTicker()
	tickc := mtt.Tick()
	tm := <-tickc
	assert.Equal(t, 0, tm.Second())
}

func TestMinuteTimeTicker_Stop(t *testing.T) {
	mtt := NewMinuteTimeTicker()
	tickc := mtt.Tick()
	mtt.Stop()
	checkTicker := time.NewTicker(time.Second)
	timeElapsed := false
	for !timeElapsed {
		select {
		case <-checkTicker.C:
			timeElapsed = true
		case _, ok := <-tickc:
			if !ok {
				return
			}
		}
	}
	t.Fatal("the minute ticker did not stop properly")
}
