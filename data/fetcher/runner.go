package fetcher

import (
	"time"
)

// Runner to trigger fetcher
type FetcherRunner interface {
	GetExchangeTicker() <-chan time.Time
	GetBlockchainTicker() <-chan time.Time
}

type TickerRunner struct {
	eticker <-chan time.Time
	bticker <-chan time.Time
}

func (self *TickerRunner) GetExchangeTicker() <-chan time.Time   { return self.eticker }
func (self *TickerRunner) GetBlockchainTicker() <-chan time.Time { return self.bticker }

func NewTickerRunner(eduration, bduration time.Duration) *TickerRunner {
	return &TickerRunner{
		time.Tick(eduration),
		time.Tick(bduration),
	}
}
