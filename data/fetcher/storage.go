package fetcher

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type Storage interface {
	StorePrice(data map[common.TokenPairID]common.OnePrice, timepoint uint64) error
	StoreRate(data common.AllRateEntry, timepoint uint64) error
	StoreAuthSnapshot(data *common.AuthDataSnapshot, timepoint uint64) error

	GetPendingActivities() ([]common.ActivityRecord, error)
	UpdateActivity(id common.ActivityID, act common.ActivityRecord) error
}
