package storage

import (
	"container/list"
	"errors"
	"strconv"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
)

type RamActivityStorage struct {
	mu             sync.RWMutex
	version        int64
	records        *list.List
	pendingRecords *list.List
}

func NewRamActivityStorage() *RamActivityStorage {
	return &RamActivityStorage{
		sync.RWMutex{}, 0, list.New(), list.New(),
	}
}

func activitiesFromList(l *list.List) []common.ActivityRecord {
	result := []common.ActivityRecord{}
	ele := l.Back()
	for {
		if ele == nil {
			break
		} else {
			result = append(result, *ele.Value.(*common.ActivityRecord))
			ele = ele.Prev()
		}
	}
	return result
}

func (self *RamActivityStorage) StoreNewData(
	action string,
	id common.ActivityID,
	destination string,
	params map[string]interface{}, result map[string]interface{},
	estatus string,
	mstatus string,
	timepoint uint64) error {

	self.mu.Lock()
	defer self.mu.Unlock()
	version := self.version + 1
	self.version = version
	record := common.ActivityRecord{
		Action:         action,
		ID:             id,
		Destination:    destination,
		Params:         params,
		Result:         result,
		ExchangeStatus: estatus,
		MiningStatus:   mstatus,
		Timestamp:      common.Timestamp(strconv.FormatUint(timepoint, 10)),
	}
	self.records.PushBack(&record)
	if record.IsPending() {
		self.pendingRecords.PushBack(&record)
	}
	return nil
}

func (self *RamActivityStorage) UpdateActivity(id common.ActivityID, activity common.ActivityRecord) error {
	self.mu.RLock()
	defer self.mu.RUnlock()
	ele := self.pendingRecords.Back()
	var updated bool = false
	for {
		if ele == nil {
			break
		} else {
			oldAct := ele.Value.(*common.ActivityRecord)
			if oldAct.ID == id {
				updated = true
				oldAct.ExchangeStatus = activity.ExchangeStatus
				oldAct.MiningStatus = activity.MiningStatus
				if !oldAct.IsPending() {
					self.pendingRecords.Remove(ele)
				}
			}
			ele = ele.Prev()
		}
	}
	if updated {
		return nil
	} else {
		return errors.New("Pending activity nout found")
	}
}

func (self *RamActivityStorage) GetAllRecords() ([]common.ActivityRecord, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return activitiesFromList(self.records), nil
}

func (self *RamActivityStorage) GetPendingRecords() ([]common.ActivityRecord, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return activitiesFromList(self.pendingRecords), nil
}

func (self *RamActivityStorage) HasPendingDeposit(token common.Token, exchange common.Exchange) bool {
	self.mu.RLock()
	defer self.mu.RUnlock()
	ele := self.pendingRecords.Back()
	for {
		if ele == nil {
			break
		} else {
			activity := ele.Value.(*common.ActivityRecord)
			if activity.Action == "deposit" && activity.Params["token"].(string) == token.ID && activity.Destination == string(exchange.ID()) {
				return true
			}
			ele = ele.Prev()
		}
	}
	return false
}
