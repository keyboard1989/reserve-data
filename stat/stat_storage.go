package stat

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type StatStorage interface {
	GetAssetVolume(fromTime, toTime uint64, freq, asset string) (common.StatTicks, error)
	GetBurnFee(fromTime, toTime uint64, freq, reserveAddr string) (common.StatTicks, error)
	GetWalletFee(fromTime, toTime uint64, freq, reserveAddr string, walletAddr string) (common.StatTicks, error)
	GetUserVolume(fromTime, toTime uint64, freq, userAddr string) (common.StatTicks, error)
	GetTradeSummary(fromTime, toTime uint64, timezone int64) (common.StatTicks, error)
	GetLastProcessedTradeLogTimepoint() (timepoint uint64, err error)
	GetWalletStats(fromTime, toTime uint64, walletAddr string, timezone int64) (common.StatTicks, error)
	GetUserStats(timestamp uint64, addr, email, wallet, country string, kycEd bool, timezone int64) (common.TradeStats, error)
	SetUserStats(timestamp uint64, addr, email, wallet, country string, kycEd bool, timezone int64, stats common.TradeStats) error

	SetTradeStats(freq string, t uint64, tradeStats common.TradeStats) error
	SetWalletAddress(walletAddr string) error
	GetWalletAddress() ([]string, error)
	SetLastProcessedTradeLogTimepoint(timepoint uint64) error

	SetCountry(country string) error
	GetCountries() ([]string, error)
	SetCountryStat(countryStats map[string]common.CountryStatsTimeZone) error
	GetCountryStats(fromTime, toTime uint64, country string, tzparam int64) (common.StatTicks, error)
}
