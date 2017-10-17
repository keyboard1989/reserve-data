package fetcher

import (
	"github.com/KyberNetwork/reserve-data/market"
	"sync"
)

type ConcurrentAllPriceData struct {
	mu   sync.RWMutex
	data *market.AllPriceData
}

func NewConcurrentAllPriceData() *ConcurrentAllPriceData {
	return &ConcurrentAllPriceData{
		mu:   sync.RWMutex{},
		data: market.NewAllPriceData(),
	}
}

func (self *ConcurrentAllPriceData) SetVersion(v int64) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.data.Version = v
}

func (self *ConcurrentAllPriceData) SetTimestamp(t string) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.data.Timestamp = t
}

func (self *ConcurrentAllPriceData) SetOnePairData(exchange Exchange, p string, d *market.ExchangeData) {
	self.mu.Lock()
	defer self.mu.Unlock()
	_, exist := self.data.AllPairData[p]
	if !exist {
		self.data.AllPairData[p] = market.OnePairData{}
	}
	self.data.AllPairData[p][exchange.ID()] = *d
}

func (self *ConcurrentAllPriceData) GetData() market.AllPriceData {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return *self.data
}
