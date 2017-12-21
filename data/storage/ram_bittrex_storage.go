package storage

import (
	"sync"
)

type RamBittrexStorage struct {
	data *sync.Map
}

func (self *RamBittrexStorage) IsNewDeposit(id uint64) bool {
	_, found := self.data.Load(id)
	if found {
		return false
	} else {
		return true
	}
}

func (self *RamBittrexStorage) RegisterDeposit(id uint64) error {
	self.data.Store(id, true)
	return nil
}

func NewRamBittrexStorage() *RamBittrexStorage {
	return &RamBittrexStorage{
		data: &sync.Map{},
	}
}
