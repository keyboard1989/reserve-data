package storage

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type RamStorage struct {
	price   *RamPriceStorage
	balance *RamBalanceStorage
}

func NewRamStorage() *RamStorage {
	return &RamStorage{
		NewRamPriceStorage(),
		NewRamBalanceStorage(),
	}
}

func (self *RamStorage) CurrentPriceVersion() (common.Version, error) {
	version, err := self.price.CurrentVersion()
	return common.Version(version), err
}

func (self *RamStorage) GetAllPrices(version common.Version) (map[common.TokenPairID]common.OnePrice, error) {
	return self.price.GetAllPrices(int64(version))
}

func (self *RamStorage) GetOnePrice(pair common.TokenPairID, version common.Version) (common.OnePrice, error) {
	return self.price.GetOnePrice(pair, int64(version))
}

func (self *RamStorage) StorePrice(data map[common.TokenPairID]common.OnePrice) error {
	return self.price.StoreNewData(data)
}

func (self *RamStorage) StoreBalance(data map[string]common.RawBalance) error {
	return self.balance.StoreNewData(data)
}
