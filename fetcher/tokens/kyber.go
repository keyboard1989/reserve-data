package tokens

import (
	"github.com/KyberNetwork/reserve-data/fetcher"
)

type KNC struct{}

func (self KNC) Name() string   { return "KyberNetwork Crystal" }
func (self KNC) Symbol() string { return "KNC" }
func (self KNC) SymbolOnExchange(exchange fetcher.Exchange) string {
	switch exchange.Name() {
	case "binance", "liqui":
		return "KNC"
	default:
		return ""
	}
}
