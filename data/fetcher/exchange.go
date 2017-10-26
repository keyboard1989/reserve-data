package fetcher

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type Exchange interface {
	ID() common.ExchangeID
	Name() string
	FetchPriceData() (map[common.TokenPairID]common.ExchangePrice, error)
	FetchEBalanceData() (common.EBalanceEntry, error)
}
