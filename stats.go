package reserve

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type ReserveStats interface {
	GetTradeLogs(fromTime uint64, toTime uint64) ([]common.TradeLog, error)
	GetAssetVolume(fromTime, toTime uint64, freq, asset string) ([]common.TradeStats, error)
	GetBurnFee(fromTime, toTime uint64, freq, reserveAddr string) ([]common.TradeStats, error)
	GetWalletFee(fromTime, toTime uint64, freq, reserveAddr, walletAddr string) ([]common.TradeStats, error)
	GetUserVolume(fromTime, toTime uint64, freq, userAddr string) ([]common.TradeStats, error)

	Run() error
	Stop() error
}
