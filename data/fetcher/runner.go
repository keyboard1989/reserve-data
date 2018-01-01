package fetcher

import (
	"time"
)

// Runner to trigger fetcher
type FetcherRunner interface {
	GetOrderbookTicker() <-chan time.Time
	GetAuthDataTicker() <-chan time.Time
	GetRateTicker() <-chan time.Time
	// Start must be non-blocking and must only return after runner
	// gets to ready state before GetOrderbookTicker() and
	// GetAuthDataTicker() get called
	Start() error
	// Stop should only be invoked when the runner is already running
	Stop() error
}

type TickerRunner struct {
	oduration time.Duration
	aduration time.Duration
	rduration time.Duration
	oclock    *time.Ticker
	aclock    *time.Ticker
	rclock    *time.Ticker
	signal    chan bool
}

func (self *TickerRunner) GetOrderbookTicker() <-chan time.Time {
	if self.oclock == nil {
		<-self.signal
	}
	return self.oclock.C
}
func (self *TickerRunner) GetAuthDataTicker() <-chan time.Time {
	if self.aclock == nil {
		<-self.signal
	}
	return self.aclock.C
}
func (self *TickerRunner) GetRateTicker() <-chan time.Time {
	if self.rclock == nil {
		<-self.signal
	}
	return self.rclock.C
}

func (self *TickerRunner) Start() error {
	self.oclock = time.NewTicker(self.oduration)
	self.signal <- true
	self.aclock = time.NewTicker(self.aduration)
	self.signal <- true
	self.rclock = time.NewTicker(self.rduration)
	self.signal <- true
	return nil
}

func (self *TickerRunner) Stop() error {
	self.oclock.Stop()
	self.aclock.Stop()
	self.rclock.Stop()
	return nil
}

func NewTickerRunner(oduration, aduration, rduration time.Duration) *TickerRunner {
	return &TickerRunner{
		oduration,
		aduration,
		rduration,
		nil,
		nil,
		nil,
		make(chan bool, 3),
	}
}
