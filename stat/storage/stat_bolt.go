package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/boltdb/bolt"
)

const (
	MAX_NUMBER_VERSION int = 1000

	TRADELOG_PROCESSOR_STATE string = "tradelog_processor_state"

	TRADE_STATS_BUCKET string = "trade_stats"

	MINUTE_BUCKET               string = "minute"
	HOUR_BUCKET                 string = "hour"
	DAY_BUCKET                  string = "day"
	TIMEZONE_BUCKET_PREFIX      string = "utc"
	START_TIMEZONE              int64  = -11
	END_TIMEZONE                int64  = 14
	ADDRESS_BUCKET_PREFIX       string = "address"
	DAILY_ADDRESS_BUCKET_PREFIX string = "daily_address"
	DAILY_USER_BUCKET_PREFIX    string = "daily_user"
	WALLET_ADDRESS_BUCKET       string = "wallet_address"
	RESERVE_RATES               string = "reserve_rates"
	EXPIRED                     uint64 = uint64(7 * time.Hour * 24)
	COUNTRY_BUCKET              string = "country_bucket"
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
		tx.CreateBucket([]byte(WALLET_ADDRESS_BUCKET))
		//create timezone buckets
		tradeStatsBk := tx.Bucket([]byte(TRADE_STATS_BUCKET))
		frequencies := []string{MINUTE_BUCKET, HOUR_BUCKET, DAY_BUCKET}

		for _, freq := range frequencies {
			tradeStatsBk.CreateBucket([]byte(freq))
		}
		for i := START_TIMEZONE; i <= END_TIMEZONE; i++ {
			tzstring := fmt.Sprintf("%s%d", TIMEZONE_BUCKET_PREFIX, i)
			tx.CreateBucket([]byte(tzstring))
			tradeStatsBk.CreateBucket([]byte(tzstring))
			dailyAddrBkname := fmt.Sprintf("%s%d", DAILY_ADDRESS_BUCKET_PREFIX, i)
			tx.CreateBucket([]byte(dailyAddrBkname))
			dailyUserBkname := fmt.Sprintf("%s%d", DAILY_USER_BUCKET_PREFIX, i)
			tx.CreateBucket([]byte(dailyUserBkname))
			addrBucketName := fmt.Sprintf("%s%d", ADDRESS_BUCKET_PREFIX, i)
			tx.CreateBucket([]byte(addrBucketName))
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
	default:
		offset, ok := strconv.ParseInt(strings.TrimPrefix(freq, "utc"), 10, 64)
		if (offset < START_TIMEZONE) || (offset > END_TIMEZONE) {
			err = errors.New("Frequency is wrong, can not get bucket name")
		}
		if ok != nil {
			err = ok
		}
		bucketName = freq
	}
	return
}

func getTimestampByFreq(t uint64, freq string) (result []byte) {
	ui64Day := uint64(time.Hour * 24)
	switch freq {
	case "m", "M":
		result = uint64ToBytes(t / uint64(time.Minute) * uint64(time.Minute))
	case "h", "H":
		result = uint64ToBytes(t / uint64(time.Hour) * uint64(time.Hour))
	case "d", "D":
		result = uint64ToBytes(t / ui64Day * ui64Day)
	default:
		// utc timezone
		offset, _ := strconv.ParseInt(strings.TrimPrefix(freq, "utc"), 10, 64)
		ui64offset := uint64(int64(time.Hour) * offset)
		if offset > 0 {
			result = uint64ToBytes((t+ui64offset)/ui64Day*ui64Day + ui64offset)
		} else {
			offset = 0 - offset
			result = uint64ToBytes((t-ui64offset)/ui64Day*ui64Day - ui64offset)
		}
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

func (self *BoltStatStorage) GetUserStats(timestamp uint64, addr, email, wallet string, kycEd bool, timezone int64) (common.TradeStats, error) {
	stats := common.TradeStats{}
	var err error

	self.db.View(func(tx *bolt.Tx) error {
		freq := fmt.Sprintf("%s%d", TIMEZONE_BUCKET_PREFIX, timezone)
		dailyTimestamp := string(getTimestampByFreq(timestamp, freq))
		dailyAddrBkname := fmt.Sprintf("%s%d", DAILY_ADDRESS_BUCKET_PREFIX, timezone)
		dailyAddrBk := tx.Bucket([]byte(dailyAddrBkname))
		dailyAddrKey := fmt.Sprintf("%s_%s", dailyTimestamp, addr)
		if v := dailyAddrBk.Get([]byte(dailyAddrKey)); v == nil {
			stats["first_trade_in_day"] = 1 // FIRST TRADE IN DAY

			addrBucketName := fmt.Sprintf("%s%d", ADDRESS_BUCKET_PREFIX, timezone)
			addrBk := tx.Bucket([]byte(addrBucketName))
			if v := addrBk.Get([]byte(addr)); v == nil {
				stats["first_trade_ever"] = 1 // FIRST TRADE EVER
			}

			if kycEd {
				dailyUserBkname := fmt.Sprintf("%s%d", DAILY_USER_BUCKET_PREFIX, timezone)
				dailyUserBk := tx.Bucket([]byte(dailyUserBkname))
				dailyUserKey := fmt.Sprintf("%s_%s", dailyTimestamp, email)
				if v := dailyUserBk.Get([]byte(dailyUserKey)); v == nil {
					stats["kyced_in_day"] = 1
				}
			}
		}

		//Get user stat for the current Wallet address
		dailyAddrWalletKey := fmt.Sprintf("%s_%s_%s", dailyTimestamp, addr, wallet)
		if v := dailyAddrBk.Get([]byte(dailyAddrWalletKey)); v == nil {
			stats[fmt.Sprintf("wallet_first_trade_in_day_%s", wallet)] = 1
			addrBucketName := fmt.Sprintf("%s%d", ADDRESS_BUCKET_PREFIX, timezone)
			addrBk := tx.Bucket([]byte(addrBucketName))
			if v := addrBk.Get([]byte(fmt.Sprintf("%s_%s", addr, wallet))); v == nil {
				stats[fmt.Sprintf("wallet_first_trade_ever_%s", wallet)] = 1
			}
			if kycEd {
				dailyUserBkname := fmt.Sprintf("%s%d", DAILY_USER_BUCKET_PREFIX, timezone)
				dailyUserBk := tx.Bucket([]byte(dailyUserBkname))
				dailyUserWalletKey := fmt.Sprintf("%s_%s_%s", dailyTimestamp, email, wallet)
				if v := dailyUserBk.Get([]byte(dailyUserWalletKey)); v == nil {
					stats[fmt.Sprintf("wallet_kyced_in_day_%s", wallet)] = 1
				}
			}

		}

		return nil
	})

	return stats, err
}

func isEalier(k, timestamp []byte) bool {
	temps := strings.Split(string(k), "_")
	return (bytes.Compare([]byte(temps[0]), timestamp) <= 0)
}

func (self *BoltStatStorage) PruneDailyBucket(timepoint uint64, timezone int64) (err error) {
	freq := fmt.Sprintf("%s%d", TIMEZONE_BUCKET_PREFIX, timezone)
	currentUintTimeStamp := bytesToUint64(getTimestampByFreq(timepoint, freq))
	expiredByteTimeStamp := uint64ToBytes(currentUintTimeStamp - EXPIRED)

	//update daily Address bucket
	dailyAddrBkname := fmt.Sprintf("%s%d", DAILY_ADDRESS_BUCKET_PREFIX, timezone)
	self.db.Update(func(tx *bolt.Tx) error {
		dailyAddrBk := tx.Bucket([]byte(dailyAddrBkname))
		c := dailyAddrBk.Cursor()
		for k, _ := c.First(); k != nil && isEalier(k, expiredByteTimeStamp); k, _ = c.Next() {
			err = dailyAddrBk.Delete([]byte(k))
		}
		return err
	})
	//update daily User bucket
	dailyUserBkname := fmt.Sprintf("%s%d", DAILY_USER_BUCKET_PREFIX, timezone)
	self.db.Update(func(tx *bolt.Tx) error {
		dailyAddrBk := tx.Bucket([]byte(dailyUserBkname))
		c := dailyAddrBk.Cursor()
		for k, _ := c.First(); k != nil && isEalier(k, expiredByteTimeStamp); k, _ = c.Next() {
			err = dailyAddrBk.Delete([]byte(k))
		}
		return err
	})
	return
}

func (self *BoltStatStorage) SetUserStats(timestamp uint64, addr, email, wallet string, kycEd bool, timezone int64, stats common.TradeStats) error {
	var err error
	freq := fmt.Sprintf("%s%d", TIMEZONE_BUCKET_PREFIX, timezone)
	if err = self.SetTradeStats(freq, timestamp, stats); err != nil {
		return err
	}
	self.PruneDailyBucket(timestamp, timezone)

	self.db.Update(func(tx *bolt.Tx) error {

		dailyTimestamp := string(getTimestampByFreq(timestamp, freq))
		if _, traded := stats["first_trade_in_day"]; traded {
			dailyAddrBkname := fmt.Sprintf("%s%d", DAILY_ADDRESS_BUCKET_PREFIX, timezone)
			dailyAddrBk := tx.Bucket([]byte(dailyAddrBkname))
			dailyAddrKey := fmt.Sprintf("%s_%s", dailyTimestamp, addr)
			if err := dailyAddrBk.Put([]byte(dailyAddrKey), []byte("1")); err != nil {
				return err
			}
			if _, traded := stats["first_trade_ever"]; traded {
				addrBucketName := fmt.Sprintf("%s%d", ADDRESS_BUCKET_PREFIX, timezone)
				addrBk := tx.Bucket([]byte(addrBucketName))
				if err := addrBk.Put([]byte(addr), []byte("1")); err != nil {
					return err
				}
			}

			if kycEd {
				dailyUserBkname := fmt.Sprintf("%s%d", DAILY_USER_BUCKET_PREFIX, timezone)
				dailyUserBk := tx.Bucket([]byte(dailyUserBkname))
				dailyUserKey := fmt.Sprintf("%s_%s", dailyTimestamp, email)
				if err := dailyUserBk.Put([]byte(dailyUserKey), []byte("1")); err != nil {
					return err
				}
			}
		}

		//Set user stat for the current Wallet address
		if _, walletTraded := stats[fmt.Sprintf("wallet_first_trade_in_day_%s", wallet)]; walletTraded {
			dailyAddrBkname := fmt.Sprintf("%s%d", DAILY_ADDRESS_BUCKET_PREFIX, timezone)
			dailyAddrBk := tx.Bucket([]byte(dailyAddrBkname))
			dailyAddrWalletKey := fmt.Sprintf("%s_%s_%s", dailyTimestamp, addr, wallet)
			if err := dailyAddrBk.Put([]byte(dailyAddrWalletKey), []byte("1")); err != nil {
				return err
			}

			if _, walletTraded := stats[fmt.Sprintf("wallet_first_trade_ever_%s", wallet)]; walletTraded {
				addrBucketName := fmt.Sprintf("%s%d", ADDRESS_BUCKET_PREFIX, timezone)
				addrBk := tx.Bucket([]byte(addrBucketName))
				if err := addrBk.Put([]byte(fmt.Sprintf("%s_%s", addr, wallet)), []byte("1")); err != nil {
					return err
				}
			}
			if kycEd {
				dailyUserBkname := fmt.Sprintf("%s%d", DAILY_USER_BUCKET_PREFIX, timezone)
				dailyUserBk := tx.Bucket([]byte(dailyUserBkname))
				dailyUserWalletKey := fmt.Sprintf("%s_%s_%s", dailyTimestamp, email, wallet)
				if err := dailyUserBk.Put([]byte(dailyUserWalletKey), []byte("1")); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

func (self *BoltStatStorage) SetWalletAddress(walletAddr string) (err error) {
	self.db.Update(func(tx *bolt.Tx) error {
		walletBucket := tx.Bucket([]byte(WALLET_ADDRESS_BUCKET))

		if err := walletBucket.Put([]byte(walletAddr), []byte("1")); err != nil {
			return err
		}

		return nil
	})
	return
}

func (self *BoltStatStorage) GetWalletAddress() ([]string, error) {
	var result []string
	self.db.View(func(tx *bolt.Tx) error {
		walletBucket := tx.Bucket([]byte(WALLET_ADDRESS_BUCKET))
		c := walletBucket.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			result = append(result, string(k[:]))
		}
		return nil
	})
	return result, nil
}

func (self *BoltStatStorage) GetWalletStats(fromTime uint64, toTime uint64, walletAddr string, timezone int64) (common.StatTicks, error) {

	walletAddr = strings.ToLower(walletAddr)
	log.Printf("fromtime %v toTime %v fre %v", fromTime, toTime, timezone)
	result := common.StatTicks{}
	tzstring := fmt.Sprintf("%s%d", TIMEZONE_BUCKET_PREFIX, timezone)
	stats, err := self.getTradeStats(fromTime, toTime, tzstring)
	if err != nil {
		return result, err
	}

	for timestamp, stat := range stats {
		trade_countstr, found := stat[fmt.Sprintf("wallet_trade_count_%s", walletAddr)]
		if !found {
			continue
		}
		tradeCount := float64(trade_countstr)
		var avgEth, avgUsd float64
		if tradeCount > 0 {
			avgEth = float64(stat[fmt.Sprintf("wallet_eth_volume_%s", walletAddr)]) / tradeCount
			avgUsd = float64(stat[fmt.Sprintf("wallet_usd_volume_%s", walletAddr)]) / tradeCount
		}
		result[timestamp] = map[string]float64{
			"total_eth_volume":     stat[fmt.Sprintf("wallet_eth_volume_%s", walletAddr)],
			"total_usd_amount":     stat[fmt.Sprintf("wallet_usd_volume_%s", walletAddr)],
			"total_burn_fee":       stat[fmt.Sprintf("wallet_burn_fee_%s", walletAddr)],
			"unique_addresses":     stat[fmt.Sprintf("wallet_first_trade_in_day_%s", walletAddr)],
			"new_unique_addresses": stat[fmt.Sprintf("wallet_first_trade_ever_%s", walletAddr)],
			"kyced_addresses":      stat[fmt.Sprintf("wallet_kyced_in_day_%s", walletAddr)],
			"total_trade":          tradeCount,
			"eth_per_trade":        avgEth,
			"usd_per_trade":        avgUsd,
		}
	}
	return result, nil
}

func (self *BoltStatStorage) SetCountry(country string) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(COUNTRY_BUCKET))
		err = b.Put([]byte(country), []byte("1"))
		return err
	})
	return err
}

func (self *BoltStatStorage) GetCountries() ([]string, error) {
	countries := []string{}
	var err error
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(COUNTRY_BUCKET))
		c := b.Cursor()
		for _, v := c.First(); v != nil; _, v = c.Next() {
			countries = append(countries, string(v))
		}
		return err
	})
	return countries, err
}

func (self *BoltStatStorage) GetAssetVolume(fromTime uint64, toTime uint64, freq string, asset string) (common.StatTicks, error) {
	result := common.StatTicks{}

	stats, err := self.getTradeStats(fromTime, toTime, freq)
	eth := common.SupportedTokens["ETH"]

	for timestamp, stat := range stats {
		if strings.ToLower(eth.Address) == asset {
			result[timestamp] = map[string]float64{
				"volume":     stat[fmt.Sprintf("assets_volume_%s", asset)],
				"usd_amount": stat[fmt.Sprintf("assets_usd_amount_%s", asset)],
			}
		} else {
			result[timestamp] = map[string]float64{
				"volume":     stat[fmt.Sprintf("assets_volume_%s", asset)],
				"eth_amount": stat[fmt.Sprintf("assets_eth_amount_%s", asset)],
				"usd_amount": stat[fmt.Sprintf("assets_usd_amount_%s", asset)],
			}
		}
	}

	return result, err
}

func (self *BoltStatStorage) GetBurnFee(fromTime uint64, toTime uint64, freq string, reserveAddr string) (common.StatTicks, error) {
	result := common.StatTicks{}

	stats, err := self.getTradeStats(fromTime, toTime, freq)
	for timestamp, stat := range stats {
		result[timestamp] = stat[fmt.Sprintf("burn_fee_%s", reserveAddr)]
	}

	return result, err
}

func (self *BoltStatStorage) GetWalletFee(fromTime uint64, toTime uint64, freq string, reserveAddr string, walletAddr string) (common.StatTicks, error) {
	result := common.StatTicks{}

	stats, err := self.getTradeStats(fromTime, toTime, freq)
	for timestamp, stat := range stats {
		result[timestamp] = stat[fmt.Sprintf("wallet_fee_%s_%s", reserveAddr, walletAddr)]
	}

	return result, err
}

func (self *BoltStatStorage) GetUserVolume(fromTime uint64, toTime uint64, freq string, userAddr string) (common.StatTicks, error) {
	result := common.StatTicks{}

	stats, err := self.getTradeStats(fromTime, toTime, freq)
	for timestamp, stat := range stats {
		result[timestamp] = stat[fmt.Sprintf("user_volume_%s", userAddr)]
	}
	return result, err
}

func (self *BoltStatStorage) GetTradeSummary(fromTime uint64, toTime uint64, timezone int64) (common.StatTicks, error) {
	log.Printf("fromtime %v toTime %v fre %v", fromTime, toTime, timezone)
	result := common.StatTicks{}
	tzstring := fmt.Sprintf("%s%d", TIMEZONE_BUCKET_PREFIX, timezone)
	stats, err := self.getTradeStats(fromTime, toTime, tzstring)
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
