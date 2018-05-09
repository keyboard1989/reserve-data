package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/boltdb/bolt"
)

const (
	TRADE_HISTORY                = "trade_history"
	MAX_GET_TRADE_HISTORY uint64 = 3 * 86400000
)

type ExchangeStorage struct {
	mu sync.RWMutex
	db *bolt.DB
}

func NewBoltExchangeStorage(path string) (*ExchangeStorage, error) {
	// init instance
	var err error
	var db *bolt.DB
	db, err = bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	// init buckets
	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(TRADE_HISTORY))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	storage := &ExchangeStorage{sync.RWMutex{}, db}
	return storage, nil
}

func uint64ToBytes(u uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, u)
	return b
}

func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func (self *ExchangeStorage) StoreTradeHistory(data common.AllTradeHistory) error {
	var err error
	err = self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TRADE_HISTORY))
		for exchange, dataHistory := range data.Data {
			exchangeBk, err := b.CreateBucketIfNotExists([]byte(exchange))
			if err != nil {
				log.Printf("Cannot create exchange history bucket: %s", err.Error())
			}
			for pair, pairHistory := range dataHistory {
				pairBk, err := exchangeBk.CreateBucketIfNotExists([]byte(pair))
				if err != nil {
					log.Printf("Cannot create pair history bucket: %s", err.Error())
				}
				for _, history := range pairHistory {
					idBytes := uint64ToBytes(history.Timestamp)
					dataJSON, err := json.Marshal(history)
					if err != nil {
						log.Printf("Cannot marshal history: %s", err.Error())
					}
					pairBk.Put(idBytes, dataJSON)
				}
			}
		}
		return err
	})
	return err
}

func (self *ExchangeStorage) GetTradeHistory(fromTime, toTime uint64) (common.AllTradeHistory, error) {
	result := common.AllTradeHistory{
		Timestamp: common.GetTimestamp(),
		Data:      map[common.ExchangeID]common.ExchangeTradeHistory{},
	}
	var err error
	if toTime-fromTime > MAX_GET_TRADE_HISTORY {
		return result, errors.New(fmt.Sprintf("Time range is too broad, it must be smaller or equal to 3 days (miliseconds)"))
	}
	min := uint64ToBytes(fromTime)
	max := uint64ToBytes(toTime)
	err = self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TRADE_HISTORY))
		c := b.Cursor()
		for k, v := c.First(); k != nil && v == nil; k, v = c.Next() {
			exchangeBk := b.Bucket(k)
			cursor := exchangeBk.Cursor()
			exchangeHistory := common.ExchangeTradeHistory{}
			for key, value := cursor.First(); key != nil && value == nil; key, value = cursor.Next() {
				pairBk := exchangeBk.Bucket(key)
				pairsHistory := []common.TradeHistory{}
				pairCursor := pairBk.Cursor()
				for pairKey, history := pairCursor.Seek(min); pairKey != nil && bytes.Compare(pairKey, max) <= 0; pairKey, history = pairCursor.Next() {
					pairHistory := common.TradeHistory{}
					json.Unmarshal(history, &pairHistory)
					pairsHistory = append(pairsHistory, pairHistory)
				}
				exchangeHistory[common.TokenPairID(key)] = pairsHistory
			}
			result.Data[common.ExchangeID(k)] = exchangeHistory
		}
		return nil
	})
	return result, err
}

func (self *ExchangeStorage) GetLastIDTradeHistory(exchange, pair string) (string, error) {
	history := common.TradeHistory{}
	err := self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TRADE_HISTORY))
		exchangeBk, err := b.CreateBucketIfNotExists([]byte(exchange))
		if err != nil {
			log.Printf("Cannot get exchange bucket: %s", err.Error())
			return err
		}
		pairBk, err := exchangeBk.CreateBucketIfNotExists([]byte(pair))
		if err != nil {
			log.Printf("Cannot get pair bucket: %s", err.Error())
			return err
		}
		k, v := pairBk.Cursor().Last()
		if k != nil {
			json.Unmarshal(v, &history)
		}
		return err
	})
	return history.ID, err
}
