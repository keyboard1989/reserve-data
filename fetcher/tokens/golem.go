package tokens

import (
	"github.com/KyberNetwork/reserve-data/fetcher"
)

type GNT struct{}

func (self GNT) Name() string   { return "Golem" }
func (self GNT) Symbol() string { return "GNT" }
func (self GNT) SymbolOnExchange(exchange fetcher.Exchange) string {
	switch exchange.Name() {
	case "bittrex", "liqui", "binance":
		return "GNT"
	default:
		return ""
	}
}
