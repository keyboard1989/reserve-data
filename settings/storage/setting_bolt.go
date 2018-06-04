package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/KyberNetwork/reserve-data/boltutil"
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/settings"
	"github.com/boltdb/bolt"
	ethereum "github.com/ethereum/go-ethereum/common"
)

const (
	TOKEN_BUCKET_BY_ID         string = "token_by_id"
	TOKEN_BUCKET_BY_ADDRESS    string = "token_by_addr"
	ADDRESS_SETTING_BUCKET     string = "address_setting"
	ADDRESS_SET_SETTING_BUCKET string = "address_set_setting"
)

type FilterFunction func(common.Token) bool

func isActive(t common.Token) bool {
	return t.Active
}

func isToken(_ common.Token) bool {
	return true
}

func isInternal(t common.Token) bool {
	return t.Active && t.Internal
}

func isExternal(t common.Token) bool {
	return t.Active && !t.Internal
}

type BoltSettingStorage struct {
	db *bolt.DB
}

func NewBoltSettingStorage(dbPath string) (*BoltSettingStorage, error) {
	var err error
	var db *bolt.DB
	db, err = bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		if _, uErr := tx.CreateBucketIfNotExists([]byte(TOKEN_BUCKET_BY_ID)); err != nil {
			return uErr
		}
		if _, uErr := tx.CreateBucketIfNotExists([]byte(TOKEN_BUCKET_BY_ADDRESS)); err != nil {
			return uErr
		}
		if _, uErr := tx.CreateBucketIfNotExists([]byte(ADDRESS_SETTING_BUCKET)); err != nil {
			return uErr
		}
		if _, uErr := tx.CreateBucketIfNotExists([]byte(ADDRESS_SET_SETTING_BUCKET)); err != nil {
			return uErr
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	storage := BoltSettingStorage{db}
	return &storage, nil
}

func addTokenByID(tx *bolt.Tx, t common.Token) error {
	b, uErr := tx.CreateBucketIfNotExists([]byte(TOKEN_BUCKET_BY_ID))
	if uErr != nil {
		return uErr
	}
	dataJSON, uErr := json.Marshal(t)
	if uErr != nil {
		return uErr
	}
	return b.Put([]byte(strings.ToLower(t.ID)), dataJSON)
}

func addTokenByAddress(tx *bolt.Tx, t common.Token) error {
	b, uErr := tx.CreateBucketIfNotExists([]byte(TOKEN_BUCKET_BY_ADDRESS))
	if uErr != nil {
		return uErr
	}
	dataJson, uErr := json.Marshal(t)
	if uErr != nil {
		return uErr
	}
	return b.Put([]byte(strings.ToLower(t.Address)), dataJson)
}

func (self *BoltSettingStorage) UpdateToken(t common.Token) error {
	err := self.db.Update(func(tx *bolt.Tx) error {
		if uErr := addTokenByID(tx, t); uErr != nil {
			return uErr
		}
		if uErr := addTokenByAddress(tx, t); uErr != nil {
			return uErr
		}
		return nil
	})
	return err
}

func (self *BoltSettingStorage) AddTokenByID(t common.Token) error {
	err := self.db.Update(func(tx *bolt.Tx) error {
		b, uErr := tx.CreateBucketIfNotExists([]byte(TOKEN_BUCKET_BY_ID))
		if uErr != nil {
			return uErr
		}
		dataJSON, uErr := json.Marshal(t)
		if uErr != nil {
			return uErr
		}
		return b.Put([]byte(strings.ToLower(t.ID)), dataJSON)
	})
	return err
}

func (self *BoltSettingStorage) AddTokenByAddress(t common.Token) error {
	err := self.db.Update(func(tx *bolt.Tx) error {
		b, uErr := tx.CreateBucketIfNotExists([]byte(TOKEN_BUCKET_BY_ADDRESS))
		if uErr != nil {
			return uErr
		}
		dataJson, uErr := json.Marshal(t)
		if uErr != nil {
			return uErr
		}
		return b.Put([]byte(strings.ToLower(t.Address)), dataJson)
	})
	return err
}

func (self *BoltSettingStorage) getTokensWithFilter(filter FilterFunction) (result []common.Token, err error) {
	err = self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TOKEN_BUCKET_BY_ID))
		if b == nil {
			return fmt.Errorf("Bucket doesn't exist yet")
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var temp common.Token
			uErr := json.Unmarshal(v, &temp)
			if uErr != nil {
				return uErr
			}
			if filter(temp) {
				result = append(result, temp)
			}
		}
		return nil
	})
	return result, err
}

func (self *BoltSettingStorage) GetAllTokens() (result []common.Token, err error) {
	return self.getTokensWithFilter(isToken)
}

func (self *BoltSettingStorage) GetActiveTokens() (result []common.Token, err error) {
	return self.getTokensWithFilter(isActive)
}

func (self *BoltSettingStorage) GetInternalTokens() (result []common.Token, err error) {
	return self.getTokensWithFilter(isInternal)
}

func (self *BoltSettingStorage) GetExternalTokens() (result []common.Token, err error) {
	return self.getTokensWithFilter(isExternal)
}

func (self *BoltSettingStorage) getTokenByIDWithFiltering(id string, filter FilterFunction) (t common.Token, err error) {
	id = strings.ToLower(id)
	err = self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TOKEN_BUCKET_BY_ID))
		if b == nil {
			return fmt.Errorf("Bucket doesn't exist yet")
		}
		c := b.Cursor()
		k, v := c.Seek([]byte(id))
		if bytes.Compare(k, []byte(id)) != 0 {
			return fmt.Errorf("Token %s is not found in current setting", id)
		}
		uErr := json.Unmarshal(v, &t)
		if uErr != nil {
			return uErr
		}
		if !filter(t) {
			return fmt.Errorf("Token %s status is not as querried", id)
		}
		return nil
	})
	return t, err
}

func (self *BoltSettingStorage) GetTokenByID(id string) (common.Token, error) {
	return self.getTokenByIDWithFiltering(id, isToken)
}

func (self *BoltSettingStorage) GetActiveTokenByID(id string) (common.Token, error) {
	return self.getTokenByIDWithFiltering(id, isActive)
}

func (self *BoltSettingStorage) GetInternalTokenByID(id string) (common.Token, error) {
	return self.getTokenByIDWithFiltering(id, isInternal)
}

func (self *BoltSettingStorage) GetExternalTokenByID(id string) (common.Token, error) {
	return self.getTokenByIDWithFiltering(id, isExternal)
}

func (self *BoltSettingStorage) getTokenByAddressWithFiltering(addr string, filter FilterFunction) (t common.Token, err error) {
	addr = strings.ToLower(addr)
	err = self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TOKEN_BUCKET_BY_ADDRESS))
		if b == nil {
			return fmt.Errorf("Bucket doesn't exist yet")
		}
		c := b.Cursor()
		k, v := c.Seek([]byte(addr))
		if bytes.Compare(k, []byte(addr)) != 0 {
			return fmt.Errorf("Token %s is not found in current setting", addr)
		}
		uErr := json.Unmarshal(v, &t)
		if uErr != nil {
			return uErr
		}
		if !filter(t) {
			return fmt.Errorf("Token %s status is not as querried", t.ID)
		}
		return nil
	})
	return t, err
}

func (self *BoltSettingStorage) GetActiveTokenByAddress(Addr ethereum.Address) (common.Token, error) {
	return self.getTokenByAddressWithFiltering(Addr.Hex(), isActive)
}

func (self *BoltSettingStorage) GetTokenByAddress(Addr ethereum.Address) (common.Token, error) {
	return self.getTokenByAddressWithFiltering(Addr.Hex(), isToken)
}

func (self *BoltSettingStorage) GetInternalTokenByAddress(Addr ethereum.Address) (common.Token, error) {
	return self.getTokenByAddressWithFiltering(Addr.Hex(), isInternal)
}

func (self *BoltSettingStorage) GetExternalTokenByAddress(Addr ethereum.Address) (common.Token, error) {
	return self.getTokenByAddressWithFiltering(Addr.Hex(), isExternal)
}

func (self *BoltSettingStorage) UpdateOneAddress(name settings.AddressName, address string) error {
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

func (self *BoltSettingStorage) GetAddress(add settings.AddressName) (string, error) {
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

func (self *BoltSettingStorage) AddAddressToSet(setName settings.AddressSetName, address string) error {
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

func (self *BoltSettingStorage) GetAddresses(setName settings.AddressSetName) ([]string, error) {
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

func (self *BoltSettingStorage) CountAddress() (uint64, error) {
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
