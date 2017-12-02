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
	balance  *RamBalanceStorage
	ebalance *RamEBalanceStorage
	rate     *RamRateStorage
	activity *RamActivityStorage
}

func NewRamStorage() *RamStorage {
	return &RamStorage{
		NewRamPriceStorage(),
		NewRamBalanceStorage(),
		NewRamEBalanceStorage(),
		NewRamRateStorage(),
		NewRamActivityStorage(),
	}
}

func (self *RamStorage) CurrentPriceVersion(timepoint uint64) (common.Version, error) {
	version, err := self.price.CurrentVersion(timepoint)
	return common.Version(version), err
}

func (self *RamStorage) CurrentBalanceVersion(timepoint uint64) (common.Version, error) {
	version, err := self.balance.CurrentVersion(timepoint)
	return common.Version(version), err
}

func (self *RamStorage) CurrentEBalanceVersion(timepoint uint64) (common.Version, error) {
	version, err := self.ebalance.CurrentVersion(timepoint)
	return common.Version(version), err
}

func (self *RamStorage) CurrentRateVersion(timepoint uint64) (common.Version, error) {
	version, err := self.rate.CurrentVersion(timepoint)
	return common.Version(version), err
}

func (self *RamStorage) GetAllPrices(version common.Version) (map[common.TokenPairID]common.OnePrice, error) {
	return self.price.GetAllPrices(int64(version))
}

func (self *RamStorage) GetOnePrice(pair common.TokenPairID, version common.Version) (common.OnePrice, error) {
	return self.price.GetOnePrice(pair, int64(version))
}

func (self *RamStorage) GetAllBalances(version common.Version) (map[string]common.BalanceEntry, error) {
	return self.balance.GetAllBalances(int64(version))
}

func (self *RamStorage) GetAllEBalances(version common.Version) (map[common.ExchangeID]common.EBalanceEntry, error) {
	return self.ebalance.GetAllBalances(int64(version))
}

func (self *RamStorage) GetAllRates(version common.Version) (common.AllRateEntry, error) {
	return self.rate.GetRates(int64(version))
}

func (self *RamStorage) StorePrice(data map[common.TokenPairID]common.OnePrice, timepoint uint64) error {
	return self.price.StoreNewData(data, timepoint)
}

func (self *RamStorage) StoreBalance(data map[string]common.BalanceEntry, timepoint uint64) error {
	return self.balance.StoreNewData(data, timepoint)
}

func (self *RamStorage) StoreEBalance(data map[common.ExchangeID]common.EBalanceEntry, timepoint uint64) error {
	return self.ebalance.StoreNewData(data, timepoint)
}

func (self *RamStorage) StoreRate(data common.AllRateEntry, timepoint uint64) error {
	return self.rate.StoreNewData(data, timepoint)
}

func (self *RamStorage) Record(
	action string,
	id common.ActivityID,
	destination string,
	params map[string]interface{}, result map[string]interface{},
	status string,
	timepoint uint64) error {
	return self.activity.StoreNewData(
		action, id, destination,
		params, result, status, timepoint,
	)
}

func (self *RamStorage) GetAllRecords() ([]common.ActivityRecord, error) {
	return self.activity.GetAllRecords()
}

func (self *RamStorage) GetPendingActivities() ([]common.ActivityRecord, error) {
	return self.activity.GetPendingRecords()
}

func (self *RamStorage) UpdateActivityStatus(action string, id common.ActivityID, destination string, status string) error {
	return self.activity.UpdateActivityStatus(action, id, destination, status)
}
