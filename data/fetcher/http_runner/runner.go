package http_runner

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// HttpRunner is an implementation of FetcherRunner
// that run a HTTP server and tick when it receives request to a certain endpoints.
type HttpRunner struct {
	port                    int
	oticker                 chan time.Time
	aticker                 chan time.Time
	rticker                 chan time.Time
	bticker                 chan time.Time
	tticker                 chan time.Time
	rsticker                chan time.Time
	lticker                 chan time.Time
	tradeLogProcessorTicker chan time.Time
	catLogProcessorTicker   chan time.Time
	globalDataTicker        chan time.Time
	server                  *HttpRunnerServer
}

// GetGlobalDataTicker returns the global data ticker.
func (self *HttpRunner) GetGlobalDataTicker() <-chan time.Time {
	return self.globalDataTicker
}

// GetTradeLogProcessorTicker returns the trade log processor ticker.
func (self *HttpRunner) GetTradeLogProcessorTicker() <-chan time.Time {
	return self.tradeLogProcessorTicker
}

// GetCatLogProcessorTicker returns the cat log processor ticker.
func (self *HttpRunner) GetCatLogProcessorTicker() <-chan time.Time {
	return self.catLogProcessorTicker
}

// GetLogTicker returns the log ticker.
func (self *HttpRunner) GetLogTicker() <-chan time.Time {
	return self.lticker
}

// GetBlockTicker returns the block ticker.
func (self *HttpRunner) GetBlockTicker() <-chan time.Time {
	return self.bticker
}

// GetOrderbookTicker returns the order book ticker.
func (self *HttpRunner) GetOrderbookTicker() <-chan time.Time {
	return self.oticker
}

// GetAuthDataTicker returns the auth data ticker.
func (self *HttpRunner) GetAuthDataTicker() <-chan time.Time {
	return self.aticker
}

// GetRateTicker returns the rate ticker.
func (self *HttpRunner) GetRateTicker() <-chan time.Time {
	return self.rticker
}

// GetTradeHistoryTicker returns the trade history ticker.
func (self *HttpRunner) GetTradeHistoryTicker() <-chan time.Time {
	return self.tticker
}

// GetReserveRatesTicker returns the reserve rates ticker.
func (self *HttpRunner) GetReserveRatesTicker() <-chan time.Time {
	return self.rsticker
}

// Start initializes and starts the ticker HTTP server.
// It returns an error if the server is started already.
// It is guaranteed that the HTTP server is ready to serve request after
// this method is returned.
func (self *HttpRunner) Start() error {
	if self.server != nil {
		return errors.New("runner start already")
	} else {
		self.server = NewHttpRunnerServer(self, fmt.Sprintf(":%d", self.port))
		go func() {
			err := self.server.Start()
			if err != nil {
				log.Printf("Http server for runner couldn't start or get stopped. Error: %s", err)
			}
		}()
		return nil
	}
}

// Stop stops the HTTP server. It returns an error if the server is already stopped.
func (self *HttpRunner) Stop() error {
	if self.server != nil {
		err := self.server.Stop()
		self.server = nil
		return err
	} else {
		return errors.New("runner stop already")
	}
}

// NewHttpRunner creates a new instance of HttpRunner.
// The HTTP server is also started after creation.
func NewHttpRunner(port int) (*HttpRunner, error) {
	ochan := make(chan time.Time)
	achan := make(chan time.Time)
	rchan := make(chan time.Time)
	bchan := make(chan time.Time)
	tchan := make(chan time.Time)
	rschan := make(chan time.Time)
	lchan := make(chan time.Time)
	tradeLogProcessorChan := make(chan time.Time)
	catLogProcessorChan := make(chan time.Time)
	globalDataChan := make(chan time.Time)
	runner := &HttpRunner{
		port,
		ochan,
		achan,
		rchan,
		bchan,
		tchan,
		rschan,
		lchan,
		tradeLogProcessorChan,
		catLogProcessorChan,
		globalDataChan,
		nil,
	}
	if err := runner.Start(); err != nil {
		return nil, err
	}
	return runner, nil
}
