package fetcher

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type Storage interface {
	StorePrice(data common.AllPriceEntry, timepoint uint64) error
	StoreRate(data common.AllRateEntry, timepoint uint64) error
	StoreAuthSnapshot(data *common.AuthDataSnapshot, timepoint uint64) error
	StoreTradeLog(stat common.TradeLog, timepoint uint64) error
	UpdateLogBlock(block uint64, timepoint uint64) error

	LastBlock() (uint64, error)
	GetPendingActivities() ([]common.ActivityRecord, error)
	UpdateActivity(id common.ActivityID, act common.ActivityRecord) error
}
