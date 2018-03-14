package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/boltdb/bolt"
)

const (
	MAX_NUMBER_VERSION int = 1000

	TRADE_STATS_BUCKET   string = "trade_stats"
	ASSETS_VOLUME_BUCKET string = "assets_volume"
	BURN_FEE_BUCKET      string = "burn_fee"
	WALLET_FEE_BUCKET    string = "wallet_fee"
	USER_VOLUME_BUCKET   string = "user_volume"
	MINUTE_BUCKET        string = "minute"
	HOUR_BUCKET          string = "hour"
	DAY_BUCKET           string = "day"

	RESERVE_RATES string = "reserve_rates"
)

type BoltStatStorage struct {
	db *bolt.DB
}

func uint64ToBytes(u uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, u)
	return b
}

func bytesToUint64(b []byte) uint64 {
	if len(b) != 8 {
		return 0
	}
	return binary.BigEndian.Uint64(b)
}

func NewBoltStatStorage(path string) (*BoltStatStorage, error) {
	// init instance
	var err error
	var db *bolt.DB
	db, err = bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	// init buckets
	db.Update(func(tx *bolt.Tx) error {
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
	storage := &BoltStatStorage{db}
	return storage, nil
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
func (self *BoltStatStorage) GetNumberOfVersion(tx *bolt.Tx, bucket string) int {
	result := 0
	b := tx.Bucket([]byte(bucket))
	c := b.Cursor()
	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		result++
	}
	return result
}

// PruneOutdatedData Remove first version out of database
func (self *BoltStatStorage) PruneOutdatedData(tx *bolt.Tx, bucket string) error {
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

func (self *BoltStatStorage) SetTradeStats(metric, freq string, t uint64, tradeStats common.TradeStats) (err error) {
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

func (self *BoltStatStorage) getTradeStats(fromTime, toTime uint64, freq, metric, key string) (common.StatTicks, error) {
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

func (self *BoltStatStorage) GetAssetVolume(
	fromTime, toTime uint64, freq, asset string) (common.StatTicks, error) {
	result, err := self.getTradeStats(fromTime, toTime, freq, ASSETS_VOLUME_BUCKET, asset)
	return result, err
}

func (self *BoltStatStorage) GetBurnFee(
	fromTime, toTime uint64, freq, reserveAddr string) (result common.StatTicks, err error) {
	result, err = self.getTradeStats(fromTime, toTime, freq, BURN_FEE_BUCKET, strings.ToLower(reserveAddr))
	return
}

func (self *BoltStatStorage) GetWalletFee(
	fromTime, toTime uint64, freq, reserveAddr, walletAddr string) (result common.StatTicks, err error) {
	key := strings.Join([]string{
		strings.ToLower(reserveAddr),
		strings.ToLower(walletAddr),
	}, "_")
	result, err = self.getTradeStats(fromTime, toTime, freq, WALLET_FEE_BUCKET, key)
	return
}

func (self *BoltStatStorage) GetUserVolume(
	fromTime, toTime uint64, freq, userAddr string) (result common.StatTicks, err error) {
	result, err = self.getTradeStats(fromTime, toTime, freq, USER_VOLUME_BUCKET, strings.ToLower(userAddr))
	return
}
