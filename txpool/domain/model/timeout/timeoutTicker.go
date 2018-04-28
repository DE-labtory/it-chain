package timeout

import (
	"time"
)

type TimeoutTicker struct{
	T *time.Ticker
	TimeoutMs time.Duration
}

func NewTimeoutTicker(timeoutMs int64) *TimeoutTicker{
	tTicker := time.NewTicker(time.Duration(timeoutMs))
	// Duration represents the time as int64 nanosecond count.
	tMs := time.Duration(timeoutMs*100000)

	tTicker.Stop()

	return &TimeoutTicker{
		T:         tTicker,
		TimeoutMs: tMs,
	}
}

func (t *TimeoutTicker) Start() {
	t.T = time.NewTicker(t.TimeoutMs)
}

func (t *TimeoutTicker) Stop() {
	t.T.Stop()
}