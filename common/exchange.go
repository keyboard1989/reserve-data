package common

import (
	"errors"
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum/common"
)

type Exchange interface {
	ID() ExchangeID
	Address(token Token) (address ethereum.Address, supported bool)
	Withdraw(token Token, amount *big.Int, address ethereum.Address, timepoint uint64) (ethereum.Hash, error)
	Trade(tradeType string, base Token, quote Token, rate float64, amount float64, timepoint uint64) (done float64, remaining float64, finished bool, err error)
	MarshalText() (text []byte, err error)
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
