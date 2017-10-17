package tokens

import (
	"github.com/KyberNetwork/reserve-data/fetcher"
)

type PAY struct{}

func (self PAY) Name() string   { return "TenX" }
func (self PAY) Symbol() string { return "PAY" }
func (self PAY) SymbolOnExchange(exchange fetcher.Exchange) string {
	switch exchange.Name() {
	case "bittrex", "liqui":
		return "PAY"
	default:
		return ""
	}
}
