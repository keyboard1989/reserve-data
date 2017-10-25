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
	"OMG":  Token{"OMG", "0xa0306f99fa78b3261937f0e5959545c40d95a87a", 18},
	"DGD":  Token{"DGD", "0x9b4c72f6f713b57c2b9e48fe88d0e331d7f85c67", 9},
	"CVC":  Token{"CVC", "0x1964ed8dd6c55a995f9ac19f1172b875326f17f3", 8},
	"FUN":  Token{"FUN", "0x2acb94b917ad26e74c677096f3ca6f06ce5356c4", 8},
	"MCO":  Token{"MCO", "0xb76c34d957a27ce7e28d29c154cf7fbe123df77b", 8},
	"GNT":  Token{"GNT", "0xb2951dabddd6b1d3e49b59ce9631a7db454cf426", 18},
	"ADX":  Token{"ADX", "0x56831302d650a193fc1e16595807b60822150ed6", 4},
	"PAY":  Token{"PAY", "0x2a42c22023b89d3920ec8f572d2496bc969a3935", 18},
	"BAT":  Token{"BAT", "0xc3af5d9fd03980dcee5c1cb7738464a649cb80d7", 18},
	"KNC":  Token{"KNC", "0x81b2e459253e0f6848fa68719878641be90f6736", 18},
	"EOS":  Token{"EOS", "0x4f394ed0d40eabc8db58b03a03fe44daea550180", 18},
	"LINK": Token{"LINK", "0xf1bbe08842b9247a045e47fd734b563ffa391922", 18},
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
