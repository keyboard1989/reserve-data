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
	"OMG":  Token{"OMG", "0x5f2396de67edfaed541f12ec4f4017e7e244a9ca", 18},
	"DGD":  Token{"DGD", "0x7ab3b88c28cdaa51ec3029759eb1a9ed8ceb5653", 9},
	"CVC":  Token{"CVC", "0x2b3a0529c8b2004aa3ed0e2e9a5f0fca94316eb1", 8},
	"FUN":  Token{"FUN", "0xf10b6be8fd8f1b4b7bc691a24e67f05ddc8d9bab", 8},
	"MCO":  Token{"MCO", "0x3c57a51144d0c9ff8511380a5d42a08b7f90289c", 8},
	"GNT":  Token{"GNT", "0x103cbe3519bfd739b01279c64f15d504773996ee", 18},
	"ADX":  Token{"ADX", "0x93aeb71991d4e5c893b9a27e7c666af7bcb9d957", 4},
	"PAY":  Token{"PAY", "0x4f1f07247055e66c2c7c9333e36a6c1a47ab02e3", 18},
	"BAT":  Token{"BAT", "0xa4cae8a0edfd148132b21d6639914e5eae0b58eb", 18},
	"KNC":  Token{"KNC", "0x744660550f19d8843d9dd5be8dc3ecf06b611952", 18},
	"EOS":  Token{"EOS", "0xf7625ea45c26843b7e3dd515904b5197531bf8e9", 18},
	"LINK": Token{"LINK", "0xa23a77deb2e77d29dff7033ff16a6046f9c29796", 18},
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
