package storage

import (
	"errors"
	"github.com/KyberNetwork/reserve-data/common"
	"sync"
)

type RamOrderStorage struct {
	mu      sync.RWMutex
	version int64
	data    map[int64]common.AllOrderEntry
}

func NewRamOrderStorage() *RamOrderStorage {
	return &RamOrderStorage{
		mu:      sync.RWMutex{},
		version: 0,
		data:    map[int64]common.AllOrderEntry{},
	}
}

func (self *RamOrderStorage) CurrentVersion(timepoint uint64) (int64, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.version, nil
}

func (self *RamOrderStorage) GetOrders(version int64) (common.AllOrderEntry, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	all, exist := self.data[version]
	if !exist {
		return common.AllOrderEntry{}, errors.New("Version doesn't exist")
	} else {
		return all, nil
	}
}

func (self *RamOrderStorage) StoreNewData(data common.AllOrderEntry, timepoint uint64) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.version = self.version + 1
	self.data[self.version] = data
	delete(self.data, self.version-1)
	return nil
}
