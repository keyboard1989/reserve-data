package http_runner

import (
	"time"
	"errors"
	"fmt"
)

type HttpRunner struct {
	port    int
	eticker <-chan time.Time
	bticker <-chan time.Time
	server  *HttpRunnerServer
}

func (self *HttpRunner) GetExchangeTicker() <-chan time.Time {
	return self.eticker
}

func (self *HttpRunner) GetBlockchainTicker() <-chan time.Time {
	return self.bticker
}

func (self *HttpRunner) Start() error {
	if self.server != nil {
		return errors.New("runner start already")
	} else {
		self.server = NewHttpRunnerServer(self, fmt.Sprintf(":%d", self.port))
		return self.server.Start()
	}
}

func (self *HttpRunner) Stop() error {
	if self.server != nil {
		err := self.server.Stop()
		self.server = nil
		return err
	} else {
		return errors.New("runner stop already")
	}
}

func NewHttpRunner(port int) *HttpRunner {
	echan := make(chan time.Time)
	bchan := make(chan time.Time)
	runner := HttpRunner{
		port,
		echan,
		bchan,
		nil,
	}
	runner.Start()
	return &runner
}
