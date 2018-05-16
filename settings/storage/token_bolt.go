package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/boltdb/bolt"
)

const (
	ACTIVE_TOKEN_BUCKET_BY_ID   string = "active_token_by_id"
	ACTIVE_TOKEN_BUCKET_BY_ADDR string = "active_token_by_addr"

	ALL_TOKEN_BUCKET_BY_ID   string = "all_token_by_id"
	ALL_TOKEN_BUCKET_BY_ADDR string = "all_token_by_addr"

	INTERNAL_ACTIVE_BUCKET_BY_ID   string = "internal_token_by_id"
	INTERNAL_ACTIVE_BUCKET_BY_ADDR string = "internal_token_by_addr"

	EXTERNAL_ACTIVE_BUCKET_BY_ID   string = "external_token_by_id"
	EXTERNAL_ACTIVE_BUCKET_BY_ADDR string = "external_token_by_addr"
)

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
		_, err = tx.CreateBucketIfNotExists([]byte(ACTIVE_TOKEN_BUCKET_BY_ID))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(ACTIVE_TOKEN_BUCKET_BY_ADDR))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(ALL_TOKEN_BUCKET_BY_ID))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(ALL_TOKEN_BUCKET_BY_ADDR))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(INTERNAL_ACTIVE_BUCKET_BY_ID))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(INTERNAL_ACTIVE_BUCKET_BY_ADDR))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(EXTERNAL_ACTIVE_BUCKET_BY_ID))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(EXTERNAL_ACTIVE_BUCKET_BY_ADDR))
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

func (self *BoltTokenStorage) AddActiveTokenByID(t common.Token) error {
	return self.setTokenToBucketByID(t, ACTIVE_TOKEN_BUCKET_BY_ID)
}

func (self *BoltTokenStorage) AddTokenByID(t common.Token) error {
	return self.setTokenToBucketByID(t, ALL_TOKEN_BUCKET_BY_ID)
}

func (self *BoltTokenStorage) AddInternalActiveTokenByID(t common.Token) error {
	return self.setTokenToBucketByID(t, INTERNAL_ACTIVE_BUCKET_BY_ID)
}

func (self *BoltTokenStorage) AddExternalActiveTokenByID(t common.Token) error {
	return self.setTokenToBucketByID(t, EXTERNAL_ACTIVE_BUCKET_BY_ID)
}

func (self *BoltTokenStorage) setTokenToBucketByID(t common.Token, bucketName string) error {
	err := self.db.Update(func(tx *bolt.Tx) error {
		b, uErr := tx.CreateBucketIfNotExists([]byte(bucketName))
		if uErr != nil {
			return uErr
		}
		dataJson, uErr := json.Marshal(t)
		return b.Put([]byte(strings.ToLower(t.ID)), dataJson)
	})
	return err
}

func (self *BoltTokenStorage) AddActiveTokenByAddr(t common.Token) error {
	return self.setTokenToBucketByAddr(t, ACTIVE_TOKEN_BUCKET_BY_ADDR)
}

func (self *BoltTokenStorage) AddTokenByAddr(t common.Token) error {
	return self.setTokenToBucketByAddr(t, ALL_TOKEN_BUCKET_BY_ADDR)
}

func (self *BoltTokenStorage) AddInternalTokenByAddr(t common.Token) error {
	return self.setTokenToBucketByAddr(t, INTERNAL_ACTIVE_BUCKET_BY_ADDR)
}

func (self *BoltTokenStorage) AddExternalActiveTokenByAddr(t common.Token) error {
	return self.setTokenToBucketByAddr(t, EXTERNAL_ACTIVE_BUCKET_BY_ADDR)
}

func (self *BoltTokenStorage) setTokenToBucketByAddr(t common.Token, bucketName string) error {
	err := self.db.Update(func(tx *bolt.Tx) error {
		b, uErr := tx.CreateBucketIfNotExists([]byte(bucketName))
		if uErr != nil {
			return uErr
		}
		dataJson, uErr := json.Marshal(t)
		return b.Put([]byte(strings.ToLower(t.Address)), dataJson)
	})
	return err
}

func (self *BoltTokenStorage) GetActiveTokens() ([]common.Token, error) {
	return self.getAllTokenFromBucket(ACTIVE_TOKEN_BUCKET_BY_ID)
}

func (self *BoltTokenStorage) GetAllTokens() ([]common.Token, error) {
	return self.getAllTokenFromBucket(ALL_TOKEN_BUCKET_BY_ID)
}

func (self *BoltTokenStorage) GetInternalTokens() ([]common.Token, error) {
	return self.getAllTokenFromBucket(INTERNAL_ACTIVE_BUCKET_BY_ID)
}

func (self *BoltTokenStorage) GetExternalTokens() ([]common.Token, error) {
	return self.getAllTokenFromBucket(EXTERNAL_ACTIVE_BUCKET_BY_ID)
}

func (self *BoltTokenStorage) getAllTokenFromBucket(bucketName string) (result []common.Token, err error) {
	err = self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
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
			result = append(result, temp)
		}
		return nil
	})
	return result, err
}

func (self *BoltTokenStorage) GetActiveTokenbyID(ID string) (common.Token, error) {
	return self.getATokenFromBucketByID(ID, ACTIVE_TOKEN_BUCKET_BY_ID)
}

func (self *BoltTokenStorage) GetTokenByID(ID string) (common.Token, error) {
	return self.getATokenFromBucketByID(ID, ALL_TOKEN_BUCKET_BY_ID)
}

func (self *BoltTokenStorage) GetInternalTokenByID(ID string) (common.Token, error) {
	return self.getATokenFromBucketByID(ID, INTERNAL_ACTIVE_BUCKET_BY_ID)
}

func (self *BoltTokenStorage) GetExternalTokenByID(ID string) (common.Token, error) {
	return self.getATokenFromBucketByID(ID, EXTERNAL_ACTIVE_BUCKET_BY_ID)
}

func (self *BoltTokenStorage) getATokenFromBucketByID(ID, bucketName string) (common.Token, error) {
	var t common.Token
	ID = strings.ToLower(ID)
	err := self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("Bucket doesn't exist yet")
		}
		c := b.Cursor()
		k, v := c.Seek([]byte(ID))
		if bytes.Compare(k, []byte(ID)) != 0 {
			return fmt.Errorf("Token %d is not found in current setting", ID)
		}
		return json.Unmarshal(v, &t)
	})
	return t, err
}

func (self *BoltTokenStorage) GetActiveTokenbyAddress(Addr string) (common.Token, error) {
	return self.getATokenFromBucketByAddress(Addr, ACTIVE_TOKEN_BUCKET_BY_ADDR)
}

func (self *BoltTokenStorage) GetTokenByAddress(Addr string) (common.Token, error) {
	return self.getATokenFromBucketByAddress(Addr, ALL_TOKEN_BUCKET_BY_ADDR)
}

func (self *BoltTokenStorage) GetInternalTokenByAddress(Addr string) (common.Token, error) {
	return self.getATokenFromBucketByAddress(Addr, INTERNAL_ACTIVE_BUCKET_BY_ADDR)
}

func (self *BoltTokenStorage) GetExternalTokenByAddress(Addr string) (common.Token, error) {
	return self.getATokenFromBucketByAddress(Addr, EXTERNAL_ACTIVE_BUCKET_BY_ADDR)
}

func (self *BoltTokenStorage) getATokenFromBucketByAddress(Addr, bucketName string) (common.Token, error) {
	var t common.Token
	Addr = strings.ToLower(Addr)
	err := self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("Bucket doesn't exist yet")
		}
		c := b.Cursor()
		k, v := c.Seek([]byte(Addr))
		if bytes.Compare(k, []byte(Addr)) != 0 {
			return fmt.Errorf("Token %d is not found in current setting", Addr)
		}
		return json.Unmarshal(v, &t)
	})
	return t, err
}
