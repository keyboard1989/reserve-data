package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/boltdb/bolt"
)

const (
	MAX_NUMBER_VERSION   int    = 1000
	MAX_GET_RATES_PERIOD uint64 = 86400000 //1 days in milisec

	LOG_BUCKET           string = "logs"
	TRADE_STATS_BUCKET   string = "trade_stats"
	ASSETS_VOLUME_BUCKET string = "assets_volume"
	BURN_FEE_BUCKET      string = "burn_fee"
	WALLET_FEE_BUCKET    string = "wallet_fee"
	USER_VOLUME_BUCKET   string = "user_volume"
	MINUTE_BUCKET        string = "minute"
	HOUR_BUCKET          string = "hour"
	DAY_BUCKET           string = "day"

	ADDRESS_CATEGORY string = "address_category"
	ADDRESS_ID       string = "address_id"
	ID_ADDRESSES     string = "id_addresses"
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
		tx.CreateBucket([]byte(LOG_BUCKET))
		tx.CreateBucket([]byte(ADDRESS_ID))
		tx.CreateBucket([]byte(ID_ADDRESSES))
		tx.CreateBucket([]byte(ADDRESS_CATEGORY))
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

func (self *BoltStorage) GetTradeLogs(fromTime uint64, toTime uint64) ([]common.TradeLog, error) {
	result := []common.TradeLog{}
	var err error
	if toTime-fromTime > MAX_GET_RATES_PERIOD {
		return result, errors.New(fmt.Sprintf("Time range is too broad, it must be smaller or equal to %d miliseconds", MAX_GET_RATES_PERIOD))
	}
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(LOG_BUCKET))
		c := b.Cursor()
		min := uint64ToBytes(fromTime * 1000000)
		max := uint64ToBytes(toTime * 1000000)
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

			sum, ok := stats[key]
			if ok {
				stats[key] = sum + value
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

func (self *BoltStorage) getTradeStats(fromTime, toTime uint64, freq, metric, key string) (common.StatTicks, error) {
	result := common.StatTicks{}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		// Get trade stats bucket
		tradeStatsBk := tx.Bucket([]byte(TRADE_STATS_BUCKET))
		metricBk := tradeStatsBk.Bucket([]byte(metric))
		// metricStats := metricBk.Stats()
		// log.Printf("metric %s bucket stats %+v", metric, metricStats)

		var freqBkName string
		freqBkName, err = getBucketNameByFreq(freq)
		if err != nil {
			return err
		}

		freqBk := metricBk.Bucket([]byte(freqBkName))
		// freqStats := freqBk.Stats()
		// log.Printf("freq %s bucket stats %+v", freqBkName, freqStats)
		c := freqBk.Cursor()
		// min := getTimestampByFreq(fromTime, freq)
		// max := getTimestampByFreq(toTime, freq)
		// log.Printf("from %d to %d", min, max)

		min := uint64ToBytes(fromTime)
		max := uint64ToBytes(toTime)

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			stats := common.TradeStats{}
			err = json.Unmarshal(v, &stats)
			// log.Printf("%v", stats)
			if err != nil {
				return err
			}

			_, ok := stats[key]
			// log.Printf("key: %s", key)
			if ok {
				timestamp := bytesToUint64(k) / 1000000 // to milis
				result[timestamp] = stats[key]
			}
		}
		return nil
	})
	return result, err
}

func (self *BoltStorage) GetAssetVolume(fromTime uint64, toTime uint64, freq string, asset string) (common.StatTicks, error) {
	result, err := self.getTradeStats(fromTime, toTime, freq, ASSETS_VOLUME_BUCKET, asset)
	return result, err
}

func (self *BoltStorage) GetBurnFee(fromTime uint64, toTime uint64, freq string, reserveAddr string) (result common.StatTicks, err error) {
	result, err = self.getTradeStats(fromTime, toTime, freq, BURN_FEE_BUCKET, strings.ToLower(reserveAddr))
	return
}

func (self *BoltStorage) GetWalletFee(fromTime uint64, toTime uint64, freq string, reserveAddr string, walletAddr string) (result common.StatTicks, err error) {
	key := strings.Join([]string{
		strings.ToLower(reserveAddr),
		strings.ToLower(walletAddr),
	}, "_")
	result, err = self.getTradeStats(fromTime, toTime, freq, WALLET_FEE_BUCKET, key)
	return
}

func (self *BoltStorage) GetUserVolume(fromTime uint64, toTime uint64, freq string, userAddr string) (result common.StatTicks, err error) {
	result, err = self.getTradeStats(fromTime, toTime, freq, USER_VOLUME_BUCKET, strings.ToLower(userAddr))
	return
}

func (self *BoltStorage) GetAddressesOfUser(user string) ([]string, error) {
	var err error
	result := []string{}
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ID_ADDRESSES))
		userBucket := b.Bucket([]byte(user))
		if userBucket != nil {
			userBucket.ForEach(func(k, v []byte) error {
				result = append(result, string(k))
				return nil
			})
		}
		return nil
	})
	return result, err
}

func (self *BoltStorage) GetUserOfAddress(addr string) (string, error) {
	addr = strings.ToLower(addr)
	var err error
	var result string
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ADDRESS_ID))
		id := b.Get([]byte(addr))
		result = string(id)
		return nil
	})
	return result, err
}

func (self *BoltStorage) GetCategory(addr string) (string, error) {
	addr = strings.ToLower(addr)
	var err error
	var result string
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ADDRESS_CATEGORY))
		cat := b.Get([]byte(addr))
		result = string(cat)
		return nil
	})
	return result, err
}

func (self *BoltStorage) UpdateUserAddresses(user string, addrs []string) error {
	user = strings.ToLower(user)
	addresses := []string{}
	for _, addr := range addrs {
		addresses = append(addresses, strings.ToLower(addr))
	}
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		for _, address := range addresses {
			// get temp user identity
			b := tx.Bucket([]byte(ADDRESS_ID))
			oldID := b.Get([]byte(address))
			// remove the addresses bucket assocciated to this temp user
			b = tx.Bucket([]byte(ID_ADDRESSES))
			b.DeleteBucket(oldID)
			// update user to each address => user
			b = tx.Bucket([]byte(ADDRESS_ID))
			b.Put([]byte(address), []byte(user))
		}
		// update addresses bucket for real user
		b := tx.Bucket([]byte(ID_ADDRESSES))
		b, err = b.CreateBucketIfNotExists([]byte(user))
		if err != nil {
			return err
		}
		for _, address := range addresses {
			b.Put([]byte(address), []byte{1})
		}
		return nil
	})
	return err
}

func (self *BoltStorage) StoreCatLog(l common.SetCatLog) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		// map address to category
		b := tx.Bucket([]byte(ADDRESS_CATEGORY))
		addrBytes := []byte(strings.ToLower(l.Address.Hex()))
		b.Put(addrBytes, []byte(strings.ToLower(l.Category)))
		// get the user of it
		b = tx.Bucket([]byte(ADDRESS_ID))
		user := b.Get(addrBytes)
		if len(user) == 0 {
			// if the user doesn't exist, we set the user to its address
			user = addrBytes
		}
		// add address to its user addresses
		b = tx.Bucket([]byte(ID_ADDRESSES))
		b, err = b.CreateBucketIfNotExists(user)
		if err != nil {
			return err
		}
		b.Put(addrBytes, []byte{1})
		// add user to map
		b = tx.Bucket([]byte(ADDRESS_ID))
		b.Put(addrBytes, user)
		return nil
	})
	return err
}
