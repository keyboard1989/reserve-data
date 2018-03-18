package storage

import (
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
	INTERMEDIATE_TX      string = "intermediate_tx"
	MAX_NUMBER_VERSION   int    = 1000
	MAX_GET_RATES_PERIOD uint64 = 86400000 //1 days in milisec
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

// GetNumberOfVersion return number of version storing in a bucket
func (self *BoltStorage) GetNumberOfVersion(tx *bolt.Tx, bucket string) int {
	result := 0
	b := tx.Bucket([]byte(bucket))
	c := b.Cursor()
	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		result++
	}
	return result
}

// PruneOutdatedData Remove first version out of database
func (self *BoltStorage) PruneOutdatedData(tx *bolt.Tx, bucket string) error {
	var err error
	b := tx.Bucket([]byte(bucket))
	c := b.Cursor()
	for self.GetNumberOfVersion(tx, bucket) >= MAX_NUMBER_VERSION {
		k, _ := c.First()
		if k == nil {
			err = errors.New(fmt.Sprintf("There no version in %s", bucket))
			return err
		}
		err = b.Delete([]byte(k))
		if err != nil {
			panic(err)
		}
	}
	return err
}
func (self *BoltStorage) StoreIntermediateTx(hash string, exchangeID string, tokenID string, miningStatus string, exchangeStatus string, Amount float64, Timestamp common.Timestamp, id common.ActivityID) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		var dataJson []byte
		b := tx.Bucket([]byte(INTERMEDIATE_TX))

		log.Printf("Version number: %d\n", self.GetNumberOfVersion(tx, INTERMEDIATE_TX))
		self.PruneOutdatedData(tx, INTERMEDIATE_TX)
		log.Printf("After prune number version: %d\n", self.GetNumberOfVersion(tx, INTERMEDIATE_TX))
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
