package mytime

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
)

var clk = clock.New()

func SetClockInTest(t *testing.T, c clock.Clock) {
	clk = c

	t.Cleanup(func() {
		clk = clock.New()
	})
}

func Now() time.Time {
	return clk.Now()
}
