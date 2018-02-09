package storage

import (
	"container/list"
	"errors"
	"sort"
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
	// all other pending set rates should be staled now
	// remove all of them
	// AFTER EXPERIMENT, THIS WILL NOT WORK
	// if action == "set_rates" {
	// 	stales := []common.ActivityRecord{}
	// 	activities := activitiesFromList(self.pendingRecords)
	// 	for _, act := range activities {
	// 		if act.Action == "set_rates" {
	// 			stales = append(stales, act)
	// 		}
	// 	}
	// 	self.RemovePendings(stales)
	// }
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

func (self *RamActivityStorage) GetAllRecords(fromTime, toTime uint64) ([]common.ActivityRecord, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	allRecords := activitiesFromList(self.records)
	toIndex := sort.Search(len(allRecords), func(i int) bool {
		timepoint := allRecords[i].ID.Timepoint
		return timepoint <= fromTime
	})
	fromIndex := sort.Search(len(allRecords), func(i int) bool {
		timepoint := allRecords[i].ID.Timepoint
		return timepoint <= toTime
	})
	from := 0
	to := len(allRecords)
	if toTime != 0 && fromIndex < len(allRecords) {
		from = fromIndex
	}
	if fromTime != 0 && toIndex < len(allRecords) {
		timePoint := allRecords[toIndex].ID.Timepoint
		if timePoint == fromTime {
			to = toIndex + 1
		} else {
			to = toIndex
		}
	}
	result := allRecords[from:to]
	return result, nil
}

func (self *RamActivityStorage) RemovePendings(stales []common.ActivityRecord) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	toRemoves := []*list.Element{}
	ele := self.pendingRecords.Back()
	for {
		if ele == nil {
			break
		} else {
			activity := ele.Value.(*common.ActivityRecord)
			for _, stale := range stales {
				if activity.ID == stale.ID {
					toRemoves = append(toRemoves, ele)
				}
			}
			ele = ele.Prev()
		}
	}
	for _, e := range toRemoves {
		self.pendingRecords.Remove(e)
	}
	return nil
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
			if activity.Action == "deposit" && activity.Params["token"].(common.Token) == token && activity.Destination == string(exchange.ID()) {
				return true
			}
			ele = ele.Prev()
		}
	}
	return false
}
