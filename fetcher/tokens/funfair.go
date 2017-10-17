package tokens

import (
	"github.com/KyberNetwork/reserve-data/fetcher"
)

type FUN struct{}

func (self FUN) Name() string   { return "FunFair" }
func (self FUN) Symbol() string { return "FUN" }
func (self FUN) SymbolOnExchange(exchange fetcher.Exchange) string {
	switch exchange.Name() {
	case "bittrex", "binance":
		return "FUN"
	default:
		return ""
	}
}
