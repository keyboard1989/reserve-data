package fetcher

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type Storage interface {
	StorePrice(map[common.TokenPairID]common.OnePrice) error
	StoreBalance(map[string]common.BalanceEntry) error
	StoreEBalance(map[common.ExchangeID]common.EBalanceEntry) error
	StoreRate(common.AllRateEntry) error
}
