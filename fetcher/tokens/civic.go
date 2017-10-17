package tokens

import (
	"github.com/KyberNetwork/reserve-data/fetcher"
)

type CVC struct{}

func (self CVC) Name() string   { return "Civic" }
func (self CVC) Symbol() string { return "CVC" }
func (self CVC) SymbolOnExchange(exchange fetcher.Exchange) string {
	switch exchange.Name() {
	case "bittrex", "liqui", "poloniex":
		return "CVC"
	default:
		return ""
	}
}
