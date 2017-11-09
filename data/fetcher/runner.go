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

type TimestampRunner struct {
	eticker <-chan time.Time
	bticker <-chan time.Time
}

func (self *TimestampRunner) GetExchangeTicker() <-chan time.Time   { return self.eticker }
func (self *TimestampRunner) GetBlockchainTicker() <-chan time.Time { return self.bticker }

func tickTimestamp(timestamp []uint64, echan chan time.Time) {
	for _, t := range timestamp {
		echan <- time.Unix(0, int64(t)*int64(time.Millisecond))
	}
}

func NewTimestampRunner(timestamps []uint64, bduration time.Duration) *TickerRunner {
	echan := make(chan time.Time)
	go tickTimestamp(timestamps, echan)
	return &TickerRunner{
		echan,
		time.Tick(bduration),
	}
}
