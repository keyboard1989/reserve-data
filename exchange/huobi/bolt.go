package huobi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/boltdb/bolt"
)

const (
	INTERMEDIATE_TX         string = "intermediate_tx"
	PENDING_INTERMEDIATE_TX string = "pending_intermediate_tx"
	TRADE_HISTORY           string = "trade_history"
	//MAX_GET_TRADE_HISTORY is 3 days
	MAX_GET_TRADE_HISTORY uint64 = 3 * 86400000
)

//BoltStorage object
type BoltStorage struct {
	mu sync.RWMutex
	db *bolt.DB
}

//NewBoltStorage return new storage instance
func NewBoltStorage(path string) (*BoltStorage, error) {
	// init instance
	var err error
	var db *bolt.DB
	db, err = bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	// init buckets
	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucket([]byte(INTERMEDIATE_TX)); err != nil {
			log.Printf("Create bucket error: %s", err.Error())
			return err
		}
		if _, err := tx.CreateBucket([]byte(PENDING_INTERMEDIATE_TX)); err != nil {
			log.Printf("Create bucket error: %s", err.Error())
			return err
		}
		return nil
	})
	storage := &BoltStorage{sync.RWMutex{}, db}
	return storage, nil
}

//GetPendingIntermediateTXs return pending transaction for first deposit phase
func (self *BoltStorage) GetPendingIntermediateTXs() (map[common.ActivityID]common.TXEntry, error) {
	result := make(map[common.ActivityID]common.TXEntry)
	var err error
	err = self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PENDING_INTERMEDIATE_TX))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			actID := common.ActivityID{}
			record := common.TXEntry{}
			if err = json.Unmarshal(k, &actID); err != nil {
				return err
			}
			if err = json.Unmarshal(v, &record); err != nil {
				return err
			}
			result[actID] = record
		}
		return nil
	})
	return result, err
}

//StorePendingIntermediateTx store pending transaction
func (self *BoltStorage) StorePendingIntermediateTx(id common.ActivityID, data common.TXEntry) error {
	var err error
	err = self.db.Update(func(tx *bolt.Tx) error {
		var dataJSON []byte
		b := tx.Bucket([]byte(PENDING_INTERMEDIATE_TX))
		dataJSON, err = json.Marshal(data)
		if err != nil {
			return err
		}
		idJSON, err := json.Marshal(id)
		if err != nil {
			return err
		}
		return b.Put(idJSON, dataJSON)
	})
	return err
}

//RemovePendingIntermediateTx remove pending transaction
func (self *BoltStorage) RemovePendingIntermediateTx(id common.ActivityID) error {
	var err error
	err = self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PENDING_INTERMEDIATE_TX))
		idJSON, err := json.Marshal(id)
		if err != nil {
			return err
		}
		err = b.Delete(idJSON)
		return err
	})
	return err
}

//StoreIntermediateTx store transaction
func (self *BoltStorage) StoreIntermediateTx(id common.ActivityID, data common.TXEntry) error {
	var err error
	err = self.db.Update(func(tx *bolt.Tx) error {
		var dataJSON []byte
		b := tx.Bucket([]byte(INTERMEDIATE_TX))
		dataJSON, err = json.Marshal(data)
		if err != nil {
			return err
		}
		idByte := id.ToBytes()
		return b.Put(idByte[:], dataJSON)
	})
	return err
}

func isTheSame(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for ix, v := range a {
		if b[ix] != v {
			return false
		}
	}
	return true
}

//GetIntermedatorTx get intermediate transaction
func (self *BoltStorage) GetIntermedatorTx(id common.ActivityID) (common.TXEntry, error) {
	var tx2 common.TXEntry
	var err error
	err = self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(INTERMEDIATE_TX))
		c := b.Cursor()
		idBytes := id.ToBytes()
		k, v := c.Seek(idBytes[:])
		if isTheSame(k, idBytes[:]) {
			return json.Unmarshal(v, &tx2)
		}
		return errors.New("Can not find 2nd transaction tx for the deposit, please try later")
	})
	return tx2, err
}

//StoreTradeHistory store trade history
func (self *BoltStorage) StoreTradeHistory(data common.ExchangeTradeHistory) error {
	var err error
	err = self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TRADE_HISTORY))
		for pair, pairHistory := range data {
			pairBk, err := b.CreateBucketIfNotExists([]byte(pair))
			if err != nil {
				log.Printf("Cannot create pair history bucket: %s", err.Error())
				return err
			}
			for _, history := range pairHistory {
				idBytes := []byte(fmt.Sprintf("%s%s", strconv.FormatUint(history.Timestamp, 10), history.ID))
				dataJSON, err := json.Marshal(history)
				if err != nil {
					log.Printf("Cannot marshal history: %s", err.Error())
					return err
				}
				err = pairBk.Put(idBytes, dataJSON)
				if err != nil {
					log.Printf("Cannot put the new data: %s", err.Error())
					return err
				}
			}
		}
		return err
	})
	return err
}

//GetTradeHistory get trade history
func (self *BoltStorage) GetTradeHistory(fromTime, toTime uint64) (common.ExchangeTradeHistory, error) {
	result := common.ExchangeTradeHistory{}
	var err error
	if toTime-fromTime > MAX_GET_TRADE_HISTORY {
		return result, errors.New("Time range is too broad, it must be smaller or equal to 3 days (miliseconds)")
	}
	min := []byte(strconv.FormatUint(fromTime, 10))
	max := []byte(strconv.FormatUint(toTime, 10))
	err = self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TRADE_HISTORY))
		c := b.Cursor()
		exchangeHistory := common.ExchangeTradeHistory{}
		for key, value := c.First(); key != nil && value == nil; key, value = c.Next() {
			pairBk := b.Bucket(key)
			pairsHistory := []common.TradeHistory{}
			pairCursor := pairBk.Cursor()
			for pairKey, history := pairCursor.Seek(min); pairKey != nil && bytes.Compare(pairKey, max) <= 0; pairKey, history = pairCursor.Next() {
				pairHistory := common.TradeHistory{}
				err = json.Unmarshal(history, &pairHistory)
				if err != nil {
					log.Printf("Cannot unmarshal history: %s", err.Error())
					return err
				}
				pairsHistory = append(pairsHistory, pairHistory)
			}
			exchangeHistory[common.TokenPairID(key)] = pairsHistory
		}
		result = exchangeHistory
		return nil
	})
	return result, err
}

//GetLastIDTradeHistory get last trade history id
func (self *BoltStorage) GetLastIDTradeHistory(exchange, pair string) (string, error) {
	history := common.TradeHistory{}
	err := self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TRADE_HISTORY))
		pairBk, err := b.CreateBucketIfNotExists([]byte(pair))
		if err != nil {
			log.Printf("Cannot get pair bucket: %s", err.Error())
			return err
		}
		k, v := pairBk.Cursor().Last()
		if k != nil {
			err = json.Unmarshal(v, &history)
			if err != nil {
				log.Printf("Cannot unmarshal history: %s", err.Error())
				return err
			}
		}
		return err
	})
	return history.ID, err
}
