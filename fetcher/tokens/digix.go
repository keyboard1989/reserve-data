package tokens

import (
	"github.com/KyberNetwork/reserve-data/fetcher"
)

type DGD struct{}

func (self DGD) Name() string   { return "Digix" }
func (self DGD) Symbol() string { return "DGD" }
func (self DGD) SymbolOnExchange(exchange fetcher.Exchange) string {
	switch exchange.Name() {
	case "bittrex", "liqui":
		return "DGD"
	default:
		return ""
	}
}
