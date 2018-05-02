package storagecontroller

import (
	"time"
)

type StorageControllerRunner interface {
	GetPriceBucketTicker() <-chan time.Time
	GetAuthBucketTicker() <-chan time.Time
	Start() error
	Stop() error
}

type ControllerTickerRunner struct {
	priceDuration time.Duration
	authDuration  time.Duration

	priceClock *time.Ticker
	authClock  *time.Ticker
	signal     chan bool
}

func (self *ControllerTickerRunner) GetPriceBucketTicker() <-chan time.Time {
	if self.priceClock == nil {
		<-self.signal
	}
	return self.priceClock.C
}

func (self *ControllerTickerRunner) GetAuthBucketTicker() <-chan time.Time {
	if self.authClock == nil {
		<-self.signal
	}
	return self.authClock.C
}

func (self *ControllerTickerRunner) Start() error {
	self.priceClock = time.NewTicker(self.priceDuration)
	self.signal <- true
	self.authClock = time.NewTicker(self.authDuration)
	self.signal <- true
	return nil
}

func (self *ControllerTickerRunner) Stop() error {
	self.priceClock.Stop()
	self.authClock.Stop()
	return nil
}

func NewStorageControllerTickerRunner(
	priceDuration, authDuration time.Duration) *ControllerTickerRunner {
	return &ControllerTickerRunner{
		priceDuration,
		authDuration,
		nil,
		nil,
		make(chan bool, 2),
	}
}
