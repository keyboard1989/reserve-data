package fetcher

import (
	"fmt"
	"github.com/KyberNetwork/reserve-data/market"
	"log"
	"strconv"
	"sync"
	"time"
)

type IntervalFetcher struct {
	storage   market.Storage
	exchanges []Exchange
	pairs     []TokenPair
	ticker    <-chan time.Time
}

func NewIntervalFetcher(storage market.Storage, duration time.Duration) *IntervalFetcher {
	return &IntervalFetcher{
		storage:   storage,
		exchanges: []Exchange{},
		pairs:     []TokenPair{},
		ticker:    time.Tick(duration),
	}
}

func (self *IntervalFetcher) AddExchange(ex Exchange) {
	self.exchanges = append(self.exchanges, ex)
}

func (self *IntervalFetcher) AddTokenPair(p TokenPair) {
	self.pairs = append(self.pairs, p)
}

func (self *IntervalFetcher) Run() {
	for _ = range self.ticker {
		self.FetchDataFromAllExchanges()
	}
}

func (self *IntervalFetcher) FetchDataFromOneExchanges(wg *sync.WaitGroup, exchange Exchange, data *ConcurrentAllPriceData) {
	defer wg.Done()
	pairs := []TokenPair{}
	for _, pair := range self.pairs {
		base := pair.Base
		quote := pair.Quote
		if base.SymbolOnExchange(exchange) != "" && quote.SymbolOnExchange(exchange) != "" {
			pairs = append(pairs, pair)
		}
	}
	if len(pairs) > 0 {
		_, exdata := exchange.FetchData(pairs)
		for pairStr, onePairData := range exdata {
			data.SetOnePairData(exchange, pairStr, onePairData)
		}
	}
}

func (self *IntervalFetcher) FetchDataFromAllExchanges() {
	fmt.Printf("Fetching data...")
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	data := NewConcurrentAllPriceData()
	data.SetTimestamp(strconv.Itoa(int(timestamp)))
	data.SetVersion(0)
	// start fetching
	wait := sync.WaitGroup{}
	for _, exchange := range self.exchanges {
		wait.Add(1)
		go self.FetchDataFromOneExchanges(&wait, exchange, data)
	}
	wait.Wait()
	fmt.Printf("Data: %v\n", data)
	err := self.storage.StoreNewData(strconv.Itoa(int(timestamp)), data.GetData())
	if err != nil {
		log.Printf("Storing data failed: %v\n", err)
	}
}
