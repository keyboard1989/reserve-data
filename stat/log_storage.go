package stat

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type LogStorage interface {
	StoreCatLog(l common.SetCatLog) error
	StoreTradeLog(stat common.TradeLog, timepoint uint64) error
	UpdateLogBlock(block uint64, timepoint uint64) error

	GetTradeLogs(fromTime uint64, toTime uint64) ([]common.TradeLog, error)
	LastBlock() (uint64, error)
}
