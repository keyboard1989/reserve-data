package storage

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/boltdb/bolt"
)

const (
	PRICE_BUCKET            string = "prices"
	RATE_BUCKET             string = "rates"
	ORDER_BUCKET            string = "orders"
	ACTIVITY_BUCKET         string = "activities"
	AUTH_DATA_BUCKET        string = "auth_data"
	PENDING_ACTIVITY_BUCKET string = "pending_activities"
	BITTREX_DEPOSIT_HISTORY string = "bittrex_deposit_history"
)

type BoltStorage struct {
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
		_, err = tx.CreateBucket([]byte(PRICE_BUCKET))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket([]byte(RATE_BUCKET))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket([]byte(ORDER_BUCKET))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket([]byte(ACTIVITY_BUCKET))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket([]byte(PENDING_ACTIVITY_BUCKET))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket([]byte(BITTREX_DEPOSIT_HISTORY))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket([]byte(AUTH_DATA_BUCKET))
		if err != nil {
			return err
		}
		return nil
	})
	return &BoltStorage{db}, nil
}

func uint64ToBytes(u uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, u)
	return b
}

func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func reverseSeek(timepoint uint64, c *bolt.Cursor) (uint64, error) {
	version, _ := c.Seek(uint64ToBytes(timepoint))
	if version == nil {
		version, _ = c.Prev()
		if version == nil {
			return 0, errors.New(fmt.Sprintf("There is no data before timepoint %d", timepoint))
		} else {
			return bytesToUint64(version), nil
		}
	} else {
		v := bytesToUint64(version)
		if v == timepoint {
			return v, nil
		} else {
			version, _ = c.Prev()
			if version == nil {
				return 0, errors.New(fmt.Sprintf("There is no data before timepoint %d", timepoint))
			} else {
				return bytesToUint64(version), nil
			}
		}
	}
}

func (self *BoltStorage) CurrentPriceVersion(timepoint uint64) (common.Version, error) {
	var result uint64
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(PRICE_BUCKET)).Cursor()
		result, err = reverseSeek(timepoint, c)
		return nil
	})
	return common.Version(result), err
}

func (self *BoltStorage) GetAllPrices(version common.Version) (map[common.TokenPairID]common.OnePrice, error) {
	result := map[common.TokenPairID]common.OnePrice{}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PRICE_BUCKET))
		data := b.Get(uint64ToBytes(uint64(version)))
		if data == nil {
			err = errors.New(fmt.Sprintf("version %s doesn't exist", version))
		} else {
			err = json.Unmarshal(data, &result)
		}
		return nil
	})
	return result, err
}

func (self *BoltStorage) GetOnePrice(pair common.TokenPairID, version common.Version) (common.OnePrice, error) {
	result := map[common.TokenPairID]common.OnePrice{}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PRICE_BUCKET))
		data := b.Get(uint64ToBytes(uint64(version)))
		if data == nil {
			err = errors.New(fmt.Sprintf("version %s doesn't exist", version))
		} else {
			err = json.Unmarshal(data, &result)
		}
		return nil
	})
	if err != nil {
		return common.OnePrice{}, err
	} else {
		pair, exist := result[pair]
		if exist {
			return pair, nil
		} else {
			return common.OnePrice{}, errors.New("Pair of token is not supported")
		}
	}
}

func (self *BoltStorage) CurrentAuthDataVersion(timepoint uint64) (common.Version, error) {
	var result uint64
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(AUTH_DATA_BUCKET)).Cursor()
		result, err = reverseSeek(timepoint, c)
		return nil
	})
	return common.Version(result), err
}

func (self *BoltStorage) GetAuthData(version common.Version) (common.AuthDataSnapshot, error) {
	result := common.AuthDataSnapshot{}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(AUTH_DATA_BUCKET))
		data := b.Get(uint64ToBytes(uint64(version)))
		if data == nil {
			err = errors.New(fmt.Sprintf("version %s doesn't exist", version))
		} else {
			err = json.Unmarshal(data, &result)
		}
		return nil
	})
	return result, err
}

func (self *BoltStorage) CurrentRateVersion(timepoint uint64) (common.Version, error) {
	var result uint64
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(RATE_BUCKET)).Cursor()
		result, err = reverseSeek(timepoint, c)
		return nil
	})
	return common.Version(result), err
}

func (self *BoltStorage) GetAllRates(version common.Version) (common.AllRateEntry, error) {
	result := common.AllRateEntry{}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RATE_BUCKET))
		data := b.Get(uint64ToBytes(uint64(version)))
		if data == nil {
			err = errors.New(fmt.Sprintf("version %s doesn't exist", version))
		} else {
			err = json.Unmarshal(data, &result)
		}
		return nil
	})
	return result, err
}

func (self *BoltStorage) StorePrice(data map[common.TokenPairID]common.OnePrice, timepoint uint64) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		var dataJson []byte
		b := tx.Bucket([]byte(PRICE_BUCKET))
		dataJson, err = json.Marshal(data)
		if err != nil {
			return err
		}
		return b.Put(uint64ToBytes(timepoint), dataJson)
	})
	return err
}

func (self *BoltStorage) StoreAuthSnapshot(
	data *common.AuthDataSnapshot, timepoint uint64) error {

	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		var dataJson []byte
		b := tx.Bucket([]byte(AUTH_DATA_BUCKET))
		dataJson, err = json.Marshal(data)
		if err != nil {
			return err
		}
		return b.Put(uint64ToBytes(timepoint), dataJson)
	})
	return err
}

func (self *BoltStorage) StoreRate(data common.AllRateEntry, timepoint uint64) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		var dataJson []byte
		b := tx.Bucket([]byte(RATE_BUCKET))
		dataJson, err = json.Marshal(data)
		if err != nil {
			return err
		}
		return b.Put(uint64ToBytes(timepoint), dataJson)
	})
	return err
}

func (self *BoltStorage) Record(
	action string,
	id common.ActivityID,
	destination string,
	params map[string]interface{}, result map[string]interface{},
	estatus string,
	mstatus string,
	timepoint uint64) error {

	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		var dataJson []byte
		b := tx.Bucket([]byte(ACTIVITY_BUCKET))
		record := common.ActivityRecord{
			Action:         action,
			ID:             id,
			Destination:    destination,
			Params:         params,
			Result:         result,
			ExchangeStatus: estatus,
			MiningStatus:   mstatus,
			Timestamp:      common.Timestamp(strconv.FormatUint(timepoint, 10)),
		}
		dataJson, err = json.Marshal(record)
		if err != nil {
			return err
		}
		idByte, _ := id.MarshalText()
		err = b.Put(idByte, dataJson)
		if err != nil {
			return err
		}
		if record.IsPending() {
			pb := tx.Bucket([]byte(PENDING_ACTIVITY_BUCKET))
			err = pb.Put(idByte, dataJson)
		}
		return err
	})
	return err
}

func (self *BoltStorage) GetAllRecords() ([]common.ActivityRecord, error) {
	result := []common.ActivityRecord{}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ACTIVITY_BUCKET))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			record := common.ActivityRecord{}
			err = json.Unmarshal(v, &record)
			if err != nil {
				return err
			}
			result = append(result, record)
		}
		return nil
	})
	return result, err
}

func (self *BoltStorage) GetPendingActivities() ([]common.ActivityRecord, error) {
	result := []common.ActivityRecord{}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PENDING_ACTIVITY_BUCKET))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			record := common.ActivityRecord{}
			err = json.Unmarshal(v, &record)
			if err != nil {
				return err
			}
			result = append(result, record)
		}
		return nil
	})
	return result, err
}

func (self *BoltStorage) UpdateActivity(id common.ActivityID, activity common.ActivityRecord) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		pb := tx.Bucket([]byte(PENDING_ACTIVITY_BUCKET))
		idBytes, _ := id.MarshalText()
		dataJson, err := json.Marshal(activity)
		if err != nil {
			return err
		}
		err = pb.Put(idBytes, dataJson)
		if err != nil {
			return err
		}
		if !activity.IsPending() {
			err = pb.Delete(idBytes)
			if err != nil {
				return err
			}
		}
		b := tx.Bucket([]byte(ACTIVITY_BUCKET))
		if err != nil {
			return err
		}
		return b.Put(idBytes, dataJson)
	})
	return err
}

func (self *BoltStorage) IsNewBittrexDeposit(id uint64, actID common.ActivityID) bool {
	res := true
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BITTREX_DEPOSIT_HISTORY))
		v := b.Get(uint64ToBytes(id))
		if v != nil && string(v) != actID.String() {
			res = false
		}
		return nil
	})
	return res
}

func (self *BoltStorage) RegisterBittrexDeposit(id uint64, actID common.ActivityID) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BITTREX_DEPOSIT_HISTORY))
		actIDBytes, _ := actID.MarshalText()
		err = b.Put(uint64ToBytes(id), actIDBytes)
		return nil
	})
	return err
}
