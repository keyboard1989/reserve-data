package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/boltdb/bolt"
)

const (
	MAX_GET_LOG_PERIOD uint64 = 86400000 //1 days in milisec

	TRADELOG_BUCKET string = "logs"
	CATLOG_BUCKET   string = "cat_logs"
)

type BoltLogStorage struct {
	mu    sync.RWMutex
	db    *bolt.DB
	block uint64
}

func NewBoltLogStorage(path string) (*BoltLogStorage, error) {
	// init instance
	var err error
	var db *bolt.DB
	db, err = bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	// init buckets
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte(TRADELOG_BUCKET))
		tx.CreateBucket([]byte(CATLOG_BUCKET))
		return nil
	})
	storage := &BoltLogStorage{sync.RWMutex{}, db, 0}
	storage.db.View(func(tx *bolt.Tx) error {
		block, _, err := storage.LoadLastLogIndex(tx)
		if err == nil {
			storage.block = block
		}
		return err
	})
	return storage, nil
}

func (self *BoltLogStorage) LoadLastCatLog(tx *bolt.Tx) (common.SetCatLog, error) {
	b := tx.Bucket([]byte(CATLOG_BUCKET))
	c := b.Cursor()
	k, v := c.Last()
	record := common.SetCatLog{}
	if k != nil {
		err := json.Unmarshal(v, &record)
		if err != nil {
			return record, err
		}
		return record, nil
	} else {
		return record, errors.New("Database is empty")
	}
}

func (self *BoltLogStorage) LoadLastLogIndex(tx *bolt.Tx) (uint64, uint, error) {
	b := tx.Bucket([]byte(TRADELOG_BUCKET))
	c := b.Cursor()
	k, v := c.Last()
	if k != nil {
		record := common.TradeLog{}
		json.Unmarshal(v, &record)
		return record.BlockNumber, record.TransactionIndex, nil
	} else {
		return 0, 0, errors.New("Database is empty")
	}
}

func (self *BoltLogStorage) StoreCatLog(l common.SetCatLog) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CATLOG_BUCKET))
		var dataJson []byte
		lastCat, berr := self.LoadLastCatLog(tx)
		if berr == nil && (lastCat.Address.Big().Cmp(l.Address.Big()) == 0) {
			err = errors.New(
				fmt.Sprintf("Duplicated log %+v", l))
			return err
		}
		dataJson, err = json.Marshal(l)
		if err != nil {
			return err
		}
		log.Printf("Storing cat log: %d", l.Timestamp)
		idByte := uint64ToBytes(l.Timestamp)
		err = b.Put(idByte, dataJson)
		return err
	})
	return err
}

func (self *BoltLogStorage) StoreTradeLog(stat common.TradeLog, timepoint uint64) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TRADELOG_BUCKET))
		var dataJson []byte
		block, txindex, berr := self.LoadLastLogIndex(tx)
		if berr == nil && (block > stat.BlockNumber || (block == stat.BlockNumber && txindex >= stat.TransactionIndex)) {
			err = errors.New(
				fmt.Sprintf("Duplicated log (new block number %s is smaller or equal to latest block number %s)", block, stat.BlockNumber))
			return err
		}
		dataJson, err = json.Marshal(stat)
		if err != nil {
			return err
		}
		log.Printf("Storing log: %d", stat.Timestamp)
		idByte := uint64ToBytes(stat.Timestamp)
		err = b.Put(idByte, dataJson)
		return err
	})
	return err
}

func (self *BoltLogStorage) UpdateLogBlock(block uint64, timepoint uint64) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.block = block
	return nil
}

func (self *BoltLogStorage) GetCatLogs(fromTime uint64, toTime uint64) ([]common.SetCatLog, error) {
	result := []common.SetCatLog{}
	var err error
	if toTime-fromTime > MAX_GET_LOG_PERIOD {
		return result, errors.New(fmt.Sprintf("Time range is too broad, it must be smaller or equal to %d miliseconds", MAX_GET_RATES_PERIOD))
	}
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CATLOG_BUCKET))
		c := b.Cursor()
		min := uint64ToBytes(fromTime * 1000000)
		max := uint64ToBytes(toTime * 1000000)
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			record := common.SetCatLog{}
			err = json.Unmarshal(v, &record)
			if err != nil {
				return err
			}
			result = append([]common.SetCatLog{record}, result...)
		}
		return nil
	})
	return result, err
}

func (self *BoltLogStorage) GetTradeLogs(fromTime uint64, toTime uint64) ([]common.TradeLog, error) {
	result := []common.TradeLog{}
	var err error
	if toTime-fromTime > MAX_GET_LOG_PERIOD {
		return result, errors.New(fmt.Sprintf("Time range is too broad, it must be smaller or equal to %d miliseconds", MAX_GET_RATES_PERIOD))
	}
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TRADELOG_BUCKET))
		c := b.Cursor()
		min := uint64ToBytes(fromTime * 1000000)
		max := uint64ToBytes(toTime * 1000000)
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			record := common.TradeLog{}
			err = json.Unmarshal(v, &record)
			if err != nil {
				return err
			}
			result = append([]common.TradeLog{record}, result...)
		}
		return nil
	})
	return result, err
}
func (self *BoltLogStorage) LastBlock() (uint64, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.block, nil
}
