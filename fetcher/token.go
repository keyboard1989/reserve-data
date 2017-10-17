package fetcher

import (
	"fmt"
)

type Token interface {
	SymbolOnExchange(exchange Exchange) string
	Symbol() string
	Name() string
}

type TokenPair struct {
	Base  Token
	Quote Token
}

func (self *TokenPair) PairString() string {
	return fmt.Sprintf("%s-%s", self.Base.Symbol(), self.Quote.Symbol())
}
