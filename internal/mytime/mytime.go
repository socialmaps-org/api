package mytime

import (
	"time"

	"github.com/benbjohnson/clock"
)

var clk = clock.New()

func SetClock(c clock.Clock) {
	clk = c
}

func Now() time.Time {
	return clk.Now()
}
