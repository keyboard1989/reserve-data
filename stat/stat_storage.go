package stat

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type StatStorage interface {
	GetAssetVolume(fromTime uint64, toTime uint64, freq string, asset string) (common.StatTicks, error)
	GetBurnFee(fromTime uint64, toTime uint64, freq string, reserveAddr string) (common.StatTicks, error)
	GetWalletFee(fromTime uint64, toTime uint64, freq string, reserveAddr string, walletAddr string) (common.StatTicks, error)
	GetUserVolume(fromTime uint64, toTime uint64, freq string, userAddr string) (common.StatTicks, error)
	GetTradeSummary(fromTime, toTime uint64) (common.StatTicks, error)
	GetLastProcessedTradeLogTimepoint() (timepoint uint64, err error)

	// SaveUserAddress(timestamp uint64, addr string) (common.TradeStats, error)

	SetTradeStats(freq string, t uint64, tradeStats common.TradeStats) error
	SetLastProcessedTradeLogTimepoint(timepoint uint64) error
}
