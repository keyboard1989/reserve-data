package data

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type Storage interface {
	CurrentPriceVersion(timepoint uint64) (common.Version, error)
	GetAllPrices(common.Version) (map[common.TokenPairID]common.OnePrice, error)
	GetOnePrice(common.TokenPairID, common.Version) (common.OnePrice, error)

	CurrentBalanceVersion(timepoint uint64) (common.Version, error)
	GetAllBalances(common.Version) (map[string]common.BalanceEntry, error)

	CurrentEBalanceVersion(timepoint uint64) (common.Version, error)
	GetAllEBalances(common.Version) (map[common.ExchangeID]common.EBalanceEntry, error)

	CurrentRateVersion(timepoint uint64) (common.Version, error)
	GetAllRates(common.Version) (common.AllRateEntry, error)
}
