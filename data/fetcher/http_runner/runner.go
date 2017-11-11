package http_runner

import (
	"time"
	"fmt"
)

type HttpRunner struct {
	port    int
	eticker <-chan time.Time
	bticker <-chan time.Time
}

func (self *HttpRunner) GetExchangeTicker() <-chan time.Time {
	return self.eticker
}

func (self *HttpRunner) GetBlockchainTicker() <-chan time.Time {
	return self.bticker
}

func NewHttpRunner(port int) *HttpRunner {
	echan := make(chan time.Time)
	bchan := make(chan time.Time)
	runner := HttpRunner{
		port,
		echan,
		bchan,
	}
	NewHttpRunnerServer(
		runner,
		fmt.Sprintf(":%d", port),
	)
	return &runner
}
