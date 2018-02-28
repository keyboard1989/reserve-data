package stat

import (
	"errors"
	"fmt"
	"strings"
	"time"

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

func validateTimeWindow(fromTime, toTime uint64, freq string) (uint64, uint64, error) {
	var from = fromTime * 1000000
	var to = toTime * 1000000

	switch freq {
	case "m", "M":
		if to-from > uint64((time.Hour * 24).Nanoseconds()) {
			return 0, 0, errors.New("Minute frequency limit is 1 day")
		}
	case "h", "H":
		if to-from > uint64((time.Hour * 24 * 180).Nanoseconds()) {
			return 0, 0, errors.New("Hour frequency limit is 180 days")
		}
	case "d", "D":
		if to-from > uint64((time.Hour * 24 * 365 * 3).Nanoseconds()) {
			return 0, 0, errors.New("Day frequency limit is 3 years")
		}
	default:
		return 0, 0, errors.New("Invalid frequencies")
	}
	return from, to, nil
}

func (self ReserveStats) GetAssetVolume(fromTime, toTime uint64, freq, asset string) ([]common.TradeStats, error) {
	data := []common.TradeStats{}

	fromTime, toTime, err := validateTimeWindow(fromTime, toTime, freq)
	if err != nil {
		return data, err
	}

	token, err := common.GetToken(asset)
	if err != nil {
		return data, errors.New(fmt.Sprintf("assets %s is not supported", asset))
	}

	data, err = self.storage.GetAssetVolume(fromTime, toTime, freq, strings.ToLower(token.Address))
	return data, err
}

func (self ReserveStats) GetBurnFee(fromTime, toTime uint64, freq, reserveAddr string) ([]common.TradeStats, error) {
	data := []common.TradeStats{}

	fromTime, toTime, err := validateTimeWindow(fromTime, toTime, freq)
	if err != nil {
		return data, err
	}

	data, err = self.storage.GetBurnFee(fromTime, toTime, freq, reserveAddr)
	return data, err
}

func (self ReserveStats) GetWalletFee(fromTime, toTime uint64, freq, reserveAddr, walletAddr string) ([]common.TradeStats, error) {
	data := []common.TradeStats{}

	fromTime, toTime, err := validateTimeWindow(fromTime, toTime, freq)
	if err != nil {
		return data, err
	}

	data, err = self.storage.GetWalletFee(fromTime, toTime, freq, reserveAddr, walletAddr)
	return data, err
}

func (self ReserveStats) GetUserVolume(fromTime, toTime uint64, freq, userAddr string) ([]common.TradeStats, error) {
	data := []common.TradeStats{}

	fromTime, toTime, err := validateTimeWindow(fromTime, toTime, freq)
	if err != nil {
		return data, err
	}

	data, err = self.storage.GetUserVolume(fromTime, toTime, freq, userAddr)
	return data, err
}

func (self ReserveStats) Run() error {
	return self.fetcher.Run()
}

func (self ReserveStats) Stop() error {
	return self.fetcher.Stop()
}
