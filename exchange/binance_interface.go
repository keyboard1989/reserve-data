package exchange

import (
	"math/big"
	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
	"sync"
)

type BinanceInterface interface {
	FetchOnePairData(
		wg *sync.WaitGroup,
		pair common.TokenPair,
		data *sync.Map,
		timepoint uint64)

	GetInfo(timepoint uint64) (Binainfo, error)

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
