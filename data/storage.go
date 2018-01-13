package data

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type Storage interface {
	CurrentPriceVersion(timepoint uint64) (common.Version, error)
	GetAllPrices(common.Version) (common.AllPriceEntry, error)
	GetOnePrice(common.TokenPairID, common.Version) (common.OnePrice, error)

	CurrentAuthDataVersion(timepoint uint64) (common.Version, error)
	GetAuthData(common.Version) (common.AuthDataSnapshot, error)

	CurrentRateVersion(timepoint uint64) (common.Version, error)
	GetAllRates(common.Version) (common.AllRateEntry, error)

	GetAllRecords() ([]common.ActivityRecord, error)
	GetPendingActivities() ([]common.ActivityRecord, error)

	GetTradeLogs(reserve ethereum.Address, fromTime uint64, toTime uint64) ([]common.TradeLog, error)
}
