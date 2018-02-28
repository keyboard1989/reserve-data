package stat

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type Storage interface {
	GetTradeLogs(fromTime uint64, toTime uint64) ([]common.TradeLog, error)

	GetAssetVolume(fromTime uint64, toTime uint64, freq string, asset string) (common.StatTicks, error)
	GetBurnFee(fromTime uint64, toTime uint64, freq string, reserveAddr string) (common.StatTicks, error)
	GetWalletFee(fromTime uint64, toTime uint64, freq string, reserveAddr string, walletAddr string) (common.StatTicks, error)
	GetUserVolume(fromTime uint64, toTime uint64, freq string, userAddr string) (common.StatTicks, error)
}
