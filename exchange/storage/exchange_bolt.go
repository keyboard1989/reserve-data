package storage

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/boltdb/bolt"
)

const (
	INTERMEDIATE_TX         string = "intermediate_tx"
	PENDING_INTERMEDIATE_TX string = "pending_intermediate_tx"
	MAX_NUMBER_VERSION      int    = 1000
	PENDING_TX2_EXPIRED     uint64 = 3600000 // 1 hour in milisecond
)

type BoltStorage struct {
	mu sync.RWMutex
	db *bolt.DB
}

func NewBoltStorage(path string) (*BoltStorage, error) {
	// init instance
	var err error
	var db *bolt.DB
	db, err = bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	// init buckets
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte(INTERMEDIATE_TX))
		tx.CreateBucket([]byte(PENDING_INTERMEDIATE_TX))
		return nil
	})
	storage := &BoltStorage{sync.RWMutex{}, db}
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

func (self *BoltStorage) StorePendingIntermediateTxID(Timestamp uint64, id common.ActivityID) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PENDING_INTERMEDIATE_TX))
		data := uint64ToBytes(Timestamp)
		var idJson []byte
		idJson, err = json.Marshal(id)
		if err != nil {
			return err
		}
		return b.Put(idJson, data)
	})
	return err
}

func (self *BoltStorage) GetPendingIntermediateTXID(TimeStamp uint64) ([]common.ActivityID, error) {
	result := []common.ActivityID{}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PENDING_INTERMEDIATE_TX))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			record := common.ActivityID{}
			atTime := bytesToUint64(v)
			err = json.Unmarshal(k, &record)
			if err != nil {
				log.Printf("Can not unmarshall ID, check Database design., %v", err)
				return err
			}
			if (TimeStamp - atTime) > PENDING_TX2_EXPIRED {
				log.Printf("Activity %s expired, remove from pending tx2 tracking record", record.EID)
				if err := b.Delete(k); err != nil {
					return err
				}
			} else {
				result = append([]common.ActivityID{record}, result...)
			}

		}
		return nil
	})
	return result, err
}

func (self *BoltStorage) StoreIntermediateTx(hash string, exchangeID string, tokenID string, miningStatus string, exchangeStatus string, Amount float64, Timestamp common.Timestamp, id common.ActivityID) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		var dataJson []byte
		b := tx.Bucket([]byte(INTERMEDIATE_TX))
		data := common.TXEntry{
			hash, exchangeID, tokenID, miningStatus, exchangeStatus, Amount, Timestamp,
		}
		dataJson, err = json.Marshal(data)
		if err != nil {
			return err
		}
		idByte := id.ToBytes()
		return b.Put(idByte[:], dataJson)
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

func (self *BoltStorage) GetIntermedatorTx(id common.ActivityID) (common.TXEntry, error) {
	var tx2 common.TXEntry
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(INTERMEDIATE_TX))
		c := b.Cursor()
		idBytes := id.ToBytes()
		k, v := c.Seek(idBytes[:])
		if isTheSame(k, idBytes[:]) {
			err = json.Unmarshal(v, &tx2)

			if err != nil {
				return err
			}
			return nil
		} else {
			err = errors.New("Can not find 2nd transaction tx for the deposit %s, please try later")
			return err
		}
	})
	return tx2, err
}
