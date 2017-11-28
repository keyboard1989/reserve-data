package fetcher

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type Storage interface {
	StorePrice(data map[common.TokenPairID]common.OnePrice, timepoint uint64) error
	StoreBalance(data map[string]common.BalanceEntry, timepoint uint64) error
	StoreEBalance(data map[common.ExchangeID]common.EBalanceEntry, timepoint uint64) error
	StoreRate(data common.AllRateEntry, timepoint uint64) error
	StoreOrder(data common.AllOrderEntry, timepoint uint64) error

	GetPendingActivities() ([]common.ActivityRecord, error)
	UpdateActivityStatus(action string, id string, destination string, status string) error
}
