package fetcher

import (
	"github.com/KyberNetwork/reserve-data/market"
)

type Exchange interface {
	Name() string
	ID() string
	FetchData(pairs []TokenPair) (string, map[string]*market.ExchangeData)
}
