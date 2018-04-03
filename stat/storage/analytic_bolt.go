package storage

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

const (
	PRICE_ANALYTIC_BUCKET string = "price_analytic"
)

type BoltAnalyticStorage struct {
	db *bolt.DB
}

func NewBoltAnalyticStorage(path string) (*BoltAnalyticStorage, error) {
	var err error
	var db *bolt.DB
	db, err = bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte(PRICE_ANALYTIC_BUCKET))
		return nil
	})
	storage := &BoltAnalyticStorage{db}
	return storage, nil
}

func (self *BoltAnalyticStorage) UpdatePriceAnalyticData(timestamp uint64, value []byte) error {
	var err error
	k := uint64ToBytes(timestamp)
	self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PRICE_ANALYTIC_BUCKET))
		c := b.Cursor()
		existedKey, _ := c.Seek(k)
		if existedKey != nil {
			err = errors.New("The timestamp is already existed.")
			return err
		}
		err = b.Put(k, value)
		return err
	})
	return err
}

func (self *BoltAnalyticStorage) GetPriceAnalyticData(fromTime uint64, toTime uint64) (map[uint64][]byte, error) {
	var err error
	min := uint64ToBytes(fromTime)
	max := uint64ToBytes(toTime)
	result := make(map[uint64][]byte)
	if toTime-fromTime > MAX_GET_LOG_PERIOD {
		return result, errors.New(fmt.Sprintf("Time range is too broad, it must be smaller or equal to %d miliseconds", MAX_GET_RATES_PERIOD))
	}

	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PRICE_ANALYTIC_BUCKET))
		c := b.Cursor()
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			timestamp := bytesToUint64(k)
			result[timestamp] = v
		}
		return nil
	})
	return result, err
}
