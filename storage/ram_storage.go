package storage

import (
	"errors"
	"fmt"
	"github.com/KyberNetwork/reserve-data/market"
	"strings"
	"sync"
)

type RamStorage struct {
	mu        sync.RWMutex
	version   int64
	timestamp map[int64]string
	data      map[int64]market.AllPriceData
}

func NewRamStorage() *RamStorage {
	return &RamStorage{
		mu:        sync.RWMutex{},
		version:   0,
		timestamp: map[int64]string{},
		data:      map[int64]market.AllPriceData{},
	}
}

func (self *RamStorage) GetVersion() int64 {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.version
}

func (self *RamStorage) GetTimestamp(version int64) string {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.timestamp[version]
}

func (self *RamStorage) GetAllPairData(version int64) (market.AllPriceData, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	all := self.data[version]
	if all.Timestamp == "" {
		return market.AllPriceData{}, errors.New("Version doesn't exist")
	} else {
		return all, nil
	}
}

func (self *RamStorage) GetOnePairData(base, quote string, version int64) (market.OnePairData, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	pairStr := fmt.Sprintf("%s-%s", strings.ToUpper(base), strings.ToUpper(quote))
	all := self.data[version]
	if all.Timestamp == "" {
		return market.OnePairData{}, errors.New("Version doesn't exist")
	} else {
		data := all.AllPairData[pairStr]
		if len(data) == 0 {
			return market.OnePairData{}, errors.New("Pair of token is not supported")
		} else {
			return data, nil
		}
	}
}

func (self *RamStorage) StoreNewData(timestamp string, data market.AllPriceData) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.version = self.version + 1
	self.timestamp[self.version] = timestamp
	delete(self.timestamp, self.version-1)
	self.data[self.version] = data
	delete(self.data, self.version-1)
	return nil
}
