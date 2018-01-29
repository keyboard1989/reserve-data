package storage

import (
	"github.com/KyberNetwork/reserve-data/common"
)

// RamStorage is a simple and fast storage that eliminate all old
// data but the newest one
// RamStorage works fine when Core doesn't use historical data
// and doesn't take timestamp into account
//
// If Core uses such data, please use other kind of storage such as
// bolt.
type RamStorage struct {
	price    *RamPriceStorage
	auth     *RamAuthStorage
	rate     *RamRateStorage
	activity *RamActivityStorage
	log      *RamLogStorage
	bittrex  *RamBittrexStorage
}

func NewRamStorage() *RamStorage {
	return &RamStorage{
		NewRamPriceStorage(),
		NewRamAuthStorage(),
		NewRamRateStorage(),
		NewRamActivityStorage(),
		NewRamLogStorage(),
		NewRamBittrexStorage(),
	}
}

func (self *RamStorage) CurrentPriceVersion(timepoint uint64) (common.Version, error) {
	version, err := self.price.CurrentVersion(timepoint)
	return common.Version(version), err
}

func (self *RamStorage) CurrentAuthDataVersion(timepoint uint64) (common.Version, error) {
	version, err := self.auth.CurrentVersion(timepoint)
	return common.Version(version), err
}

func (self *RamStorage) CurrentRateVersion(timepoint uint64) (common.Version, error) {
	version, err := self.rate.CurrentVersion(timepoint)
	return common.Version(version), err
}

func (self *RamStorage) GetAllPrices(version common.Version) (common.AllPriceEntry, error) {
	return self.price.GetAllPrices(int64(version))
}

func (self *RamStorage) GetOnePrice(pair common.TokenPairID, version common.Version) (common.OnePrice, error) {
	return self.price.GetOnePrice(pair, int64(version))
}

func (self *RamStorage) GetAuthData(version common.Version) (common.AuthDataSnapshot, error) {
	return self.auth.GetSnapshot(int64(version))
}

func (self *RamStorage) GetRate(version common.Version) (common.AllRateEntry, error) {
	return self.rate.GetRate(int64(version))
}

func (self *RamStorage) GetRates(fromTime, toTime uint64) ([]common.AllRateEntry, error) {
	return self.rate.GetRates(fromTime, toTime)
}

func (self *RamStorage) StorePrice(data common.AllPriceEntry, timepoint uint64) error {
	return self.price.StoreNewData(data, timepoint)
}

func (self *RamStorage) StoreAuthSnapshot(
	data *common.AuthDataSnapshot,
	timepoint uint64) error {
	return self.auth.StoreNewSnapshot(data, timepoint)
}

func (self *RamStorage) StoreRate(data common.AllRateEntry, timepoint uint64) error {
	return self.rate.StoreNewData(data, timepoint)
}

func (self *RamStorage) UpdateActivity(id common.ActivityID, activity common.ActivityRecord) error {
	return self.activity.UpdateActivity(id, activity)
}

func (self *RamStorage) Record(
	action string,
	id common.ActivityID,
	destination string,
	params map[string]interface{}, result map[string]interface{},
	estatus string,
	mstatus string,
	timepoint uint64) error {
	return self.activity.StoreNewData(
		action, id, destination,
		params, result, estatus, mstatus, timepoint,
	)
}

func (self *RamStorage) GetAllRecords(fromTime, toTime uint64) ([]common.ActivityRecord, error) {
	return self.activity.GetAllRecords(fromTime, toTime)
}

func (self *RamStorage) GetPendingActivities() ([]common.ActivityRecord, error) {
	return self.activity.GetPendingRecords()
}

func (self *RamStorage) IsNewBittrexDeposit(id uint64, actID common.ActivityID) bool {
	return self.bittrex.IsNewDeposit(id, actID)
}

func (self *RamStorage) RegisterBittrexDeposit(id uint64, actID common.ActivityID) error {
	return self.bittrex.RegisterDeposit(id, actID)
}

func (self *RamStorage) HasPendingDeposit(token common.Token, exchange common.Exchange) bool {
	return self.activity.HasPendingDeposit(token, exchange)
}

func (self *RamStorage) UpdateLogBlock(block uint64, timepoint uint64) error {
	return self.log.UpdateLogBlock(block, timepoint)
}

func (self *RamStorage) LastBlock() (uint64, error) {
	return self.log.LastBlock()
}

func (self *RamStorage) GetTradeLogs(fromTime uint64, toTime uint64) ([]common.TradeLog, error) {
	return self.log.GetTradeLogs(fromTime, toTime)
}

func (self *RamStorage) StoreTradeLog(stat common.TradeLog, timepoint uint64) error {
	return self.log.StoreTradeLog(stat, timepoint)
}
