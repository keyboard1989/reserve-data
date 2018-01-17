package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/metric"
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
	METRIC_BUCKET           string = "metrics"
	LOG_BUCKET              string = "logs"
	MAX_NUMBER_VERSION      int    = 1000
)

type BoltStorage struct {
	mu    sync.RWMutex
	db    *bolt.DB
	block uint64
	index uint
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
		_, err = tx.CreateBucket([]byte(METRIC_BUCKET))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket([]byte(LOG_BUCKET))
		if err != nil {
			return err
		}
		return nil
	})
	storage := &BoltStorage{sync.RWMutex{}, db, 0, 0}
	storage.db.View(func(tx *bolt.Tx) error {
		block, index, err := storage.LoadLastLogIndex(tx)
		if err == nil {
			storage.block = block
			storage.index = index
		}
		return err
	})
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

func (self *BoltStorage) GetAllPrices(version common.Version) (common.AllPriceEntry, error) {
	result := common.AllPriceEntry{}
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
	result := common.AllPriceEntry{}
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
		pair, exist := result.Data[pair]
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

func (self *BoltStorage) StorePrice(data common.AllPriceEntry, timepoint uint64) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		var dataJson []byte
		b := tx.Bucket([]byte(PRICE_BUCKET))

		// remove outdated data from bucket
		log.Printf("Version number: %d\n", self.GetNumberOfVersion(tx, PRICE_BUCKET))
		self.PruneOutdatedData(tx, PRICE_BUCKET)
		log.Printf("After prune number version: %d\n", self.GetNumberOfVersion(tx, PRICE_BUCKET))

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
		// idByte, _ := id.MarshalText()
		idByte := id.ToBytes()
		err = b.Put(idByte[:], dataJson)
		if err != nil {
			return err
		}
		if record.IsPending() {
			pb := tx.Bucket([]byte(PENDING_ACTIVITY_BUCKET))
			err = pb.Put(idByte[:], dataJson)
		}
		return err
	})
	return err
}

func formatTimepointToActivityID(timepoint uint64, id []byte) []byte {
	if timepoint == 0 {
		return id
	} else {
		activityID := common.NewActivityID(timepoint, "")
		byteID := activityID.ToBytes()
		return byteID[:]
	}
}

func (self *BoltStorage) GetAllRecords(fromTime, toTime uint64) ([]common.ActivityRecord, error) {
	result := []common.ActivityRecord{}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ACTIVITY_BUCKET))
		c := b.Cursor()
		fkey, _ := c.First()
		lkey, _ := c.Last()
		min := formatTimepointToActivityID(fromTime, fkey)
		max := formatTimepointToActivityID(toTime, lkey)

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			record := common.ActivityRecord{}
			err = json.Unmarshal(v, &record)
			if err != nil {
				return err
			}
			result = append([]common.ActivityRecord{record}, result...)
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
			result = append(
				[]common.ActivityRecord{record}, result...)
		}
		return nil
	})
	return result, err
}

func (self *BoltStorage) UpdateActivity(id common.ActivityID, activity common.ActivityRecord) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		pb := tx.Bucket([]byte(PENDING_ACTIVITY_BUCKET))
		// idBytes, _ := id.MarshalText()
		idBytes := id.ToBytes()
		dataJson, err := json.Marshal(activity)
		if err != nil {
			return err
		}
		err = pb.Put(idBytes[:], dataJson)
		if err != nil {
			return err
		}
		if !activity.IsPending() {
			err = pb.Delete(idBytes[:])
			if err != nil {
				return err
			}
		}
		b := tx.Bucket([]byte(ACTIVITY_BUCKET))
		if err != nil {
			return err
		}
		return b.Put(idBytes[:], dataJson)
	})
	return err
}

func (self *BoltStorage) IsNewBittrexDeposit(id uint64, actID common.ActivityID) bool {
	res := true
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BITTREX_DEPOSIT_HISTORY))
		v := b.Get(uint64ToBytes(id))
		if v != nil && string(v) != actID.String() {
			log.Printf("bolt: stored act id - current act id: %s - %s", string(v), actID.String())
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
		// actIDBytes, _ := actID.MarshalText()
		actIDBytes, _ := actID.MarshalText()
		err = b.Put(uint64ToBytes(id), actIDBytes)
		return nil
	})
	return err
}

func (self *BoltStorage) HasPendingDeposit(token common.Token, exchange common.Exchange) bool {
	result := false
	self.db.View(func(tx *bolt.Tx) error {
		pb := tx.Bucket([]byte(PENDING_ACTIVITY_BUCKET))
		c := pb.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			record := common.ActivityRecord{}
			json.Unmarshal(v, &record)
			if record.Action == "deposit" && record.Params["token"].(string) == token.ID && record.Destination == string(exchange.ID()) {
				result = true
			}
		}
		return nil
	})
	return result
}

func (self *BoltStorage) StoreMetric(data *metric.MetricEntry, timepoint uint64) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		var dataJson []byte
		b := tx.Bucket([]byte(METRIC_BUCKET))
		dataJson, err = json.Marshal(data)
		if err != nil {
			return err
		}
		idByte := uint64ToBytes(data.Timestamp)
		err = b.Put(idByte, dataJson)
		return err
	})
	return err
}

func (self *BoltStorage) GetMetric(tokens []common.Token, fromTime, toTime uint64) (map[string]metric.MetricList, error) {
	imResult := map[string]*metric.MetricList{}
	for _, tok := range tokens {
		imResult[tok.ID] = &metric.MetricList{}
	}

	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(METRIC_BUCKET))
		c := b.Cursor()
		min := uint64ToBytes(fromTime)
		max := uint64ToBytes(toTime)

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			data := metric.MetricEntry{}
			err = json.Unmarshal(v, &data)
			if err != nil {
				return err
			}
			for tok, m := range data.Data {
				metricList, found := imResult[tok]
				if found {
					*metricList = append(*metricList, metric.TokenMetricResponse{
						Timestamp: data.Timestamp,
						AfpMid:    m.AfpMid,
						Spread:    m.Spread,
					})
				}
			}
		}
		return nil
	})
	result := map[string]metric.MetricList{}
	for k, v := range imResult {
		result[k] = *v
	}
	return result, err
}

func (self *BoltStorage) LoadLastLogIndex(tx *bolt.Tx) (uint64, uint, error) {
	b := tx.Bucket([]byte(LOG_BUCKET))
	c := b.Cursor()
	k, v := c.Last()
	if k != nil {
		record := common.TradeLog{}
		json.Unmarshal(v, &record)
		return record.BlockNumber, record.TransactionIndex, nil
	} else {
		return 0, 0, errors.New("Database is empty")
	}
}

func (self *BoltStorage) UpdateLogBlock(block uint64, timepoint uint64) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.block = block
	return nil
}

func (self *BoltStorage) LastBlock() (uint64, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.block, nil
}

func (self *BoltStorage) GetTradeLogs(fromTime uint64, toTime uint64) ([]common.TradeLog, error) {
	result := []common.TradeLog{}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(LOG_BUCKET))
		c := b.Cursor()
		min := uint64ToBytes(fromTime)
		max := uint64ToBytes(toTime)
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			record := common.TradeLog{}
			err = json.Unmarshal(v, &record)
			if err != nil {
				return err
			}
			result = append([]common.TradeLog{record}, result...)
		}
		return nil
	})
	return result, err
}

func (self *BoltStorage) StoreTradeLog(stat common.TradeLog, timepoint uint64) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(LOG_BUCKET))
		var dataJson []byte
		block, txindex, berr := self.LoadLastLogIndex(tx)
		if berr == nil && (block > stat.BlockNumber || (block == stat.BlockNumber && txindex >= stat.TransactionIndex)) {
			err = errors.New(
				fmt.Sprintf("Duplicated log (new block number %s is smaller or equal to latest block number %s)", block, stat.BlockNumber))
			return err
		}
		dataJson, err = json.Marshal(stat)
		if err != nil {
			return err
		}
		log.Printf("Storing log: %d", stat.Timestamp)
		idByte := uint64ToBytes(stat.Timestamp)
		err = b.Put(idByte, dataJson)
		return err
	})
	return err
}
