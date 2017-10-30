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
	"OMG":  Token{"OMG", "0xd19034874ba588b5d4cb0c998afd6852ad8f4718", 18},
	"DGD":  Token{"DGD", "0xd667cd12cafeeb55678461f7925549e26f07a7a5", 9},
	"CVC":  Token{"CVC", "0xf3e24a4e6df726d839d44aa21badb8770763855a", 8},
	"FUN":  Token{"FUN", "0xf0282a5a5d3c1bb06162e1524d51fdf2c6daadb0", 8},
	"MCO":  Token{"MCO", "0xd8a32ef63c0b21fe276dfd04950a745f6f158a33", 8},
	"GNT":  Token{"GNT", "0xb1ea7e08bc945b699548dba15a0296b5f0ce60b1", 18},
	"ADX":  Token{"ADX", "0xcad1ebc2244738722f43efe4d1859bd16f7ad87d", 4},
	"PAY":  Token{"PAY", "0x6022dce3b1f6d953d9f0256b2aa24da99bbd3187", 18},
	"BAT":  Token{"BAT", "0x209dba78da77cba9ff7d788512a8128e4ab854a8", 18},
	"KNC":  Token{"KNC", "0x03b50798dcc087953a5bd2e36e6112ad1092ceed", 18},
	"EOS":  Token{"EOS", "0x9f0ca0b0c755006f313d951502b1c603b5edeaf1", 18},
	"LINK": Token{"LINK", "0x61d3e3b963add2a43d6a93608b3c66b5baef111d", 18},
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
