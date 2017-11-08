package exchange

import (
	"math/big"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type BitfinexInterface interface {
	Depth(tokens string, timepoint uint64) (Bitfresp, error)

	GetInfo(timepoint uint64) (Bitfinfo, error)

	Withdraw(
		token common.Token,
		amount *big.Int,
		address ethereum.Address,
		timepoint uint64) error

	Trade(
		tradeType string,
		base, quote common.Token,
		rate, amount float64,
		timepoint uint64) (done float64, remaining float64, finished bool, err error)
}
