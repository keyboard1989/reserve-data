package fetcher

import (
	"sort"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
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

func (self *ConcurrentAllPriceData) UpdatePrice(
	oldPrice []common.PriceEntry,
	newPrice []common.PriceEntry) []common.PriceEntry {
	for _, price := range newPrice {
		if price.Quantity == 0 {
			// find the rate
			i := sort.Search(len(oldPrice), func(i int) bool { return oldPrice[i].Rate == price.Rate })
			// if exist, remove it
			if i < len(oldPrice) {
				oldPrice = append(oldPrice[:i], oldPrice[i+1:]...)
			}
		} else {
			// insert to the right place
			i := sort.Search(len(oldPrice), func(i int) bool { return oldPrice[i].Rate >= price.Rate })
			if i < len(oldPrice) && oldPrice[i].Rate == price.Rate {
				oldPrice[i] = price
			} else {
				oldPrice = append(oldPrice[:i], append([]common.PriceEntry{price}, oldPrice[i:]...)...)
			}
		}
	}
	return oldPrice
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

func (self *ConcurrentAllPriceData) UpdateOnePrice(
	exchange common.ExchangeID,
	pair common.TokenPairID,
	d common.ExchangePrice) {
	self.mu.Lock()
	defer self.mu.Unlock()
	_, exist := self.data[pair]
	if !exist {
		self.data[pair] = common.OnePrice{}
	}
	exchangePrice := self.data[pair][exchange]
	exchangePrice.Valid = d.Valid
	exchangePrice.Timestamp = d.Timestamp
	exchangePrice.ReturnTime = d.ReturnTime
	exchangePrice.Error = d.Error
	exchangePrice.Bids = self.UpdatePrice(exchangePrice.Bids, d.Bids)
	exchangePrice.Asks = self.UpdatePrice(exchangePrice.Asks, d.Asks)
	self.data[pair][exchange] = exchangePrice
}

func (self *ConcurrentAllPriceData) GetData() map[common.TokenPairID]common.OnePrice {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.data
}
