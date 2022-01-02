package main

import (
	"testing"
	"time"
)

func TestStartIndexRateTicker(t *testing.T) {
	t.Run("normal functioning", func(t *testing.T) {
		NewGenericTimeTicker(time.Millisecond * 50)
	})
}
