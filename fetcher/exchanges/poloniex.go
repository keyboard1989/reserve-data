package exchanges

import (
	"github.com/KyberNetwork/reserve-data/fetcher"
	"github.com/KyberNetwork/reserve-data/market"
)

type Poloniex struct {
}

func (self Poloniex) ID() string {
	return "poloniex"
}

func (self Poloniex) Name() string {
	return "poloniex"
}

func (self Poloniex) FetchData(pairs []fetcher.TokenPair) (string, map[string]market.ExchangeData) {
	// TODO
	return "timestamp", map[string]market.ExchangeData{}
}
