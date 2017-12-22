package fetcher

import (
	"time"
)

// Runner to trigger fetcher
type FetcherRunner interface {
	GetOrderbookTicker() <-chan time.Time
	GetBlockchainTicker() <-chan time.Time
	// Start must be non-blocking and must only return after runner
	// gets to ready state before GetExchangeTicker() and
	// GetBlockchainTicker() get called
	Start() error
	// Stop should only be invoked when the runner is already running
	Stop() error
}

type TickerRunner struct {
	eduration time.Duration
	bduration time.Duration
	eclock    *time.Ticker
	bclock    *time.Ticker
	signal    chan bool
}

func (self *TickerRunner) GetOrderbookTicker() <-chan time.Time {
	if self.eclock == nil {
		<-self.signal
	}
	return self.eclock.C
}
func (self *TickerRunner) GetBlockchainTicker() <-chan time.Time {
	if self.bclock == nil {
		<-self.signal
	}
	return self.bclock.C
}

func (self *TickerRunner) Start() error {
	self.eclock = time.NewTicker(self.eduration)
	self.signal <- true
	self.bclock = time.NewTicker(self.bduration)
	self.signal <- true
	return nil
}

func (self *TickerRunner) Stop() error {
	self.eclock.Stop()
	self.bclock.Stop()
	return nil
}

func NewTickerRunner(eduration, bduration time.Duration) *TickerRunner {
	return &TickerRunner{
		eduration,
		bduration,
		nil,
		nil,
		make(chan bool, 2),
	}
}

type TimestampRunner struct {
	bduration  time.Duration
	eticker    chan time.Time
	bclock     *time.Ticker
	timestamps []uint64
}

func (self *TimestampRunner) GetOrderbookTicker() <-chan time.Time {
	return (<-chan time.Time)(self.eticker)
}
func (self *TimestampRunner) GetBlockchainTicker() <-chan time.Time { return self.bclock.C }

func (self *TimestampRunner) Start() error {
	self.bclock = time.NewTicker(self.bduration)
	go tickTimestamp(self.timestamps, self.eticker)
	return nil
}

func (self *TimestampRunner) Stop() error {
	// todo: stop echan, still can't search for close() chan time.Time function
	return nil
}

func tickTimestamp(timestamp []uint64, echan chan time.Time) {
	for _, t := range timestamp {
		echan <- time.Unix(0, int64(t)*int64(time.Millisecond))
	}
}

func NewTimestampRunner(timestamps []uint64, bduration time.Duration) *TimestampRunner {
	echan := make(chan time.Time)
	return &TimestampRunner{
		bduration,
		echan,
		nil,
		timestamps,
	}
}
