package common

import (
	"errors"
	"fmt"

	ethereum "github.com/ethereum/go-ethereum/common"
)

type Exchange interface {
	ID() ExchangeID
	Address(token Token) (address ethereum.Address, supported bool)
}

var SupportedExchanges = map[ExchangeID]Exchange{}

func GetExchange(id string) (Exchange, error) {
	ex := SupportedExchanges[ExchangeID(id)]
	if ex == nil {
		return ex, errors.New(fmt.Sprintf("Exchange %s is not supported", id))
	} else {
		return ex, nil
	}
}
