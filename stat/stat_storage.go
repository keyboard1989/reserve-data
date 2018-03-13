package stat

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type StatStorage interface {
	GetAssetVolume(fromTime uint64, toTime uint64, freq string, asset string) (common.StatTicks, error)
	GetBurnFee(fromTime uint64, toTime uint64, freq string, reserveAddr string) (common.StatTicks, error)
	GetWalletFee(fromTime uint64, toTime uint64, freq string, reserveAddr string, walletAddr string) (common.StatTicks, error)
	GetUserVolume(fromTime uint64, toTime uint64, freq string, userAddr string) (common.StatTicks, error)

	SetTradeStats(metric, freq string, t uint64, tradeStats common.TradeStats) error
}
