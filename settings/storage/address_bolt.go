package storage

import (
	"errors"
	"fmt"
	"strings"

	"github.com/KyberNetwork/reserve-data/boltutil"
	"github.com/KyberNetwork/reserve-data/settings"
	"github.com/boltdb/bolt"
)

const (
	ADDRESS_SETTING_BUCKET     string = "address_setting"
	ADDRESS_SET_SETTING_BUCKET string = "address_set_setting"
)

type BoltAddressStorage struct {
	db *bolt.DB
}

func NewBoltAddressStorage(dbPath string) (*BoltAddressStorage, error) {
	var err error
	var db *bolt.DB
	db, err = bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		if _, uErr := tx.CreateBucketIfNotExists([]byte(ADDRESS_SETTING_BUCKET)); err != nil {
			return uErr
		}
		if _, uErr := tx.CreateBucketIfNotExists([]byte(ADDRESS_SET_SETTING_BUCKET)); err != nil {
			return uErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &BoltAddressStorage{db}, nil
}

func (self *BoltAddressStorage) UpdateOneAddress(name settings.AddressName, address string) error {
	address = strings.ToLower(address)
	err := self.db.Update(func(tx *bolt.Tx) error {
		b, uErr := tx.CreateBucketIfNotExists([]byte(ADDRESS_SETTING_BUCKET))
		if uErr != nil {
			return uErr
		}
		return b.Put(boltutil.Uint64ToBytes(uint64(name)), []byte(address))
	})
	return err
}

func (self *BoltAddressStorage) GetAddress(add settings.AddressName) (string, error) {
	var address string
	err := self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ADDRESS_SETTING_BUCKET))
		if b == nil {
			return fmt.Errorf("Bucket doesn't exist yet")
		}
		data := b.Get(boltutil.Uint64ToBytes(uint64(add)))
		if data == nil {
			return fmt.Errorf("Key %s is not found", add)
		}
		address = string(data)
		return nil
	})
	return address, err
}

func (self *BoltAddressStorage) AddAddressToSet(setName settings.AddressSetName, address string) error {
	address = strings.ToLower(address)
	defaultValue := "1"
	err := self.db.Update(func(tx *bolt.Tx) error {
		b, uErr := tx.CreateBucketIfNotExists([]byte(ADDRESS_SET_SETTING_BUCKET))
		if uErr != nil {
			return uErr
		}
		s, uErr := b.CreateBucketIfNotExists(boltutil.Uint64ToBytes(uint64(setName)))
		if uErr != nil {
			return uErr
		}
		return s.Put([]byte(address), []byte(defaultValue))
	})
	return err
}

func (self *BoltAddressStorage) GetAddresses(setName settings.AddressSetName) ([]string, error) {
	result := []string{}
	err := self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ADDRESS_SET_SETTING_BUCKET))
		if b == nil {
			return errors.New("Bucket doesn't exist yet")
		}
		s := b.Bucket(boltutil.Uint64ToBytes(uint64(setName)))
		if s == nil {
			return fmt.Errorf("Address set with name %s doesn't exist", setName)
		}
		c := s.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			result = append(result, string(k))
		}
		return nil
	})
	return result, err
}

func CountAddressSetBucket(b *bolt.Bucket) (uint64, error) {
	var result uint64
	c := b.Cursor()
	for setName, _ := c.First(); setName != nil; setName, _ = c.Next() {
		s := b.Bucket(setName)
		if s == nil {
			return 0, fmt.Errorf("bucket %s is not available", string(setName))
		}
		sc := s.Cursor()
		for addr, _ := sc.First(); addr != nil; addr, _ = sc.Next() {
			result += 1
		}
	}
	return result, nil
}

func CountAddressBucket(b *bolt.Bucket) uint64 {
	var result uint64
	c := b.Cursor()
	for name, _ := c.First(); name != nil; name, _ = c.Next() {
		result += 1
	}
	return result
}

func (self *BoltAddressStorage) CountAddress() (uint64, error) {
	var result uint64
	err := self.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(ADDRESS_SET_SETTING_BUCKET))
		if b == nil {
			return fmt.Errorf("bucket %s hasn't existed yet", ADDRESS_SET_SETTING_BUCKET)
		}
		count, uErr := CountAddressSetBucket(b)
		if uErr != nil {
			return uErr
		}
		result += count
		b = tx.Bucket([]byte(ADDRESS_SETTING_BUCKET))
		if b == nil {
			return fmt.Errorf("bucket %s hasn't existed yet", ADDRESS_SETTING_BUCKET)
		}
		result += CountAddressBucket(b)
		return nil
	})
	if err != nil {
		return 0, err
	}
	return result, nil
}
