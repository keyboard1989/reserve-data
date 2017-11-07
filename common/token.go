package common

import (
	"errors"
	"fmt"
	"strings"
)

type Token struct {
	ID      string
	Address string
	Decimal int64
}

func (self Token) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf(
		"%s-%s", self.ID, self.Address,
	)), nil
}

func (self Token) IsETH() bool {
	return self.ID == "ETH"
}

type TokenPair struct {
	Base  Token
	Quote Token
}

func (self *TokenPair) PairID() TokenPairID {
	return NewTokenPairID(self.Base.ID, self.Quote.ID)
}

func NewTokenPair(base, quote string) (TokenPair, error) {
	bToken, err1 := GetToken(base)
	qToken, err2 := GetToken(quote)
	if err1 != nil || err2 != nil {
		return TokenPair{}, errors.New(fmt.Sprintf("%s or %s is not supported", base, quote))
	} else {
		return TokenPair{bToken, qToken}, nil
	}
}

func MustCreateTokenPair(base, quote string) TokenPair {
	pair, err := NewTokenPair(base, quote)
	if err != nil {
		panic(err)
	} else {
		return pair
	}
}

var SupportedTokens = map[string]Token{
	"ETH":  Token{"ETH", "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", 18},
	"OMG":  Token{"OMG", "0x6b662ffde8f1d2240eb4eefa211463be0eb258a1", 18},
	"DGD":  Token{"DGD", "0xd27763c026260bb8cfcf47a3d2ca18f03cb9da55", 9},
	"CVC":  Token{"CVC", "0x3d1bdb333d4bbd0bf84519c506c953ef869ef179", 8},
	"FUN":  Token{"FUN", "0x0f679d211f23764c3020e2dca0d6277b9abb5b72", 8},
	"MCO":  Token{"MCO", "0xf596502b120689a119dd961b77426e6866e73d2a", 8},
	"GNT":  Token{"GNT", "0x79c800440c5ebac80a8072e7659fa0c7c92da7df", 18},
	"ADX":  Token{"ADX", "0x66cd4fbe38c31094682b9b8cbe306efb4fde895f", 4},
	"PAY":  Token{"PAY", "0xcccc987398f87cc3b14d29e951ba779e3a4b30b7", 18},
	"BAT":  Token{"BAT", "0x8726f7961b39c0a49501b943874ac92ed7240559", 18},
	"KNC":  Token{"KNC", "0xb4ac19f6495df29f32878182be06a2f0572f9763", 18},
	"EOS":  Token{"EOS", "0x07ae1a78a58b01f077b3ca700d352a3db1e11392", 18},
	"LINK": Token{"LINK", "0x829e5df8ba4014021a3b3ba4232c54e9c17ddf70", 18},
}

func GetToken(id string) (Token, error) {
	t := SupportedTokens[strings.ToUpper(id)]
	if t.ID == "" {
		return t, errors.New(fmt.Sprintf("Token %s is not supported", id))
	} else {
		return t, nil
	}
}

func MustGetToken(id string) Token {
	t, e := GetToken(id)
	if e != nil {
		panic(e)
	}
	return t
}
