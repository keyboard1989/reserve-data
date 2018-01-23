package common

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Token struct {
	ID      string
	Address string
	Decimal int64
}

func (self Token) MarshalText() (text []byte, err error) {
	// return []byte(fmt.Sprintf(
	// 	"%s-%s", self.ID, self.Address,
	// )), nil
	return []byte(self.ID), nil
}

func (self Token) IsETH() bool {
	return self.ID == "ETH"
}

type TokenPair struct {
	Base  Token
	Quote Token
}

type TargetQty struct {
	ReserveTargetQty float64 `json:"reserve_target"`
	TotalTargetQty   float64 `json:"total_target"`
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

var SupportedTokens map[string]Token
var TokenTargetQty map[string]TargetQty

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

func UpdateTokenTargetQty(data string) error {
	tokensValues := strings.Split(data, "|")
	for _, tok := range tokensValues {
		values := strings.Split(tok, "_")
		t := SupportedTokens[strings.ToUpper(values[0])]
		if t.ID == "" {
			return errors.New(fmt.Sprintf("Token %s is not supported", values[0]))
		}
		totalValue, _ := strconv.ParseFloat(values[1], 64)
		reserveValue, _ := strconv.ParseFloat(values[2], 64)
		TokenTargetQty[values[0]] = TargetQty{
			TotalTargetQty:   totalValue,
			ReserveTargetQty: reserveValue,
		}
	}
	return nil
}
