package storage

import (
	"errors"
	"github.com/KyberNetwork/reserve-data/common"
	"sync"
)

type RamStorage struct {
	price *RamPriceStorage
}

func NewRamStorage() *RamStorage {
	return &RamStorage{
		NewRamPriceStorage(),
	}
}

func (self *RamStorage) CurrentPriceVersion() (common.Version, error) {
	version, err := self.price.CurrentVersion()
	return common.Version(version), err
}

func (self *RamStorage) GetAllPrices(version common.Version) (map[common.TokenPairID]common.OnePrice, error) {
	return self.price.GetAllPrices(int64(version))
}

func (self *RamStorage) GetOnePrice(pair common.TokenPairID, version common.Version) (common.OnePrice, error) {
	return self.price.GetOnePrice(pair, int64(version))
}

func (self *RamStorage) StorePrice(data map[common.TokenPairID]common.OnePrice) error {
	return self.price.StoreNewData(data)
}

type RamPriceStorage struct {
	mu      sync.RWMutex
	version int64
	data    map[int64]map[common.TokenPairID]common.OnePrice
}

func NewRamPriceStorage() *RamPriceStorage {
	return &RamPriceStorage{
		mu:      sync.RWMutex{},
		version: 0,
		data:    map[int64]map[common.TokenPairID]common.OnePrice{},
	}
}

func (self *RamPriceStorage) CurrentVersion() (int64, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.version, nil
}

func (self *RamPriceStorage) GetAllPrices(version int64) (map[common.TokenPairID]common.OnePrice, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	all := self.data[version]
	if all == nil {
		return map[common.TokenPairID]common.OnePrice{}, errors.New("Version doesn't exist")
	} else {
		return all, nil
	}
}

func (self *RamPriceStorage) GetOnePrice(pair common.TokenPairID, version int64) (common.OnePrice, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	all := self.data[version]
	if all == nil {
		return common.OnePrice{}, errors.New("Version doesn't exist")
	} else {
		data := all[pair]
		if len(data) == 0 {
			return common.OnePrice{}, errors.New("Pair of token is not supported")
		} else {
			return data, nil
		}
	}
}

func (self *RamPriceStorage) StoreNewData(data map[common.TokenPairID]common.OnePrice) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.version = self.version + 1
	self.data[self.version] = data
	delete(self.data, self.version-1)
	return nil
}
