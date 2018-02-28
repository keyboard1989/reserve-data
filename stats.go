package reserve

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type ReserveStats interface {
	GetTradeLogs(fromTime uint64, toTime uint64) ([]common.TradeLog, error)
	GetAssetVolume(fromTime, toTime uint64, freq, asset string) (common.StatTicks, error)
	GetBurnFee(fromTime, toTime uint64, freq, reserveAddr string) (common.StatTicks, error)
	GetWalletFee(fromTime, toTime uint64, freq, reserveAddr, walletAddr string) (common.StatTicks, error)
	GetUserVolume(fromTime, toTime uint64, freq, userAddr string) (common.StatTicks, error)

	ExceedDailyLimit(addr string) (bool, error)

	Run() error
	Stop() error
}
