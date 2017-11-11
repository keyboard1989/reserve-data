package storage

import (
	"errors"
	"github.com/KyberNetwork/reserve-data/common"
	"sync"
)

type RamBalanceStorage struct {
	mu      sync.RWMutex
	version int64
	data    map[int64]map[string]common.BalanceEntry
}

func NewRamBalanceStorage() *RamBalanceStorage {
	return &RamBalanceStorage{
		mu:      sync.RWMutex{},
		version: 0,
		data:    map[int64]map[string]common.BalanceEntry{},
	}
}

func (self *RamBalanceStorage) CurrentVersion(timepoint uint64) (int64, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.version, nil
}

func (self *RamBalanceStorage) GetAllBalances(version int64) (map[string]common.BalanceEntry, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	all := self.data[version]
	if all == nil {
		return map[string]common.BalanceEntry{}, errors.New("Version doesn't exist")
	} else {
		return all, nil
	}
}

func (self *RamBalanceStorage) StoreNewData(data map[string]common.BalanceEntry, timepoint uint64) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.version = self.version + 1
	self.data[self.version] = data
	delete(self.data, self.version-1)
	return nil
}
