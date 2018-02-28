package stat

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type ReserveStats struct {
	storage Storage
	fetcher Fetcher
}

func NewReserveStats(storage Storage, fetcher Fetcher) *ReserveStats {
	return &ReserveStats{
		storage: storage,
		fetcher: fetcher,
	}
}

func (self ReserveStats) GetAssetVolume(fromTime, toTime uint64, freq, asset string) ([]common.TradeStats, error) {
	data, err := self.storage.GetAssetVolume(fromTime, toTime, freq, asset)
	return data, err
}

func (self ReserveStats) GetBurnFee(fromTime, toTime uint64, freq, reserveAddr string) ([]common.TradeStats, error) {
	data, err := self.storage.GetBurnFee(fromTime, toTime, freq, reserveAddr)
	return data, err
}

func (self ReserveStats) GetWalletFee(fromTime, toTime uint64, freq, reserveAddr, walletAddr string) ([]common.TradeStats, error) {
	data, err := self.storage.GetWalletFee(fromTime, toTime, freq, reserveAddr, walletAddr)
	return data, err
}

func (self ReserveStats) GetUserVolume(fromTime, toTime uint64, freq, userAddr string) ([]common.TradeStats, error) {
	data, err := self.storage.GetUserVolume(fromTime, toTime, freq, userAddr)
	return data, err
}

func (self ReserveStats) GetTradeLogs(fromTime uint64, toTime uint64) ([]common.TradeLog, error) {
	return self.storage.GetTradeLogs(fromTime, toTime)
}

func (self ReserveStats) Run() error {
	return self.fetcher.Run()
}

func (self ReserveStats) Stop() error {
	return self.fetcher.Stop()
}
