package storage

import (
	"errors"
	"github.com/KyberNetwork/reserve-data/common"
	"sync"
)

type RamBalanceStorage struct {
	mu      sync.RWMutex
	version int64
	data    map[int64]map[string]common.RawBalance
}

func NewRamBalanceStorage() *RamBalanceStorage {
	return &RamBalanceStorage{
		mu:      sync.RWMutex{},
		version: 0,
		data:    map[int64]map[string]common.RawBalance{},
	}
}

func (self *RamBalanceStorage) CurrentVersion() (int64, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.version, nil
}

func (self *RamBalanceStorage) GetAllBalances(version int64) (map[string]common.RawBalance, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	all := self.data[version]
	if all == nil {
		return map[string]common.RawBalance{}, errors.New("Version doesn't exist")
	} else {
		return all, nil
	}
}

func (self *RamBalanceStorage) StoreNewData(data map[string]common.RawBalance) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.version = self.version + 1
	self.data[self.version] = data
	delete(self.data, self.version-1)
	return nil
}
