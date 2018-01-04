package fetcher

import (
	"github.com/KyberNetwork/reserve-data/common"
	"sync"
)

type ConcurrentAllPriceData struct {
	mu   sync.RWMutex
	data common.AllPriceEntry
}

func NewConcurrentAllPriceData() *ConcurrentAllPriceData {
	return &ConcurrentAllPriceData{
		mu: sync.RWMutex{},
		data: common.AllPriceEntry{
			Data:  map[common.TokenPairID]common.OnePrice{},
			Block: 0,
		},
	}
}

func (self *ConcurrentAllPriceData) SetBlockNumber(block uint64) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.data.Block = block
}

func (self *ConcurrentAllPriceData) SetOnePrice(
	exchange common.ExchangeID,
	pair common.TokenPairID,
	d common.ExchangePrice) {
	self.mu.Lock()
	defer self.mu.Unlock()
	_, exist := self.data.Data[pair]
	if !exist {
		self.data.Data[pair] = common.OnePrice{}
	}
	self.data.Data[pair][exchange] = d
}

func (self *ConcurrentAllPriceData) GetData() common.AllPriceEntry {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.data
}
