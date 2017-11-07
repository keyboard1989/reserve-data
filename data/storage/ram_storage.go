package storage

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type RamStorage struct {
	price    *RamPriceStorage
	balance  *RamBalanceStorage
	ebalance *RamEBalanceStorage
	rate     *RamRateStorage
}

func NewRamStorage() *RamStorage {
	return &RamStorage{
		NewRamPriceStorage(),
		NewRamBalanceStorage(),
		NewRamEBalanceStorage(),
		NewRamRateStorage(),
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

func (self *RamStorage) StorePrice(data map[common.TokenPairID]common.OnePrice) error {
	return self.price.StoreNewData(data)
}

func (self *RamStorage) StoreBalance(data map[string]common.BalanceEntry) error {
	return self.balance.StoreNewData(data)
}

func (self *RamStorage) StoreEBalance(data map[common.ExchangeID]common.EBalanceEntry) error {
	return self.ebalance.StoreNewData(data)
}

func (self *RamStorage) StoreRate(data common.AllRateEntry) error {
	return self.rate.StoreNewData(data)
}
