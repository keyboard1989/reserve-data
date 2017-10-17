package tokens

import (
	"github.com/KyberNetwork/reserve-data/fetcher"
)

type ADX struct{}

func (self ADX) Name() string   { return "AdEx" }
func (self ADX) Symbol() string { return "ADX" }
func (self ADX) SymbolOnExchange(exchange fetcher.Exchange) string {
	switch exchange.Name() {
	case "bittrex", "liqui":
		return "ADX"
	default:
		return ""
	}
}
