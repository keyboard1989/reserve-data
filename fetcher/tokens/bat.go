package tokens

import (
	"github.com/KyberNetwork/reserve-data/fetcher"
)

type BAT struct{}

func (self BAT) Name() string   { return "Basic Attention Token" }
func (self BAT) Symbol() string { return "BAT" }
func (self BAT) SymbolOnExchange(exchange fetcher.Exchange) string {
	switch exchange.Name() {
	case "bittrex", "liqui":
		return "BAT"
	default:
		return ""
	}
}
