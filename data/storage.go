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
	GetRate(common.Version) (common.AllRateEntry, error)
	GetRates(fromTime, toTime uint64) ([]common.AllRateEntry, error)

	GetAllRecords(fromTime, toTime uint64) ([]common.ActivityRecord, error)
	GetPendingActivities() ([]common.ActivityRecord, error)

	GetTradeLogs(fromTime uint64, toTime uint64) ([]common.TradeLog, error)
	GetTradeHistory(common.Version) (common.AllTradeHistory, error)
}
