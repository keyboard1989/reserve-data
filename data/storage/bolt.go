package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

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
	METRIC_TARGET_QUANTITY  string = "target_quantity"
	PENDING_TARGET_QUANTITY string = "pending_target_quantity"
	LOG_BUCKET              string = "logs"
	TRADE_HISTORY           string = "trade_history"
	ENABLE_REBALANCE        string = "enable_rebalance"
	SETRATE_CONTROL         string = "setrate_control"
	MAX_NUMBER_VERSION      int    = 1000
	MAX_GET_RATES_PERIOD    uint64 = 86400000 //1 days in milisec

	TRADE_STATS_BUCKET   string = "trade_stats"
	ASSETS_VOLUME_BUCKET string = "assets_volume"
	BURN_FEE_BUCKET      string = "burn_fee"
	WALLET_FEE_BUCKET    string = "wallet_fee"
	USER_VOLUME_BUCKET   string = "user_volume"
	MINUTE_BUCKET        string = "minute"
	HOUR_BUCKET          string = "hour"
	DAY_BUCKET           string = "day"
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
		tx.CreateBucket([]byte(PRICE_BUCKET))
		tx.CreateBucket([]byte(RATE_BUCKET))
		tx.CreateBucket([]byte(ORDER_BUCKET))
		tx.CreateBucket([]byte(ACTIVITY_BUCKET))
		tx.CreateBucket([]byte(PENDING_ACTIVITY_BUCKET))
		tx.CreateBucket([]byte(BITTREX_DEPOSIT_HISTORY))
		tx.CreateBucket([]byte(AUTH_DATA_BUCKET))
		tx.CreateBucket([]byte(METRIC_BUCKET))
		tx.CreateBucket([]byte(METRIC_TARGET_QUANTITY))
		tx.CreateBucket([]byte(PENDING_TARGET_QUANTITY))
		tx.CreateBucket([]byte(LOG_BUCKET))
		tx.CreateBucket([]byte(TRADE_HISTORY))
		tx.CreateBucket([]byte(ENABLE_REBALANCE))
		tx.CreateBucket([]byte(SETRATE_CONTROL))
		tx.CreateBucket([]byte(TRADE_STATS_BUCKET))

		tradeStatsBk := tx.Bucket([]byte(TRADE_STATS_BUCKET))
		metrics := []string{ASSETS_VOLUME_BUCKET, BURN_FEE_BUCKET, WALLET_FEE_BUCKET, USER_VOLUME_BUCKET}
		frequencies := []string{MINUTE_BUCKET, HOUR_BUCKET, DAY_BUCKET}

		for _, metric := range metrics {
			tradeStatsBk.CreateBucket([]byte(metric))
			metricBk := tradeStatsBk.Bucket([]byte(metric))
			for _, freq := range frequencies {
				metricBk.CreateBucket([]byte(freq))
			}
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

func (self *BoltStorage) GetRates(fromTime, toTime uint64) ([]common.AllRateEntry, error) {
	result := []common.AllRateEntry{}
	if toTime-fromTime > MAX_GET_RATES_PERIOD {
		return result, errors.New(fmt.Sprintf("Time range is too broad, it must be smaller or equal to %d miliseconds", MAX_GET_RATES_PERIOD))
	}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RATE_BUCKET))
		c := b.Cursor()
		min := uint64ToBytes(fromTime)
		max := uint64ToBytes(toTime)

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			data := common.AllRateEntry{}
			err = json.Unmarshal(v, &data)
			if err != nil {
				return err
			}
			result = append([]common.AllRateEntry{data}, result...)
		}
		return nil
	})
	return result, err
}

func (self *BoltStorage) GetRate(version common.Version) (common.AllRateEntry, error) {
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
	var lastEntryjson common.AllRateEntry
	self.db.Update(func(tx *bolt.Tx) error {
		var dataJson []byte
		b := tx.Bucket([]byte(RATE_BUCKET))
		c := b.Cursor()
		_, lastEntry := c.Last()
		json.Unmarshal(lastEntry, &lastEntryjson)
		if lastEntryjson.BlockNumber != data.BlockNumber {
			dataJson, err = json.Marshal(data)
			if err != nil {
				return err
			}
			return b.Put(uint64ToBytes(timepoint), dataJson)
		}
		return err
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
			// all other pending set rates should be staled now
			// remove all of them
			// AFTER EXPERIMENT, THIS WILL NOT WORK
			// log.Printf("===> Trying to remove staled set rates")
			// if record.Action == "set_rates" {
			// 	stales := []common.ActivityRecord{}
			// 	c := pb.Cursor()
			// 	for k, v := c.First(); k != nil; k, v = c.Next() {
			// 		record := common.ActivityRecord{}
			// 		log.Printf("===> staled act: %+v", record)
			// 		err = json.Unmarshal(v, &record)
			// 		if err != nil {
			// 			return err
			// 		}
			// 		if record.Action == "set_rates" {
			// 			stales = append(stales, record)
			// 		}
			// 	}
			// 	log.Printf("===> removing staled acts: %+v", stales)
			// 	self.RemoveStalePendingActivities(tx, stales)
			// }
			// after remove all of them, put new set rate activity
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

func getLastPendingSetrate(pendings []common.ActivityRecord, minedNonce uint64) (*common.ActivityRecord, error) {
	var maxNonce uint64 = 0
	var maxPrice uint64 = 0
	var result *common.ActivityRecord
	for _, act := range pendings {
		if act.Action == "set_rates" {
			log.Printf("looking for pending set_rates: %+v", act)
			var nonce uint64
			actNonce := act.Result["nonce"]
			if actNonce != nil {
				nonce, _ = strconv.ParseUint(actNonce.(string), 10, 64)
			} else {
				nonce = 0
			}
			if nonce < minedNonce {
				// this is a stale actitivity, ignore it
				continue
			}
			var gasPrice uint64
			actPrice := act.Result["gasPrice"]
			if actPrice != nil {
				gasPrice, _ = strconv.ParseUint(actPrice.(string), 10, 64)
			} else {
				gasPrice = 0
			}
			if nonce == maxNonce {
				if gasPrice > maxPrice {
					maxNonce = nonce
					result = &act
					maxPrice = gasPrice
				}
			} else if nonce > maxNonce {
				maxNonce = nonce
				result = &act
				maxPrice = gasPrice
			}
		}
	}
	return result, nil
}

func (self *BoltStorage) RemoveStalePendingActivities(tx *bolt.Tx, stales []common.ActivityRecord) error {
	pb := tx.Bucket([]byte(PENDING_ACTIVITY_BUCKET))
	for _, stale := range stales {
		idBytes := stale.ID.ToBytes()
		if err := pb.Delete(idBytes[:]); err != nil {
			return err
		}
	}
	return nil
}

func (self *BoltStorage) PendingSetrate(minedNonce uint64) (*common.ActivityRecord, error) {
	pendings, err := self.GetPendingActivities()
	if err != nil {
		return nil, err
	} else {
		return getLastPendingSetrate(pendings, minedNonce)
	}
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
		// only update when it exists in pending activity bucket because
		// It might be deleted if it is replaced by another activity
		found := pb.Get(idBytes[:])
		if found != nil {
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

func (self *BoltStorage) GetPendingTargetQty() (metric.TokenTargetQty, error) {
	var err error
	var tokenTargetQty metric.TokenTargetQty
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PENDING_TARGET_QUANTITY))
		_, data := b.Cursor().Last()
		if data == nil {
			err = errors.New("There no pending target quantity")
		} else {
			err = json.Unmarshal(data, &tokenTargetQty)
			if err != nil {
				log.Printf("Cannot unmarshal: %s", err.Error())
			}
		}
		return nil
	})
	return tokenTargetQty, err
}

func (self *BoltStorage) StorePendingTargetQty(data, dataType string) error {
	var err error
	timepoint := common.GetTimepoint()
	tokenTargetQty := metric.TokenTargetQty{}
	self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PENDING_TARGET_QUANTITY))
		_, lastPending := b.Cursor().Last()
		if lastPending != nil {
			err = errors.New("There is another pending target quantity. Please confirm or cancel it before setting new target.")
			return err
		} else {
			tokenTargetQty.ID = timepoint
			tokenTargetQty.Status = "unconfirmed"
			tokenTargetQty.Data = data
			tokenTargetQty.Type, _ = strconv.ParseInt(dataType, 10, 64)
			idByte := uint64ToBytes(timepoint)
			var dataJson []byte
			dataJson, err = json.Marshal(tokenTargetQty)
			if err != nil {
				return err
			}
			log.Printf("Target to save: %v", dataJson)
			return b.Put(idByte, dataJson)
		}
		return err
	})
	return err
}

func (self *BoltStorage) RemovePendingTargetQty() error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PENDING_TARGET_QUANTITY))
		k, lastPending := b.Cursor().Last()
		log.Printf("Last key: %s", k)
		if lastPending == nil {
			return errors.New("There is no pending target quantity.")
		} else {
			b.Delete([]byte(k))
			return nil
		}
		return err
	})
	return err
}

func (self *BoltStorage) CurrentTargetQtyVersion(timepoint uint64) (common.Version, error) {
	var result uint64
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(METRIC_TARGET_QUANTITY)).Cursor()
		result, err = reverseSeek(timepoint, c)
		return nil
	})
	return common.Version(result), err
}

func (self *BoltStorage) GetTokenTargetQty() (metric.TokenTargetQty, error) {
	tokenTargetQty := metric.TokenTargetQty{}
	version, err := self.CurrentTargetQtyVersion(common.GetTimepoint())
	log.Printf("Current version: %s", version)
	if err != nil {
		log.Printf("Cannot get version: %s", err.Error())
	}
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(METRIC_TARGET_QUANTITY))
		data := b.Get(uint64ToBytes(uint64(version)))
		if data == nil {
			err = errors.New(fmt.Sprintf("version %s doesn't exist", version))
		} else {
			err = json.Unmarshal(data, &tokenTargetQty)
			if err != nil {
				log.Printf("Cannot unmarshal: %s", err.Error())
			}
		}
		return nil
	})
	return tokenTargetQty, err
}

func (self *BoltStorage) StoreTokenTargetQty(id, data string) error {
	var err error
	var tokenTargetQty metric.TokenTargetQty
	var dataJson []byte
	self.db.Update(func(tx *bolt.Tx) error {
		pending := tx.Bucket([]byte(PENDING_TARGET_QUANTITY))
		_, pendingTargetQty := pending.Cursor().Last()

		if pendingTargetQty == nil {
			err = errors.New("There is no pending target activity to confirm.")
			return err
		} else {
			// verify confirm data
			json.Unmarshal(pendingTargetQty, &tokenTargetQty)
			pendingData := tokenTargetQty.Data
			idInt, _ := strconv.ParseUint(id, 10, 64)
			if tokenTargetQty.ID != idInt {
				err = errors.New("Pending target quantity ID does not match")
				return err
			}
			if data != pendingData {
				err = errors.New("Pending target quantity data does not match")
				return err
			}

			// Save to confirmed target quantity
			tokenTargetQty.Status = "confirmed"
			b := tx.Bucket([]byte(METRIC_TARGET_QUANTITY))
			dataJson, err = json.Marshal(tokenTargetQty)
			if err != nil {
				return err
			}
			idByte := uint64ToBytes(common.GetTimepoint())
			return b.Put(idByte, dataJson)
		}
	})
	if err == nil {
		// remove pending target qty
		return self.RemovePendingTargetQty()
	}
	return err
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

func getBucketNameByFreq(freq string) (bucketName string, err error) {
	switch freq {
	case "m", "M":
		bucketName = MINUTE_BUCKET
	case "h", "H":
		bucketName = HOUR_BUCKET
	case "d", "D":
		bucketName = DAY_BUCKET
	default:
		err = errors.New(fmt.Sprintf("Invalid frequencies"))
	}

	return
}

func getTimestampByFreq(t uint64, freq string) (result []byte) {
	switch freq {
	case "m", "M":
		result = uint64ToBytes(t / uint64(time.Minute) * uint64(time.Minute))
	case "h", "H":
		result = uint64ToBytes(t / uint64(time.Hour) * uint64(time.Hour))
	case "d", "D":
		result = uint64ToBytes(t / uint64(time.Hour*24) * uint64(time.Hour*24))
	}
	return
}

func (self *BoltStorage) SetTradeStats(metric, freq string, t uint64, tradeStats common.TradeStats) (err error) {
	self.db.Update(func(tx *bolt.Tx) error {
		tradeStatsBk := tx.Bucket([]byte(TRADE_STATS_BUCKET))
		metricBk := tradeStatsBk.Bucket([]byte(metric))

		freqBkName, err := getBucketNameByFreq(freq)
		if err != nil {
			return err
		}
		freqBk := metricBk.Bucket([]byte(freqBkName))

		timestamp := getTimestampByFreq(t, freq)
		rawStats := freqBk.Get(timestamp)
		var stats common.TradeStats
		if rawStats != nil {
			json.Unmarshal(rawStats, &stats)
		} else {
			stats = common.TradeStats{}
		}

		for key, value := range tradeStats {

			if v, ok := value.(*big.Int); ok {
				if v == nil {
					continue
				}
			}

			sum, ok := stats[key]
			if ok {
				switch v := value.(type) {
				case float64:
					stats[key] = sum.(float64) + v
				case *big.Int:
					s := new(big.Int)
					s.SetUint64(uint64(sum.(float64)))
					stats[key] = s.Add(s, v)
				}
			} else {
				stats[key] = value
			}
		}

		dataJSON, err := json.Marshal(stats)
		if err != nil {
			return err
		}

		if err := freqBk.Put(timestamp, dataJSON); err != nil {
			return err
		}

		return err
	})
	return
}

func (self *BoltStorage) getTradeStats(fromTime, toTime uint64, freq, metric, key string) (result []common.TradeStats, err error) {
	self.db.View(func(tx *bolt.Tx) error {
		// Get trade stats bucket
		tradeStatsBk := tx.Bucket([]byte(TRADE_STATS_BUCKET))
		metricBk := tradeStatsBk.Bucket([]byte(metric))

		freqBkName, err := getBucketNameByFreq(freq)
		if err != nil {
			return err
		}

		freqBk := metricBk.Bucket([]byte(freqBkName))
		c := freqBk.Cursor()
		min := getTimestampByFreq(fromTime, freq)
		max := getTimestampByFreq(toTime, freq)

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			stats := common.TradeStats{}
			err = json.Unmarshal(v, &stats)
			if err != nil {
				return err
			}

			_, ok := stats[key]
			if ok {
				record := common.TradeStats{
					strconv.FormatUint(bytesToUint64(k), 10): stats[key],
				}
				result = append([]common.TradeStats{record}, result...)
			}
		}

		return err
	})
	return
}

func (self *BoltStorage) GetAssetVolume(fromTime uint64, toTime uint64, freq string, asset string) (result []common.TradeStats, err error) {
	token, err := common.GetToken(asset)
	if err != nil {
		return
	}
	result, err = self.getTradeStats(fromTime, toTime, freq, ASSETS_VOLUME_BUCKET, token.Address)
	return
}

func (self *BoltStorage) GetBurnFee(fromTime uint64, toTime uint64, freq string, reserveAddr string) (result []common.TradeStats, err error) {
	result, err = self.getTradeStats(fromTime, toTime, freq, BURN_FEE_BUCKET, reserveAddr)
	return
}

func (self *BoltStorage) GetWalletFee(fromTime uint64, toTime uint64, freq string, reserveAddr string, walletAddr string) (result []common.TradeStats, err error) {
	key := strings.Join([]string{reserveAddr, walletAddr}, "_")
	result, err = self.getTradeStats(fromTime, toTime, freq, WALLET_FEE_BUCKET, key)
	return
}

func (self *BoltStorage) GetUserVolume(fromTime uint64, toTime uint64, freq string, userAddr string) (result []common.TradeStats, err error) {
	result, err = self.getTradeStats(fromTime, toTime, freq, USER_VOLUME_BUCKET, userAddr)
	return
}

func (self *BoltStorage) GetTradeHistory(timepoint uint64) (common.AllTradeHistory, error) {
	result := common.AllTradeHistory{}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TRADE_HISTORY))
		_, data := b.Cursor().First()
		if data == nil {
			err = errors.New(fmt.Sprintf("There no data before timepoint %s", timepoint))
		} else {
			err = json.Unmarshal(data, &result)
		}
		return nil
	})
	return result, err
}

func (self *BoltStorage) StoreTradeHistory(data common.AllTradeHistory, timepoint uint64) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		var dataJson []byte
		b := tx.Bucket([]byte(TRADE_HISTORY))
		// prune out old data
		c := b.Cursor()
		k, _ := c.First()
		if k != nil {
			b.Delete([]byte(k))
		}

		// add new data
		dataJson, err = json.Marshal(data)
		if err != nil {
			return err
		}
		idByte := uint64ToBytes(timepoint)
		return b.Put(idByte, dataJson)
	})
	return err
}

func (self *BoltStorage) GetRebalanceControl() (metric.RebalanceControl, error) {
	var err error
	var result metric.RebalanceControl
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ENABLE_REBALANCE))
		_, data := b.Cursor().First()
		if data == nil {
			result = metric.RebalanceControl{
				Status: false,
			}
			self.StoreRebalanceControl(false)
		} else {
			json.Unmarshal(data, &result)
		}
		return nil
	})
	return result, err
}

func (self *BoltStorage) StoreRebalanceControl(status bool) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		var dataJson []byte
		b := tx.Bucket([]byte(ENABLE_REBALANCE))
		// prune out old data
		c := b.Cursor()
		k, _ := c.First()
		if k != nil {
			b.Delete([]byte(k))
		}

		// add new data
		data := metric.RebalanceControl{
			Status: status,
		}
		dataJson, err = json.Marshal(data)
		if err != nil {
			return err
		}
		idByte := uint64ToBytes(common.GetTimepoint())
		return b.Put(idByte, dataJson)
	})
	return err
}

func (self *BoltStorage) GetSetrateControl() (metric.SetrateControl, error) {
	var err error
	var result metric.SetrateControl
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(SETRATE_CONTROL))
		_, data := b.Cursor().First()
		if data == nil {
			result = metric.SetrateControl{
				Status: false,
			}
			self.StoreSetrateControl(false)
		} else {
			json.Unmarshal(data, &result)
		}
		return nil
	})
	return result, err
}

func (self *BoltStorage) StoreSetrateControl(status bool) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		var dataJson []byte
		b := tx.Bucket([]byte(SETRATE_CONTROL))
		// prune out old data
		c := b.Cursor()
		k, _ := c.First()
		if k != nil {
			b.Delete([]byte(k))
		}

		// add new data
		data := metric.SetrateControl{
			Status: status,
		}
		dataJson, err = json.Marshal(data)
		if err != nil {
			return err
		}
		idByte := uint64ToBytes(common.GetTimepoint())
		return b.Put(idByte, dataJson)
	})
	return err
}
