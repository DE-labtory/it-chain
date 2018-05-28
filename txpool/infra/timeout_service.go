package infra

import (
	"time"

	"github.com/it-chain/it-chain-Engine/txpool"
)

type TimeoutTicker struct {
	T         *time.Ticker
	TimeoutMs time.Duration
}

func NewTimeoutTicker(timeoutMs int64) *TimeoutTicker {
	tTicker := time.NewTicker(time.Duration(timeoutMs))
	tTicker.Stop()

	// Duration represents the time as int64 nanosecond count.
	tMs := time.Duration(timeoutMs * 100000)

	return &TimeoutTicker{
		T:         tTicker,
		TimeoutMs: tMs,
	}
}

func (t *TimeoutTicker) TransferTxStart() {
	t.T = time.NewTicker(t.TimeoutMs)
	txService := txpool.TxPeriodicTransferService{}

	go func() {
		for _ := range t.T.C {
			txService.TransferTxToLeader()
		}
	}()
}

func (t *TimeoutTicker) Stop() {
	t.T.Stop()
}
