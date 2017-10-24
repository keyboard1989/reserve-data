package fetcher

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Fetcher struct {
	storage    Storage
	exchanges  []Exchange
	ticker     <-chan time.Time
	blockchain Blockchain
	bticker    <-chan time.Time
	rmaddr     common.Address
}

func NewFetcher(
	storage Storage,
	duration, bduration time.Duration,
	address common.Address) *Fetcher {
	return &Fetcher{
		storage:    storage,
		exchanges:  []Exchange{},
		ticker:     time.Tick(duration),
		blockchain: nil,
		bticker:    time.Tick(bduration),
		rmaddr:     address,
	}
}

func (self *Fetcher) SetBlockchain(blockchain Blockchain) {
	self.blockchain = blockchain
}

func (self *Fetcher) AddExchange(exchange Exchange) {
	self.exchanges = append(self.exchanges, exchange)
}

func (self *Fetcher) fetchingFromExchanges() {
	for _ = range self.ticker {
		self.fetchAllFromExchanges()
	}
}

func (self *Fetcher) fetchingFromBlockchain() {
	for _ = range self.bticker {
		self.fetchAllFromBlockchain()
	}
}

func (self *Fetcher) Run() error {
	go self.fetchingFromExchanges()
	go self.fetchingFromBlockchain()
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
	err := self.storage.StorePrice(data.GetData())
	if err != nil {
		log.Printf("Storing data failed: %s\n", err)
	}
}

func (self *Fetcher) fetchAllBalances() {
	data, err := self.blockchain.FetchBalanceData(self.rmaddr)
	if err != nil {
		log.Printf("Fetching data from blockchain failed: %s\n", err)
	}
	err = self.storage.StoreBalance(data)
	fmt.Printf("balance data: %v\n", data)
	if err != nil {
		log.Printf("Storing balance data failed: %s\n", err)
	}
}

func (self *Fetcher) fetchAllFromExchanges() {
	fmt.Printf("Fetching data...")
	self.fetchAllPrices()
}

func (self *Fetcher) fetchAllFromBlockchain() {
	fmt.Printf("Fetching data from blockchain...")
	self.fetchAllBalances()
}
