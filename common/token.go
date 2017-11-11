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
	"OMG":  Token{"OMG", "0x879b5383c9e1269be9dcf73ae9740c26b91e9802", 18},
	"DGD":  Token{"DGD", "0xc94c72978bdcc50d763a541695d90a8416f050b2", 9},
	"CVC":  Token{"CVC", "0x91cacf7aea3b0d945a2873eff7595cb4de0d7297", 8},
	"FUN":  Token{"FUN", "0xd3b0286ad5edac328bc9e625327853057e1a0e72", 8},
	"MCO":  Token{"MCO", "0xb3ca241e04f2b9a94b58c9857ce854bd56efc8ee", 8},
	"GNT":  Token{"GNT", "0xee45f2ff517f892e8c0d16b341d66f14a1372cff", 18},
	"ADX":  Token{"ADX", "0xf15f87db547796266cb33da7bd52a9aae6055698", 4},
	"PAY":  Token{"PAY", "0xda0e5f258734959982d58f3b17457f104d6dcb68", 18},
	"BAT":  Token{"BAT", "0xc12e72373eae8f3b901f6d47b7124e025e55fb2b", 18},
	"KNC":  Token{"KNC", "0x88c29c3f40b4e15989176f9546b80a1cff4a6b0d", 18},
	"EOS":  Token{"EOS", "0x44fb6a08ad67ac0b4ef57519de84bda74f99d0f6", 18},
	"LINK": Token{"LINK", "0xdd1de6eac7ff7ae02b139526e024fd53a128c164", 18},
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
