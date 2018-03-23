package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/boltdb/bolt"
)

const (
	MAX_NUMBER_VERSION int = 1000

	TRADELOG_PROCESSOR_STATE string = "tradelog_processor_state"

	TRADE_STATS_BUCKET string = "trade_stats"

	MINUTE_BUCKET string = "minute"
	HOUR_BUCKET   string = "hour"
	DAY_BUCKET    string = "day"

	ADDRESS_BUCKET       string = "address"
	DAILY_ADDRESS_BUCKET string = "daily_address"
	DAILY_USER_BUCKET    string = "daily_user"

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
		tx.CreateBucket([]byte(TRADELOG_PROCESSOR_STATE))

		tx.CreateBucket([]byte(ADDRESS_BUCKET))
		tx.CreateBucket([]byte(DAILY_ADDRESS_BUCKET))
		tx.CreateBucket([]byte(DAILY_USER_BUCKET))

		tradeStatsBk := tx.Bucket([]byte(TRADE_STATS_BUCKET))
		frequencies := []string{MINUTE_BUCKET, HOUR_BUCKET, DAY_BUCKET}

		for _, freq := range frequencies {
			tradeStatsBk.CreateBucket([]byte(freq))
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

func (self *BoltStatStorage) SetLastProcessedTradeLogTimepoint(timepoint uint64) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TRADELOG_PROCESSOR_STATE))
		err = b.Put([]byte("last_timepoint"), uint64ToBytes(timepoint))
		return err
	})
	return err
}

func (self *BoltStatStorage) GetLastProcessedTradeLogTimepoint() (uint64, error) {
	var result uint64
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TRADELOG_PROCESSOR_STATE))
		result = bytesToUint64(b.Get([]byte("last_timepoint")))
		return nil
	})
	return result, nil
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

func (self *BoltStatStorage) SetTradeStats(freq string, t uint64, tradeStats common.TradeStats) (err error) {
	self.db.Update(func(tx *bolt.Tx) error {
		tradeStatsBk := tx.Bucket([]byte(TRADE_STATS_BUCKET))

		freqBkName, err := getBucketNameByFreq(freq)
		if err != nil {
			return err
		}
		freqBk := tradeStatsBk.Bucket([]byte(freqBkName))

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

func (self *BoltStatStorage) getTradeStats(fromTime, toTime uint64, freq string) (map[uint64]common.TradeStats, error) {
	result := map[uint64]common.TradeStats{}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		// Get trade stats bucket
		tradeStatsBk := tx.Bucket([]byte(TRADE_STATS_BUCKET))

		var freqBkName string
		freqBkName, err = getBucketNameByFreq(freq)
		if err != nil {
			return err
		}
		freqBk := tradeStatsBk.Bucket([]byte(freqBkName))
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

			timestamp := bytesToUint64(k) / 1000000 // to milis
			result[timestamp] = stats
		}
		return nil
	})
	return result, err
}

func (self *BoltStatStorage) GetUserStats(timestamp uint64, addr, email string, kycEd bool) (common.TradeStats, error) {
	stats := common.TradeStats{}
	var err error

	self.db.View(func(tx *bolt.Tx) error {
		dailyTimestamp := string(getTimestampByFreq(timestamp, "D"))

		dailyAddrBk := tx.Bucket([]byte(DAILY_ADDRESS_BUCKET))
		dailyAddrKey := strings.Join([]string{dailyTimestamp, addr}, "_")
		if v := dailyAddrBk.Get([]byte(dailyAddrKey)); v == nil {
			stats["first_trade_in_day"] = 1 // FIRST TRADE IN DAY

			addrBk := tx.Bucket([]byte(ADDRESS_BUCKET))
			if v := addrBk.Get([]byte(addr)); v == nil {
				stats["first_trade_ever"] = 1 // FIRST TRADE EVER
			}

			if kycEd {
				dailyUserBk := tx.Bucket([]byte(DAILY_USER_BUCKET))
				dailyUserKey := strings.Join([]string{dailyTimestamp, email}, "_")
				if v := dailyUserBk.Get([]byte(dailyUserKey)); v == nil {
					stats["kyced_in_day"] = 1
				}
			}
		}

		return nil
	})

	return stats, err
}

func (self *BoltStatStorage) SetUserStats(timestamp uint64, addr, email string, kycEd bool, stats common.TradeStats) error {
	var err error

	if err = self.SetTradeStats("D", timestamp, stats); err != nil {
		return err
	}

	self.db.Update(func(tx *bolt.Tx) error {

		dailyTimestamp := string(getTimestampByFreq(timestamp, "D"))

		if _, traded := stats["first_trade_in_day"]; traded {
			dailyAddrBk := tx.Bucket([]byte(DAILY_ADDRESS_BUCKET))
			dailyAddrKey := strings.Join([]string{dailyTimestamp, addr}, "_")
			if err := dailyAddrBk.Put([]byte(dailyAddrKey), []byte("1")); err != nil {
				return err
			}

			if _, traded := stats["first_trade_ever"]; traded {
				addrBk := tx.Bucket([]byte(ADDRESS_BUCKET))
				if err := addrBk.Put([]byte(addr), []byte("1")); err != nil {
					return err
				}
			}

			if kycEd {
				dailyUserBk := tx.Bucket([]byte(DAILY_USER_BUCKET))
				dailyUserKey := strings.Join([]string{dailyTimestamp, email}, "_")
				if err := dailyUserBk.Put([]byte(dailyUserKey), []byte("1")); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

func findField(fieldName string, stat common.TradeStats) float64 {
	val, found := stat[fieldName]
	if !found {
		return 0
	}
	return val
}

func (self *BoltStatStorage) GetWalletStats(fromTime uint64, toTime uint64, freq string, walletAddr string) ([]common.WalletStats, error) {
	log.Println("LIEMTEST: 3 got here")
	walletAddr = strings.ToLower(walletAddr)
	log.Printf("LIEMTEST: wlladdr is %s \n", walletAddr)
	var result []common.WalletStats
	stats, err := self.getTradeStats(fromTime, toTime, freq)
	log.Printf("LIEMTEST: len of stats is %d", len(stats))
	if err != nil {
		return result, err
	}
	for timestamp, stat := range stats {
		log.Printf("LIEMTEST: stat is %v", stat)

		eth_amount, found := stat[strings.Join([]string{"wallet_eth_volume", walletAddr}, "_")]
		if !found {
			continue
		}
		usd_amount, _ := stat[strings.Join([]string{"wallet_usd_volume", walletAddr}, "_")]
		burn_fee, _ := stat[strings.Join([]string{"wallet_burn_fee", walletAddr}, "_")]
		trade_count, _ := stat[strings.Join([]string{"wallet_trade_count", walletAddr}, "_")]
		record := common.WalletStats{
			timestamp,
			eth_amount,
			usd_amount,
			burn_fee,
			uint64(trade_count),
		}
		result = append(result, record)

	}
	return result, nil
}

func (self *BoltStatStorage) GetAssetVolume(fromTime uint64, toTime uint64, freq string, asset string) (common.StatTicks, error) {
	result := common.StatTicks{}

	stats, err := self.getTradeStats(fromTime, toTime, freq)
	eth := common.SupportedTokens["ETH"]

	for timestamp, stat := range stats {
		if strings.ToLower(eth.Address) == asset {
			result[timestamp] = map[string]float64{
				"volume":     stat[strings.Join([]string{"assets_volume", asset}, "_")],
				"usd_amount": stat[strings.Join([]string{"assets_usd_amount", asset}, "_")],
			}
		} else {
			result[timestamp] = map[string]float64{
				"volume":     stat[strings.Join([]string{"assets_volume", asset}, "_")],
				"eth_amount": stat[strings.Join([]string{"assets_eth_amount", asset}, "_")],
				"usd_amount": stat[strings.Join([]string{"assets_usd_amount", asset}, "_")],
			}
		}
	}

	return result, err
}

func (self *BoltStatStorage) GetBurnFee(fromTime uint64, toTime uint64, freq string, reserveAddr string) (common.StatTicks, error) {
	result := common.StatTicks{}

	stats, err := self.getTradeStats(fromTime, toTime, freq)
	for timestamp, stat := range stats {
		result[timestamp] = stat[strings.Join([]string{"burn_fee", reserveAddr}, "_")]
	}

	return result, err
}

func (self *BoltStatStorage) GetWalletFee(fromTime uint64, toTime uint64, freq string, reserveAddr string, walletAddr string) (common.StatTicks, error) {
	result := common.StatTicks{}

	stats, err := self.getTradeStats(fromTime, toTime, freq)
	for timestamp, stat := range stats {
		result[timestamp] = stat[strings.Join([]string{"wallet_fee", reserveAddr, walletAddr}, "_")]
	}

	return result, err
}

func (self *BoltStatStorage) GetUserVolume(fromTime uint64, toTime uint64, freq string, userAddr string) (common.StatTicks, error) {
	result := common.StatTicks{}

	stats, err := self.getTradeStats(fromTime, toTime, freq)
	for timestamp, stat := range stats {
		result[timestamp] = stat[strings.Join([]string{"user_volume", userAddr}, "_")]
	}
	return result, err
}

func (self *BoltStatStorage) GetTradeSummary(fromTime uint64, toTime uint64) (common.StatTicks, error) {
	result := common.StatTicks{}

	stats, err := self.getTradeStats(fromTime, toTime, "D")
	if err != nil {
		return result, err
	}

	for timestamp, stat := range stats {
		tradeCount := float64(stat["trade_count"])
		var avgEth, avgUsd float64
		if tradeCount > 0 {
			avgEth = float64(stat["eth_volume"]) / tradeCount
			avgUsd = float64(stat["usd_volume"]) / tradeCount
		}

		result[timestamp] = map[string]float64{
			"total_eth_volume":     stat["eth_volume"],
			"total_usd_amount":     stat["usd_volume"],
			"total_burn_fee":       stat["burn_fee"],
			"unique_addresses":     stat["first_trade_in_day"],
			"new_unique_addresses": stat["first_trade_ever"],
			"kyced_addresses":      stat["kyced_in_day"],
			"total_trade":          tradeCount,
			"eth_per_trade":        avgEth,
			"usd_per_trade":        avgUsd,
		}
	}

	return result, err
}
