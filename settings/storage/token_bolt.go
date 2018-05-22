package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/boltdb/bolt"
)

const (
	TOKEN_BUCKET_BY_ID      string = "token_by_id"
	TOKEN_BUCKET_BY_ADDRESS string = "token_by_addr"
)

type FilterFunction func(common.Token) bool

func isActive(t common.Token) bool {
	return t.Active
}

func isToken(t common.Token) bool {
	return true
}

func isInternal(t common.Token) bool {
	return t.Active && t.Internal
}

func isExternal(t common.Token) bool {
	return t.Active && !t.Internal
}

type BoltTokenStorage struct {
	db *bolt.DB
}

func NewBoltTokenStorage(dbPath string) (*BoltTokenStorage, error) {
	var err error
	var db *bolt.DB
	db, err = bolt.Open(dbPath, 0600, nil)
	if err != nil {
		panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(TOKEN_BUCKET_BY_ID))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(TOKEN_BUCKET_BY_ADDRESS))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	storage := BoltTokenStorage{db}
	return &storage, nil
}

func (self *BoltTokenStorage) AddTokenByID(t common.Token) error {
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

func (self *BoltTokenStorage) AddTokenByAddress(t common.Token) error {
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

func (self *BoltTokenStorage) getTokensWithFilter(filter FilterFunction) (result []common.Token, err error) {
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

func (self *BoltTokenStorage) GetAllTokens() (result []common.Token, err error) {
	return self.getTokensWithFilter(isToken)
}

func (self *BoltTokenStorage) GetActiveTokens() (result []common.Token, err error) {
	return self.getTokensWithFilter(isActive)
}

func (self *BoltTokenStorage) GetInternalTokens() (result []common.Token, err error) {
	return self.getTokensWithFilter(isInternal)
}

func (self *BoltTokenStorage) GetExternalTokens() (result []common.Token, err error) {
	return self.getTokensWithFilter(isExternal)
}

func (self *BoltTokenStorage) getTokenByIDWithFiltering(id string, filter FilterFunction) (t common.Token, err error) {
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

func (self *BoltTokenStorage) GetTokenByID(id string) (common.Token, error) {
	return self.getTokenByIDWithFiltering(id, isToken)
}

func (self *BoltTokenStorage) GetActiveTokenByID(id string) (common.Token, error) {
	return self.getTokenByIDWithFiltering(id, isActive)
}

func (self *BoltTokenStorage) GetInternalTokenByID(id string) (common.Token, error) {
	return self.getTokenByIDWithFiltering(id, isInternal)
}

func (self *BoltTokenStorage) GetExternalTokenByID(id string) (common.Token, error) {
	return self.getTokenByIDWithFiltering(id, isExternal)
}

func (self *BoltTokenStorage) getTokenByAddressWithFiltering(addr string, filter FilterFunction) (t common.Token, err error) {
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

func (self *BoltTokenStorage) GetActiveTokenByAddress(Addr ethereum.Address) (common.Token, error) {
	return self.getTokenByAddressWithFiltering(Addr.Hex(), isActive)
}

func (self *BoltTokenStorage) GetTokenByAddress(Addr ethereum.Address) (common.Token, error) {
	return self.getTokenByAddressWithFiltering(Addr.Hex(), isToken)
}

func (self *BoltTokenStorage) GetInternalTokenByAddress(Addr ethereum.Address) (common.Token, error) {
	return self.getTokenByAddressWithFiltering(Addr.Hex(), isInternal)
}

func (self *BoltTokenStorage) GetExternalTokenByAddress(Addr ethereum.Address) (common.Token, error) {
	return self.getTokenByAddressWithFiltering(Addr.Hex(), isExternal)
}
