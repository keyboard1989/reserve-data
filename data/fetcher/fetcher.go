package fetcher

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Fetcher struct {
	storage    Storage
	exchanges  []Exchange
	ticker     <-chan time.Time
	blockchain Blockchain
	bticker    <-chan time.Time
	rmaddr     ethereum.Address
}

func NewFetcher(
	storage Storage,
	duration, bduration time.Duration,
	address ethereum.Address) *Fetcher {
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

func (self *Fetcher) fetchEBalanceFromExchange(wg *sync.WaitGroup, exchange Exchange, data *sync.Map) {
	defer wg.Done()
	exdata, err := exchange.FetchEBalanceData()
	if err != nil {
		log.Printf("Fetching exchange balances from %s failed: %v\n", exchange.Name(), err)
	}
	data.Store(exchange.ID(), exdata)
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

func (self *Fetcher) fetchAllPrices(w *sync.WaitGroup) {
	defer w.Done()
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

func (self *Fetcher) fetchAllEBalances(w *sync.WaitGroup) {
	defer w.Done()
	data := sync.Map{}
	// start fetching
	wait := sync.WaitGroup{}
	for _, exchange := range self.exchanges {
		wait.Add(1)
		go self.fetchEBalanceFromExchange(&wait, exchange, &data)
	}
	wait.Wait()
	ebalances := map[common.ExchangeID]common.EBalanceEntry{}
	data.Range(func(key, value interface{}) bool {
		ebalances[key.(common.ExchangeID)] = value.(common.EBalanceEntry)
		return true
	})
	err := self.storage.StoreEBalance(ebalances)
	if err != nil {
		log.Printf("Storing exchange balances failed: %s\n", err)
	}
}

func (self *Fetcher) fetchAllBalances(w *sync.WaitGroup) {
	defer w.Done()
	data, err := self.blockchain.FetchBalanceData(self.rmaddr)
	if err != nil {
		log.Printf("Fetching data from blockchain failed: %s\n", err)
	}
	err = self.storage.StoreBalance(data)
	// fmt.Printf("balance data: %v\n", data)
	if err != nil {
		log.Printf("Storing balance data failed: %s\n", err)
	}
}

func (self *Fetcher) fetchAllFromExchanges() {
	fmt.Printf("Fetching data...")
	wait := sync.WaitGroup{}
	wait.Add(1)
	go self.fetchAllPrices(&wait)
	wait.Add(1)
	go self.fetchAllEBalances(&wait)
	wait.Wait()
}

func (self *Fetcher) fetchAllFromBlockchain() {
	fmt.Printf("Fetching data from blockchain...")
	wait := sync.WaitGroup{}
	wait.Add(1)
	self.fetchAllBalances(&wait)
	wait.Wait()
}
