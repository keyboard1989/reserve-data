package storage

import (
	"encoding/json"
	"fmt"

	"github.com/KyberNetwork/reserve-data/boltutil"
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/settings"
	"github.com/boltdb/bolt"
)

// GetFee returns a map[tokenID]exchangeFees and error if occur
func (boltSettingStorage *BoltSettingStorage) GetFee(ex settings.ExchangeName) (result common.ExchangeFees, err error) {
	err = boltSettingStorage.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(EXCHANGE_FEE_BUCKET))
		if b == nil {
			return fmt.Errorf("bucket %s hasn't existed yet", EXCHANGE_FEE_BUCKET)
		}
		data := b.Get(boltutil.Uint64ToBytes(uint64(ex)))
		if data == nil {
			return fmt.Errorf("key %s hasn't existed yet", string(ex))
		}
		uErr := json.Unmarshal(data, &result)
		if uErr != nil {
			return uErr
		}
		return nil
	})
	return result, err
}

// StoreFee stores the fee with exchangeName as key into database and return error if occur
func (boltSettingStorage *BoltSettingStorage) StoreFee(ex settings.ExchangeName, data common.ExchangeFees) error {
	err := boltSettingStorage.db.Update(func(tx *bolt.Tx) error {
		b, uErr := tx.CreateBucketIfNotExists([]byte(EXCHANGE_FEE_BUCKET))
		if uErr != nil {
			return uErr
		}
		dataJSON, uErr := json.Marshal(data)
		if uErr != nil {
			return uErr
		}
		return b.Put(boltutil.Uint64ToBytes(uint64(ex)), dataJSON)
	})
	return err
}

// GetMinDeposit returns a map[tokenID]MinDeposit and error if occur
func (boltSettingStorage *BoltSettingStorage) GetMinDeposit(ex settings.ExchangeName) (result common.ExchangesMinDeposit, err error) {
	err = boltSettingStorage.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(EXCHANGE_MIN_DEPOSIT_BUCKET))
		if b == nil {
			return fmt.Errorf("bucket %s hasn't existed yet", EXCHANGE_MIN_DEPOSIT_BUCKET)
		}
		data := b.Get(boltutil.Uint64ToBytes(uint64(ex)))
		if data == nil {
			return fmt.Errorf("key %s hasn't existed yet", string(ex))
		}
		uErr := json.Unmarshal(data, &result)
		if uErr != nil {
			return uErr
		}
		return nil
	})
	return result, err
}

// StoreMinDeposit stores the minDeposit with exchangeName as key into database and return error if occur
func (boltSettingStorage *BoltSettingStorage) StoreMinDeposit(ex settings.ExchangeName, data common.ExchangesMinDeposit) error {
	err := boltSettingStorage.db.Update(func(tx *bolt.Tx) error {
		b, uErr := tx.CreateBucketIfNotExists([]byte(EXCHANGE_MIN_DEPOSIT_BUCKET))
		if uErr != nil {
			return uErr
		}
		dataJSON, uErr := json.Marshal(data)
		if uErr != nil {
			return uErr
		}
		return b.Put(boltutil.Uint64ToBytes(uint64(ex)), dataJSON)
	})
	return err
}

// GetDepositAddress returns a map[tokenID]DepositAddress and error if occur
func (boltSettingStorage *BoltSettingStorage) GetDepositAddress(ex settings.ExchangeName) (result common.ExchangeAddresses, err error) {
	err = boltSettingStorage.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(EXCHANGE_DEPOSIT_ADDRESS))
		if b == nil {
			return fmt.Errorf("bucket %s hasn't existed yet", EXCHANGE_DEPOSIT_ADDRESS)
		}
		data := b.Get(boltutil.Uint64ToBytes(uint64(ex)))
		if data == nil {
			return fmt.Errorf("key %s hasn't existed yet", string(ex))
		}
		uErr := json.Unmarshal(data, &result)
		if uErr != nil {
			return uErr
		}
		return nil
	})
	return result, err

}

// StoreDepositAddress stores the depositAddress with exchangeName as key into database and
// return error if occur
func (boltSettingStorage *BoltSettingStorage) StoreDepositAddress(ex settings.ExchangeName, addrs common.ExchangeAddresses) error {
	err := boltSettingStorage.db.Update(func(tx *bolt.Tx) error {
		b, uErr := tx.CreateBucketIfNotExists([]byte(EXCHANGE_DEPOSIT_ADDRESS))
		if uErr != nil {
			return uErr
		}
		dataJSON, uErr := json.Marshal(addrs)
		if uErr != nil {
			return uErr
		}
		return b.Put(boltutil.Uint64ToBytes(uint64(ex)), dataJSON)
	})
	return err
}

// Count Fee return the number of element in fee database, and error if occcur.
// It is used mainly to check if the fee database is empty
func (boltSettingStorage *BoltSettingStorage) CountFee() (uint64, error) {
	var result uint64 = 0
	err := boltSettingStorage.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(EXCHANGE_FEE_BUCKET))
		if b == nil {
			return fmt.Errorf("bucket %s hasn't existed yet", EXCHANGE_DEPOSIT_ADDRESS)
		}
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			result++
		}
		return nil
	})
	return result, err
}
