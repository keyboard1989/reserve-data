package storage

import (
	"errors"
	"github.com/KyberNetwork/reserve-data/common"
	"sync"
)

type RamRateStorage struct {
	mu      sync.RWMutex
	version int64
	data    map[int64]map[common.TokenPairID]common.RateEntry
}

func NewRamRateStorage() *RamRateStorage {
	return &RamRateStorage{
		mu:      sync.RWMutex{},
		version: 0,
		data:    map[int64]map[common.TokenPairID]common.RateEntry{},
	}
}

func (self *RamRateStorage) CurrentVersion() (int64, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.version, nil
}

func (self *RamRateStorage) GetRates(version int64) (map[common.TokenPairID]common.RateEntry, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	all := self.data[version]
	if all == nil {
		return map[common.TokenPairID]common.RateEntry{}, errors.New("Version doesn't exist")
	} else {
		return all, nil
	}
}

func (self *RamRateStorage) StoreNewData(data map[common.TokenPairID]common.RateEntry) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.version = self.version + 1
	self.data[self.version] = data
	delete(self.data, self.version-1)
	return nil
}
