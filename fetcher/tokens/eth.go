package tokens

import (
	"github.com/KyberNetwork/reserve-data/fetcher"
)

type ETH struct{}

func (self ETH) Name() string   { return "Ether" }
func (self ETH) Symbol() string { return "ETH" }
func (self ETH) SymbolOnExchange(exchange fetcher.Exchange) string {
	return "ETH"
}
