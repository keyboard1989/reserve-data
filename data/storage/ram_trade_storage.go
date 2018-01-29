package storage

import (
	"errors"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
)

type RamTradeStorage struct {
	mu      sync.RWMutex
	version int64
	data    map[int64]common.AllTradeHistory
}

func NewRamTradeStorage() *RamTradeStorage {
	return &RamTradeStorage{
		mu:      sync.RWMutex{},
		version: 0,
		data:    map[int64]common.AllTradeHistory{},
	}
}

func (self *RamTradeStorage) CurrentVersion(timepoint uint64) (int64, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.version, nil
}

func (self *RamTradeStorage) GetTradeHistory(version common.Version) (common.AllTradeHistory, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	all, found := self.data[int64(version)]
	if found {
		return all, nil
	} else {
		return common.AllTradeHistory{}, errors.New("Version doesn't exist")
	}
}

func (self *RamTradeStorage) StoreTradeHistory(data common.AllTradeHistory, timepoint uint64) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.version = self.version + 1
	self.data[self.version] = data
	delete(self.data, self.version-1)
	return nil
}
