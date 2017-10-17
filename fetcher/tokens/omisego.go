package tokens

import (
	"github.com/KyberNetwork/reserve-data/fetcher"
)

type OMG struct{}

func (self OMG) Name() string   { return "OmiseGo" }
func (self OMG) Symbol() string { return "OMG" }
func (self OMG) SymbolOnExchange(exchange fetcher.Exchange) string {
	return "OMG"
}
