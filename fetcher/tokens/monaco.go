package tokens

import (
	"github.com/KyberNetwork/reserve-data/fetcher"
)

type MCO struct{}

func (self MCO) Name() string   { return "Monaco" }
func (self MCO) Symbol() string { return "MCO" }
func (self MCO) SymbolOnExchange(exchange fetcher.Exchange) string {
	switch exchange.Name() {
	case "bittrex", "liqui", "binance":
		return "MCO"
	default:
		return ""
	}
}
