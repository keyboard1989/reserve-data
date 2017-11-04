package storage

import (
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
)

type RamStorage struct {
	mu      sync.RWMutex
	version int64
	records []common.ActivityRecord
}

func NewRamStorage() *RamStorage {
	return &RamStorage{
		sync.RWMutex{}, 0, []common.ActivityRecord{},
	}
}

func (self *RamStorage) Record(action string, params map[string]interface{}, result interface{}) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	timestamp := common.GetTimestamp()
	id := self.version + 1
	self.version = id
	self.records = append(self.records, common.ActivityRecord{
		ID:        id,
		Timestamp: timestamp,
		Action:    action,
		Params:    params,
		Result:    result,
	})
	return nil
}

func (self *RamStorage) GetAllRecords() ([]common.ActivityRecord, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.records, nil
}
