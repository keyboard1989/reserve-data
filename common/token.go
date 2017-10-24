package common

import (
	"errors"
	"fmt"
	"strings"
)

type Token struct {
	ID      string
	Address string
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

var supportedTokens = map[string]Token{
	"ETH":  Token{"ETH", "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"},
	"OMG":  Token{"OMG", "0xd3210dcf33b66d19dd8e9d61fe9fec9c55d73e9a"},
	"DGD":  Token{"DGD", "0xe8a72c46ec0b4b028020895247c306e05d5a671c"},
	"CVC":  Token{"CVC", "0x05aab29c326b6c50cf8b11b9b169bbf39bdfc792"},
	"FUN":  Token{"FUN", "0x922023fccfbb17070ca9d3dd299445e9a2fbf0f2"},
	"MCO":  Token{"MCO", "0xd9a9c05310077250956837494a336f6f24506653"},
	"GNT":  Token{"GNT", "0xa7b94a4720dfc991a9186edaa891b68026aacd88"},
	"ADX":  Token{"ADX", "0x8c04b0b0856fe9fdfc983852675e8e3bc40c5fc1"},
	"PAY":  Token{"PAY", "0x746d341a65d7a14b3f2e63b7a9b47e8233acdc0a"},
	"BAT":  Token{"BAT", "0x20a1b8a052e02c4e1443fc1b9589c6e1c881ddb0"},
	"KNC":  Token{"KNC", "0xbd46bb7cf321b4acf0a703422f4c2dd69ad0dba0"},
	"EOS":  Token{"EOS", "0xab830bf5a5397a67743e1a290c33fdcb0b4a1632"},
	"LINK": Token{"LINK", "0x3700f5565002159f5031bcdeee3b4d99ca34486f"},
}

func GetToken(id string) (Token, error) {
	t := supportedTokens[strings.ToUpper(id)]
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
