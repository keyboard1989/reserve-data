package fetcher

import (
	"github.com/KyberNetwork/reserve-data/common"
	"sync"
)

type ConcurrentAllPriceData struct {
	mu   sync.RWMutex
	data map[common.TokenPairID]common.OnePrice
}

func NewConcurrentAllPriceData() *ConcurrentAllPriceData {
	return &ConcurrentAllPriceData{
		mu:   sync.RWMutex{},
		data: map[common.TokenPairID]common.OnePrice{},
	}
}

func (self *ConcurrentAllPriceData) SetOnePrice(
	exchange common.ExchangeID,
	pair common.TokenPairID,
	d common.ExchangePrice) {
	self.mu.Lock()
	defer self.mu.Unlock()
	_, exist := self.data[pair]
	if !exist {
		self.data[pair] = common.OnePrice{}
	}
	self.data[pair][exchange] = d
}

func (self *ConcurrentAllPriceData) GetData() map[common.TokenPairID]common.OnePrice {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.data
}
