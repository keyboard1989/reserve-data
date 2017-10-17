package exchanges

import (
	"github.com/KyberNetwork/reserve-data/market"
	"sync"
)

type ConcurrentResultMap struct {
	mu   sync.RWMutex
	data map[string]*market.ExchangeData
}

func NewConcurrentResultMap() *ConcurrentResultMap {
	return &ConcurrentResultMap{
		sync.RWMutex{},
		map[string]*market.ExchangeData{},
	}
}

func (self *ConcurrentResultMap) Set(key string, value *market.ExchangeData) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.data[key] = value
}

func (self *ConcurrentResultMap) GetData() map[string]*market.ExchangeData {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.data
}
