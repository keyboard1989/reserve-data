package storage

import (
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
)

type RamBittrexStorage struct {
	data *sync.Map
}

func (self *RamBittrexStorage) IsNewDeposit(id uint64, actID common.ActivityID) bool {
	v, found := self.data.Load(id)
	if found && v.(common.ActivityID) != actID {
		return false
	} else {
		return true
	}
}

func (self *RamBittrexStorage) RegisterDeposit(id uint64, actID common.ActivityID) error {
	self.data.Store(id, actID)
	return nil
}

func NewRamBittrexStorage() *RamBittrexStorage {
	return &RamBittrexStorage{
		data: &sync.Map{},
	}
}
