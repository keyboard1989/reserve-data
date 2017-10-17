package tokens

import (
	"github.com/KyberNetwork/reserve-data/fetcher"
)

type EOS struct{}

func (self EOS) Name() string   { return "EOS" }
func (self EOS) Symbol() string { return "EOS" }
func (self EOS) SymbolOnExchange(exchange fetcher.Exchange) string {
	switch exchange.Name() {
	case "binance", "liqui", "bitfinex":
		return "EOS"
	default:
		return ""
	}
}
