package fetcher

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Fetcher struct {
	storage   Storage
	exchanges []Exchange
	ticker    <-chan time.Time
}

func NewFetcher(storage Storage, duration time.Duration) *Fetcher {
	return &Fetcher{
		storage:   storage,
		exchanges: []Exchange{},
		ticker:    time.Tick(duration),
	}
}

func (self *Fetcher) AddExchange(exchange Exchange) {
	self.exchanges = append(self.exchanges, exchange)
}

func (self *Fetcher) Run() error {
	for _ = range self.ticker {
		self.fetchAll()
	}
	return nil
}

func (self *Fetcher) fetchPriceFromExchange(wg *sync.WaitGroup, exchange Exchange, data *ConcurrentAllPriceData) {
	defer wg.Done()
	exdata, err := exchange.FetchPriceData()
	if err != nil {
		log.Printf("Fetching data from %s failed: %v\n", exchange.Name(), err)
	}
	for pair, exchangeData := range exdata {
		data.SetOnePrice(exchange.ID(), pair, exchangeData)
	}
}

func (self *Fetcher) fetchAllPrices() {
	data := NewConcurrentAllPriceData()
	// start fetching
	wait := sync.WaitGroup{}
	for _, exchange := range self.exchanges {
		wait.Add(1)
		go self.fetchPriceFromExchange(&wait, exchange, data)
	}
	wait.Wait()
	fmt.Printf("Data: %v\n", data)
	err := self.storage.StorePrice(data.GetData())
	if err != nil {
		log.Printf("Storing data failed: %v\n", err)
	}
}

func (self *Fetcher) fetchAll() {
	fmt.Printf("Fetching data...")
	self.fetchAllPrices()
}
