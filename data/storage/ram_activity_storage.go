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
	status string,
	timepoint uint64) error {

	self.mu.Lock()
	defer self.mu.Unlock()
	version := self.version + 1
	self.version = version
	record := common.ActivityRecord{
		Action:      action,
		ID:          id,
		Destination: destination,
		Params:      params,
		Result:      result,
		Status:      status,
		Timestamp:   common.Timestamp(strconv.FormatUint(timepoint, 10)),
	}
	self.records.PushBack(&record)
	if status == "submitted" {
		self.pendingRecords.PushBack(&record)
	}
	return nil
}

func (self *RamActivityStorage) UpdateActivityStatus(action string, id common.ActivityID, destination string, status string) error {
	self.mu.RLock()
	defer self.mu.RUnlock()
	ele := self.pendingRecords.Back()
	var updated bool = false
	for {
		if ele == nil {
			break
		} else {
			activity := ele.Value.(*common.ActivityRecord)
			if activity.Action == action && activity.ID == id && activity.Destination == destination {
				updated = true
				activity.Status = status
				self.pendingRecords.Remove(ele)
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
