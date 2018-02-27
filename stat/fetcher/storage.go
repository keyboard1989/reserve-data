package fetcher

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type Storage interface {
	LastBlock() (uint64, error)
	UpdateLogBlock(block uint64, timepoint uint64) error
	StoreTradeLog(stat common.TradeLog, timepoint uint64) error
	SetTradeStats(metric, freq string, t uint64, tradeStats common.TradeStats) error
}
