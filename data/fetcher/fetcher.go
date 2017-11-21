package fetcher

import (
	"log"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Fetcher struct {
	storage    Storage
	exchanges  []Exchange
	blockchain Blockchain
	runner     FetcherRunner
	rmaddr     ethereum.Address
}

func NewFetcher(
	storage Storage,
	runner FetcherRunner,
	address ethereum.Address) *Fetcher {
	return &Fetcher{
		storage:    storage,
		exchanges:  []Exchange{},
		blockchain: nil,
		runner:     runner,
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
	for {
		log.Printf("waiting for signal from runner for exchange ticker")
		t := <-self.runner.GetExchangeTicker()
		log.Printf("got signal in exchange ticker")
		self.fetchAllFromExchanges(common.TimeToTimepoint(t))
		log.Printf("fetched data from exchanges")
	}
}

func (self *Fetcher) fetchingFromBlockchain() {
	for {
		t := <-self.runner.GetBlockchainTicker()
		self.fetchAllFromBlockchain(common.TimeToTimepoint(t))
	}
}

func (self *Fetcher) Stop() error {
	return self.runner.Stop()
}

func (self *Fetcher) Run() error {
	log.Printf("Fetcher runner is starting...")
	self.runner.Start()
	log.Printf("Fetcher runner is running...")
	go self.fetchingFromExchanges()
	go self.fetchingFromBlockchain()
	return nil
}

func (self *Fetcher) fetchEBalanceFromExchange(wg *sync.WaitGroup, exchange Exchange, data *sync.Map, timepoint uint64) {
	defer wg.Done()
	exdata, err := exchange.FetchEBalanceData(timepoint)
	if err != nil {
		log.Printf("Fetching exchange balances from %s failed: %v\n", exchange.Name(), err)
	}
	data.Store(exchange.ID(), exdata)
}

func (self *Fetcher) fetchPriceFromExchange(wg *sync.WaitGroup, exchange Exchange, data *ConcurrentAllPriceData, timepoint uint64) {
	defer wg.Done()
	exdata, err := exchange.FetchPriceData(timepoint)
	if err != nil {
		log.Printf("Fetching data from %s failed: %v\n", exchange.Name(), err)
	}
	for pair, exchangeData := range exdata {
		data.SetOnePrice(exchange.ID(), pair, exchangeData)
	}
}

func (self *Fetcher) fetchAllPrices(w *sync.WaitGroup, timepoint uint64) {
	defer w.Done()
	data := NewConcurrentAllPriceData()
	// start fetching
	wait := sync.WaitGroup{}
	for _, exchange := range self.exchanges {
		wait.Add(1)
		go self.fetchPriceFromExchange(&wait, exchange, data, timepoint)
	}
	wait.Wait()
	err := self.storage.StorePrice(data.GetData(), timepoint)
	if err != nil {
		log.Printf("Storing data failed: %s\n", err)
	}
}

func (self *Fetcher) fetchAllEBalances(w *sync.WaitGroup, timepoint uint64) {
	defer w.Done()
	data := sync.Map{}
	// start fetching
	wait := sync.WaitGroup{}
	for _, exchange := range self.exchanges {
		wait.Add(1)
		go self.fetchEBalanceFromExchange(&wait, exchange, &data, timepoint)
	}
	wait.Wait()
	ebalances := map[common.ExchangeID]common.EBalanceEntry{}
	data.Range(func(key, value interface{}) bool {
		ebalances[key.(common.ExchangeID)] = value.(common.EBalanceEntry)
		return true
	})
	err := self.storage.StoreEBalance(ebalances, timepoint)
	if err != nil {
		log.Printf("Storing exchange balances failed: %s\n", err)
	}
}

func (self *Fetcher) fetchOrderFromExchange(wg *sync.WaitGroup, exchange Exchange, data *sync.Map, timepoint uint64) {
	defer wg.Done()
	orderData, err := exchange.FetchOrderData(timepoint)
	if err != nil {
		log.Printf("Fetching orders from %s failed: %v\n", exchange.Name(), err)
	}
	data.Store(exchange.ID(), orderData)
}

func (self *Fetcher) fetchAllOrders(w *sync.WaitGroup, timepoint uint64) {
	defer w.Done()
	data := sync.Map{}
	// start fetching
	wait := sync.WaitGroup{}
	for _, exchange := range self.exchanges {
		wait.Add(1)
		go self.fetchOrderFromExchange(&wait, exchange, &data, timepoint)
	}
	wait.Wait()
	orders := common.AllOrderEntry{}
	data.Range(func(key, value interface{}) bool {
		orders[key.(common.ExchangeID)] = value.(common.OrderEntry)
		return true
	})
	err := self.storage.StoreOrder(orders, timepoint)
	if err != nil {
		log.Printf("Storing orders failed: %s\n", err)
	}
}

func (self *Fetcher) fetchAllBalances(w *sync.WaitGroup, timepoint uint64) {
	defer w.Done()
	data, err := self.blockchain.FetchBalanceData(self.rmaddr, timepoint)
	if err != nil {
		log.Printf("Fetching data from blockchain failed: %s\n", err)
	}
	err = self.storage.StoreBalance(data, timepoint)
	// fmt.Printf("balance data: %v\n", data)
	if err != nil {
		log.Printf("Storing balance data failed: %s\n", err)
	}
}

func (self *Fetcher) fetchAllRates(w *sync.WaitGroup, timepoint uint64) {
	defer w.Done()
	log.Printf("Fetching all rates from blockchain...")
	sources := []common.Token{}
	dests := []common.Token{}
	pairs := map[common.TokenPairID]bool{}
	for _, ex := range self.exchanges {
		tokenPairs := ex.TokenPairs()
		for _, p := range tokenPairs {
			_, exist := pairs[p.PairID()]
			if !exist {
				pairs[p.PairID()] = true
				sources = append(sources, p.Base)
				dests = append(dests, p.Quote)
			}
		}
	}
	data, err := self.blockchain.FetchRates(sources, dests, timepoint)
	if err != nil {
		log.Printf("Fetching data from blockchain failed: %s\n", err)
	}
	err = self.storage.StoreRate(data, timepoint)
	// fmt.Printf("balance data: %v\n", data)
	if err != nil {
		log.Printf("Storing balance data failed: %s\n", err)
	}
}

func (self *Fetcher) fetchAllFromExchanges(timepoint uint64) {
	log.Printf("Fetching all data from exchanges...")
	wait := sync.WaitGroup{}
	wait.Add(1)
	go self.fetchAllPrices(&wait, timepoint)
	wait.Add(1)
	go self.fetchAllEBalances(&wait, timepoint)
	// wait.Add(1)
	// go self.fetchAllOrders(&wait, timepoint)
	log.Printf("Waiting price, balance, order data from exchanges...")
	wait.Wait()
}

func (self *Fetcher) fetchAllFromBlockchain(timepoint uint64) {
	log.Printf("Fetching data from blockchain...")
	wait := sync.WaitGroup{}
	wait.Add(1)
	self.fetchAllBalances(&wait, timepoint)
	wait.Add(1)
	self.fetchAllRates(&wait, timepoint)
	wait.Wait()
}
