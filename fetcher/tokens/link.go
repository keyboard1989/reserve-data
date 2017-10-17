package tokens

import (
	"github.com/KyberNetwork/reserve-data/fetcher"
)

type LINK struct{}

func (self LINK) Name() string   { return "Chainlink" }
func (self LINK) Symbol() string { return "LINK" }
func (self LINK) SymbolOnExchange(exchange fetcher.Exchange) string {
	switch exchange.Name() {
	case "binance":
		return "LINK"
	default:
		return ""
	}
}
