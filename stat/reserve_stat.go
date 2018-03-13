package stat

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type ReserveStats struct {
	statStorage StatStorage
	logStorage  LogStorage
	userStorage UserStorage
	rateStorage RateStorage
	fetcher     *Fetcher
}

func NewReserveStats(
	statStorage StatStorage,
	logStorage LogStorage,
	rateStorage RateStorage,
	userStorage UserStorage,
	fetcher *Fetcher) *ReserveStats {
	return &ReserveStats{
		statStorage: statStorage,
		logStorage:  logStorage,
		rateStorage: rateStorage,
		userStorage: userStorage,
		fetcher:     fetcher,
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

func (self ReserveStats) GetAssetVolume(fromTime, toTime uint64, freq, asset string) (common.StatTicks, error) {
	data := common.StatTicks{}

	fromTime, toTime, err := validateTimeWindow(fromTime, toTime, freq)
	if err != nil {
		return data, err
	}

	token, err := common.GetToken(asset)
	if err != nil {
		return data, errors.New(fmt.Sprintf("assets %s is not supported", asset))
	}

	data, err = self.statStorage.GetAssetVolume(fromTime, toTime, freq, strings.ToLower(token.Address))
	return data, err
}

func (self ReserveStats) GetBurnFee(fromTime, toTime uint64, freq, reserveAddr string) (common.StatTicks, error) {
	data := common.StatTicks{}

	fromTime, toTime, err := validateTimeWindow(fromTime, toTime, freq)
	if err != nil {
		return data, err
	}

	data, err = self.statStorage.GetBurnFee(fromTime, toTime, freq, reserveAddr)
	return data, err
}

func (self ReserveStats) GetWalletFee(fromTime, toTime uint64, freq, reserveAddr, walletAddr string) (common.StatTicks, error) {
	data := common.StatTicks{}

	fromTime, toTime, err := validateTimeWindow(fromTime, toTime, freq)
	if err != nil {
		return data, err
	}

	data, err = self.statStorage.GetWalletFee(fromTime, toTime, freq, reserveAddr, walletAddr)
	return data, err
}

func (self ReserveStats) GetUserVolume(fromTime, toTime uint64, freq, userAddr string) (common.StatTicks, error) {
	data := common.StatTicks{}

	fromTime, toTime, err := validateTimeWindow(fromTime, toTime, freq)
	if err != nil {
		return data, err
	}

	data, err = self.statStorage.GetUserVolume(fromTime, toTime, freq, userAddr)
	return data, err
}

func (self ReserveStats) GetTradeLogs(fromTime uint64, toTime uint64) ([]common.TradeLog, error) {
	return self.logStorage.GetTradeLogs(fromTime, toTime)
}

func (self ReserveStats) GetPendingAddresses() ([]string, error) {
	return self.userStorage.GetPendingAddresses()
}

func (self ReserveStats) Run() error {
	return self.fetcher.Run()
}

func (self ReserveStats) Stop() error {
	return self.fetcher.Stop()
}

func (self ReserveStats) GetCapByAddress(addr ethereum.Address) (*common.UserCap, error) {
	category, err := self.userStorage.GetCategory(addr.Hex())
	if err != nil {
		return nil, err
	}
	if category == "0x4" {
		return common.KycedCap(), nil
	} else {
		return common.NonKycedCap(), nil
	}
}

func (self ReserveStats) GetCapByUser(userID string) (*common.UserCap, error) {
	addresses, _, err := self.userStorage.GetAddressesOfUser(userID)
	if err != nil {
		return nil, err
	}
	if len(addresses) == 0 {
		log.Printf("Couldn't find any associated addresses. User %s is not kyced.", userID)
		return common.NonKycedCap(), nil
	} else {
		return self.GetCapByAddress(ethereum.HexToAddress(addresses[0]))
	}
}

func (self ReserveStats) GetReserveRates(fromTime, toTime uint64, reserveAddr ethereum.Address) ([]common.ReserveRates, error) {
	var result []common.ReserveRates
	var err error
	result, err = self.rateStorage.GetReserveRates(fromTime, toTime, reserveAddr.Hex())
	log.Printf("Get reserve rate: %v", result)
	return result, err
}

func (self ReserveStats) UpdateUserAddresses(userID string, addrs []ethereum.Address, timestamps []uint64) error {
	addresses := []string{}
	for _, addr := range addrs {
		addresses = append(addresses, addr.Hex())
	}
	return self.userStorage.UpdateUserAddresses(userID, addresses, timestamps)
}

func (self ReserveStats) ExceedDailyLimit(address ethereum.Address) (bool, error) {
	user, _, err := self.userStorage.GetUserOfAddress(address.Hex())
	if err != nil {
		return false, err
	}
	addrs, _, err := self.userStorage.GetAddressesOfUser(user)
	if err != nil {
		return false, err
	}
	today := common.GetTimepoint() / uint64(24*time.Hour/time.Millisecond) * uint64(24*time.Hour/time.Millisecond)
	var totalVolume float64 = 0.0
	for _, addr := range addrs {
		volumeStats, err := self.GetUserVolume(today-1, today, "D", addr)
		if err == nil {
			log.Printf("volumes: %+v", volumeStats)
			if len(volumeStats) == 0 {
			} else if len(volumeStats) > 1 {
				log.Printf("Got more than 1 day stats. This is a bug in GetUserVolume")
			} else {
				for _, volume := range volumeStats {
					totalVolume += volume
					break
				}
			}
		}
	}
	cap, err := self.GetCapByAddress(address)
	if err == nil && totalVolume >= cap.DailyLimit {
		return true, nil
	} else {
		return false, nil
	}
}
